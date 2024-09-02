package swan

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/swanchain/go-swan-sdk/contract"
	"log"
	"math/big"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type APIClient struct {
	apiKey         string
	httpClient     *HttpClient
	contractDetail ContractDetail
}

func NewAPIClient(apiKey string, testnet ...bool) *APIClient {
	host := gatewayMainnet
	if len(testnet) > 0 && testnet[0] {
		host = gatewayTestnet
	}

	header := make(http.Header)
	header.Add("Authorization", "Bearer "+apiKey)

	var apiClient = APIClient{
		apiKey:     apiKey,
		httpClient: NewHttpClient(host, header),
	}

	contractDetail, err := apiClient.getContractInfo(false)
	if err != nil {
		log.Fatalf("failed to get contract detail, error: %v", err)
	}
	apiClient.contractDetail = contractDetail
	return &apiClient
}

func (c *APIClient) Hardwares() ([]*Hardware, error) {
	var result HardwareResult

	if err := c.httpClient.Get(apiMachines, nil, NewResult(&result)); err != nil {
		return nil, err
	}

	return result.Hardware, nil
}

func (c *APIClient) TaskInfo(taskUUID string) (*TaskInfo, error) {
	var result TaskInfo

	if err := c.httpClient.Get(fmt.Sprintf("%s/%s", apiTask, taskUUID), nil, NewResult(&result)); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *APIClient) Tasks(req *TaskQueryReq) (total int64, list []*TaskInfo, err error) {
	api := apiTasks
	if req != nil {
		api += fmt.Sprintf("?wallet=%s&size=%d&page=%d", req.Wallet, req.Size, req.Page)

	}

	var result PageResult
	result.List = &list

	if err = c.httpClient.Get(api, nil, NewResult(&result)); err != nil {
		return
	}
	total = result.Total
	return
}

func (c *APIClient) CreateTask(req *CreateTaskReq) (CreateTaskResp, error) {
	var createTaskResp CreateTaskResp

	if req.WalletAddress == "" {
		return createTaskResp, fmt.Errorf("no wallet_address provided, please pass in a wallet_address")
	}

	if req.AutoPay && req.PrivateKey == "" {
		return createTaskResp, fmt.Errorf("please provide private_key if using auto_pay")
	}

	if req.Region == "" {
		req.Region = "global"
	}

	if req.StartIn == 0 {
		req.StartIn = 300
	}

	if req.Duration == 0 {
		req.Duration = 3600
	}

	if strings.TrimSpace(req.InstanceType) == "" {
		req.InstanceType = "C1ae.small"
	}

	if _, err := c.getHardwareByInstanceType(req.InstanceType); err != nil {
		return createTaskResp, err
	}
	log.Printf("Using %s machine, region=%s  duration=%d (seconds) \n", req.InstanceType, req.Region, req.Duration)

	if req.JobSourceUri == "" {
		if req.RepoUri != "" {
			sourceUri, err := c.GetSourceUri(req.RepoUri, req.WalletAddress, req.InstanceType, req.RepoBranch, req.RepoOwner, req.RepoName)
			if err != nil {
				return createTaskResp, fmt.Errorf("please provide JobSourceUri, or RepoUri, error: %v", err)
			}
			req.JobSourceUri = sourceUri
		}
	}

	if req.JobSourceUri == "" {
		return createTaskResp, fmt.Errorf("cannot get JobSourceUri. make sure `RepoUri` or `JobSourceUri` is correct")
	}

	var preferredCp string
	if len(req.PreferredCpList) > 0 {
		preferredCp = strings.Join(req.PreferredCpList, ",")
	}

	var params map[string]interface{}
	if c.verifyHardwareRegion(req.InstanceType, req.Region) {
		params["duration"] = req.Duration
		params["cfg_name"] = req.InstanceType
		params["region"] = req.Region
		params["start_in"] = req.StartIn
		params["wallet"] = req.WalletAddress
		params["job_source_uri"] = req.JobSourceUri
	} else {
		return createTaskResp, fmt.Errorf("no %s machine in %s", req.InstanceType, req.Region)
	}

	if preferredCp != "" {
		params["preferred_cp"] = preferredCp
	}

	var task Task
	if err := c.httpClient.PostJSON(apiTask, params, NewResult(&task)); err != nil {
		return createTaskResp, fmt.Errorf("failed to create task, error: %v", err)
	}
	taskUuid := task.UUID

	createTaskResp.Task = task
	createTaskResp.TaskUuid = taskUuid
	createTaskResp.InstanceType = req.InstanceType
	createTaskResp.Id = taskUuid

	estimatePrice, err := c.EstimatePayment(req.InstanceType, req.Duration)
	if err != nil {
		return createTaskResp, err
	}
	createTaskResp.Price = estimatePrice

	var txHash string
	if req.AutoPay {
		payment, err := c.PayAndDeployTask(taskUuid, req.PrivateKey, req.Duration, req.InstanceType)
		if err != nil {
			return createTaskResp, err
		}
		createTaskResp.ConfigOrder = payment.ConfigOrder
		createTaskResp.TxHash = payment.TxHash
		log.Printf("Task created successfully, taskUuid=%s, txHash=%s, instanceType=%s", taskUuid, txHash, req.InstanceType)
	}
	return createTaskResp, nil
}

func (c *APIClient) PayAndDeployTask(taskUuid, privateKey string, duration int, instanceType string) (PaymentResult, error) {
	var paymentResult PaymentResult

	if strings.TrimSpace(instanceType) == "" {
		return paymentResult, fmt.Errorf("invalid instanceType")
	}
	if privateKey == "" {
		return paymentResult, fmt.Errorf("no privateKey provided")
	}

	submitPaymentTx, err := c.submitPayment(taskUuid, privateKey, duration, instanceType)
	if err != nil {
		return paymentResult, err
	}
	time.Sleep(3 * time.Second)
	validatePaymentResult, err := c.validatePayment(submitPaymentTx, taskUuid)
	if err != nil {
		return paymentResult, err
	}
	paymentResult.ConfigOrder = validatePaymentResult.ConfigOrder
	log.Printf("Payment submitted and validated successfully, taskUuid=%s, tx_hash=%s", taskUuid, submitPaymentTx)

	return paymentResult, nil

}

func (c *APIClient) EstimatePayment(instanceType string, duration int) (int64, error) {
	hardwareBaseInfo, err := c.getHardwareByInstanceType(instanceType)
	if err != nil {
		return 0, err
	}

	priceInt, err := strconv.ParseInt(hardwareBaseInfo.Price, 10, 64)
	if err != nil {
		return 0, err
	}
	return priceInt * int64(duration/3600), nil
}

func (c *APIClient) GetSourceUri(repoUri, walletAddress string, instanceType string, repoBranch, repoOwner, repoName string) (string, error) {
	var jobSourceUriResult JobSourceUriResult

	hardwareBaseInfo, err := c.getHardwareByInstanceType(instanceType)
	if err != nil {
		return "", err
	}

	if walletAddress == "" {
		return "", fmt.Errorf("no wallet_address provided")
	}
	var reqData = map[string]interface{}{
		"repo_owner":     repoOwner,
		"repo_name":      repoName,
		"repo_branch":    repoBranch,
		"wallet_address": walletAddress,
		"hardware_id":    hardwareBaseInfo.ID,
		"repo_uri":       repoUri,
	}

	if err := c.httpClient.PostJSON(apiSourceUri, reqData, NewResult(&jobSourceUriResult)); err != nil {
		return "", err
	}
	return jobSourceUriResult.JobSourceUri, nil
}

func (c *APIClient) ReNewTask(taskUuid, txHash, privateKey string, instanceType string, duration int, autoPay bool) (*ReNewTaskResp, error) {
	if strings.TrimSpace(instanceType) == "" {
		return nil, fmt.Errorf("invalid instanceType")
	}

	if !autoPay && privateKey == "" && txHash == "" {
		return nil, fmt.Errorf("auto_pay off or tx_hash not provided, please provide a tx_hash or set auto_pay to True and provide private_key")
	}
	if txHash == "" {
		// renew_payment
		reNewPaymentTxHash, err := c.RenewPayment(taskUuid, privateKey, instanceType, duration)
		if err != nil {
			return nil, err
		}
		txHash = reNewPaymentTxHash
	} else {
		log.Printf("Using given payment transaction hash, txHash=%s", txHash)
	}

	if txHash != "" && taskUuid != "" {
		var params = map[string]interface{}{
			"task_uuid": taskUuid,
			"duration":  duration,
			"tx_hash":   txHash,
		}

		var reNewTaskResp ReNewTaskResp
		if err := c.httpClient.PostJSON(apiReNewTask, params, NewResult(&reNewTaskResp)); err != nil {
			return nil, err
		}
		return &reNewTaskResp, nil
	} else {
		return nil, fmt.Errorf("txHash or taskUuid invalid")
	}

}

func (c *APIClient) RenewPayment(taskUuid, privateKey string, instanceType string, duration int) (string, error) {
	if strings.TrimSpace(instanceType) == "" {
		return "", fmt.Errorf("invalid instanceType")
	}
	if privateKey == "" {
		return "", fmt.Errorf("no privateKey provided")
	}

	hardwareBaseInfo, err := c.getHardwareByInstanceType(instanceType)
	if err != nil {
		return "", err
	}
	var hardwareId = hardwareBaseInfo.ID

	estimatePrice, err := c.EstimatePayment(instanceType, duration)
	if err != nil {
		return "", err
	}
	priceBigInt := new(big.Int).SetInt64(estimatePrice)

	client, err := ethclient.Dial(c.contractDetail.RpcUrl)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// call token contract approve
	tokenContract, err := contract.NewToken(common.HexToAddress(c.contractDetail.SwanTokenContractAddress), client)
	if err != nil {
		return "", err
	}

	// call payment contract submit_payment
	paymentContract, err := contract.NewPaymentContract(common.HexToAddress(c.contractDetail.ClientContractAddress), client)
	if err != nil {
		return "", err
	}

	// client contract address
	approveTransactOpts, err := CreateTransactOpts(client, privateKey)
	if err != nil {
		return "", err
	}
	approve, err := tokenContract.Approve(approveTransactOpts, common.HexToAddress(c.contractDetail.ClientContractAddress), priceBigInt)
	if err != nil {
		return "", err
	}
	log.Printf("token approve tx: %v \n", approve.Hash().String())

	paymentTransactOpts, err := CreateTransactOpts(client, privateKey)
	if err != nil {
		return "", err
	}

	hardwareIdBigInt := new(big.Int).SetInt64(hardwareId)
	durationBigInt := new(big.Int).SetInt64(int64(duration))
	transaction, err := paymentContract.RenewPayment(paymentTransactOpts, taskUuid, hardwareIdBigInt, durationBigInt)
	if err != nil {
		return "", fmt.Errorf("failed to renew payment, error: %v", err)
	}
	log.Printf("txHash: %v", transaction.Hash().String())
	return transaction.Hash().String(), nil
}

func (c *APIClient) verifyHardwareRegion(instanceType, region string) bool {
	hardwareList, err := c.Hardwares()
	if err != nil {
		log.Printf("failed to get hardware, error: %v", err)
		return false
	}

	for _, hardware := range hardwareList {
		if hardware.Name == instanceType {
			for _, r := range hardware.Region {
				if region == r || (strings.ToLower(region) == "global" && hardware.Status == "available") {
					return true
				}
			}
		}
	}
	return false
}

// submitPayment Submit payment for a task
//
// Args:
// task_uuid: unique id returned by `swan_api.create_task`
// hardware_id: id of cp/hardware configuration set
// duration: duration of service runtime (seconds).
//
// Returns:
// tx_hash
func (c *APIClient) submitPayment(taskUuid, privateKey string, duration int, instanceType string) (string, error) {
	hardwareBaseInfo, err := c.getHardwareByInstanceType(instanceType)
	if err != nil {
		return "", err
	}

	if privateKey == "" {
		return "", fmt.Errorf("no privateKey provided")
	}
	var hardwareId = hardwareBaseInfo.ID
	estimatePrice, err := c.EstimatePayment(instanceType, duration)
	if err != nil {
		return "", err
	}
	priceBigInt := new(big.Int).SetInt64(estimatePrice)

	client, err := ethclient.Dial(c.contractDetail.RpcUrl)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// call token contract approve
	tokenContract, err := contract.NewToken(common.HexToAddress(c.contractDetail.SwanTokenContractAddress), client)
	if err != nil {
		return "", err
	}

	// call payment contract submit_payment
	paymentContract, err := contract.NewPaymentContract(common.HexToAddress(c.contractDetail.ClientContractAddress), client)
	if err != nil {
		return "", err
	}

	// client contract address
	approveTransactOpts, err := CreateTransactOpts(client, privateKey)
	if err != nil {
		return "", err
	}

	hardwareIdBigInt := new(big.Int).SetInt64(hardwareId)
	durationBigInt := new(big.Int).SetInt64(int64(duration))
	approve, err := tokenContract.Approve(approveTransactOpts, common.HexToAddress(c.contractDetail.ClientContractAddress), priceBigInt)
	if err != nil {
		return "", err
	}
	log.Printf("token approve tx: %v \n", approve.Hash().String())

	paymentTransactOpts, err := CreateTransactOpts(client, privateKey)
	if err != nil {
		return "", err
	}
	transaction, err := paymentContract.SubmitPayment(paymentTransactOpts, taskUuid, hardwareIdBigInt, durationBigInt)
	if err != nil {
		return "", fmt.Errorf("failed to submit payment, error: %v", err)
	}
	log.Printf("Payment submitted, task_uuid=%s, duration=%d, hardwareId=%d", taskUuid, duration, hardwareId)
	return transaction.Hash().String(), nil
}

func (c *APIClient) validatePayment(txHash, taskUuid string) (*ValidatePaymentResult, error) {
	var validatePaymentResult ValidatePaymentResult
	if txHash != "" && taskUuid != "" {
		var params map[string]string
		params["tx_hash"] = txHash
		params["task_uuid"] = taskUuid
		if err := c.httpClient.PostJSON(apiValidatePayment, params, NewResult(&validatePaymentResult)); err != nil {
			return nil, err
		}
		log.Printf("Payment validation request sent, task_uuid=%s, tx_hash=%s \n", taskUuid, txHash)
		return &validatePaymentResult, nil
	} else {
		return nil, fmt.Errorf("tx_hash or task_uuid invalid")
	}
}

func (c *APIClient) getContractInfo(validate bool) (ContractDetail, error) {
	var contractResult ContractResult
	if err := c.httpClient.Get(apiContract, nil, NewResult(&contractResult)); err != nil {
		return ContractDetail{}, err
	}

	if validate {
		if !contractInfoVerified(contractResult.ContractInfo, contractResult.Signature, OrchestratorPublicAddressTestnet) {
			return ContractDetail{}, fmt.Errorf("failed to verified contract")
		}
	}
	return contractResult.ContractInfo.ContractDetail, nil
}

func (c *APIClient) getHardwareByInstanceType(instanceType string) (HardwareBaseInfo, error) {
	var baseInfo HardwareBaseInfo

	if strings.TrimSpace(instanceType) == "" {
		return baseInfo, fmt.Errorf("invalid instanceType")
	}

	hardwares, err := c.Hardwares()
	if err != nil {
		return baseInfo, fmt.Errorf("failed to get hardware, error: %v", err)
	}

	for _, hardware := range hardwares {
		if hardware.Name == instanceType {
			baseInfo = hardware.HardwareBaseInfo
			break
		}
	}

	if baseInfo.Description != "" {
		return baseInfo, nil
	}
	return baseInfo, fmt.Errorf("invalid instanceType: %s", instanceType)
}

type Result struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func NewResult(dest any) *Result {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		panic("dest must be a pointer")
	}
	var result Result
	result.Data = dest
	return &result
}

func (r *Result) Check() error {
	if r.Status != "success" {
		return errors.New(r.Message)
	}
	return nil
}
