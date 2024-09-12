package swan

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/swanchain/go-swan-sdk/contract"
)

type APIClient struct {
	apiKey         string
	httpClient     *HttpClient
	contractDetail ContractDetail
	instanceList   []*InstanceResource
}

func NewAPIClient(apiKey string, isTestnet ...bool) (*APIClient, error) {
	host := gatewayMainnet
	if len(isTestnet) > 0 && isTestnet[0] {
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
		return nil, fmt.Errorf("failed to get contract detail, error: %v", err)
	}

	instanceList, err := apiClient.InstanceResources()
	if err != nil {
		return nil, fmt.Errorf("failed to get instance resources, error: %v", err)
	}

	chainId, netWorkName, err := getNetWorkInfo(contractDetail.RpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain info, error: %v", err)
	}
	fmt.Printf("Logging in Swan Chain %s(%d) \n", netWorkName, chainId)

	apiClient.instanceList = instanceList
	apiClient.contractDetail = contractDetail
	return &apiClient, nil
}

func (c *APIClient) login() error {
	var token string
	if err := c.httpClient.PostForm(apiLogin, url.Values{"api_key": {c.apiKey}}, NewResult(&token)); err != nil {
		return err
	}
	c.httpClient.header.Set("Authorization", "Bearer "+token)
	return nil
}

func (c *APIClient) InstanceResources(available ...bool) ([]*InstanceResource, error) {
	var result InstanceResult

	if err := c.httpClient.Get(apiMachines, nil, NewResult(&result)); err != nil {
		return nil, err
	}

	if len(available) > 0 && available[0] {
		var data []*InstanceResource
		for _, instance := range result.Instances {
			instanceCp := instance
			if instanceCp.Status == "available" {
				data = append(data, instanceCp)
			}
		}
		return data, nil
	}
	return result.Instances, nil
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

func (c *APIClient) CreateTask(req *CreateTaskReq) (*CreateTaskResp, error) {
	var createTaskResp CreateTaskResp

	if req.WalletAddress == "" && req.PrivateKey == "" {
		return nil, fmt.Errorf("please provide WalletAddress or PrivateKey")
	}

	var walletAddress = req.WalletAddress
	if walletAddress != "" {
		walletAddress = req.WalletAddress
	} else {
		publicKeyAddress, err := privateKeyToPublicKey(req.PrivateKey)
		if err != nil {
			return nil, err
		}
		walletAddress = publicKeyAddress.String()
	}

	if req.Region == "" {
		req.Region = "global"
	}

	if req.StartIn == 0 {
		req.StartIn = 300
	}

	if req.Duration.Seconds() < 3600 {
		return nil, fmt.Errorf("duration must be no less than 3600 seconds")
	}

	if strings.TrimSpace(req.InstanceType) == "" {
		req.InstanceType = "C1ae.small"
	}

	if _, err := c.getInstanceByInstanceType(req.InstanceType); err != nil {
		return nil, err
	}
	log.Printf("Using %s machine, region=%s  duration=%d (seconds) \n", req.InstanceType, req.Region, req.Duration)

	if req.JobSourceUri == "" {
		if req.RepoUri != "" {
			sourceUri, err := c.getSourceUri(req.RepoUri, walletAddress, req.InstanceType, req.RepoBranch)
			if err != nil {
				return nil, fmt.Errorf("please provide JobSourceUri, or RepoUri, error: %v", err)
			}
			req.JobSourceUri = sourceUri
		}
	}

	if req.JobSourceUri == "" {
		return nil, fmt.Errorf("cannot get JobSourceUri. make sure `RepoUri` or `JobSourceUri` is correct")
	}

	var preferredCp string
	if len(req.PreferredCpList) > 0 {
		preferredCp = strings.Join(req.PreferredCpList, ",")
	}

	if !c.verifyHardwareRegion(req.InstanceType, req.Region) {
		return nil, fmt.Errorf("no %s machine in %s", req.InstanceType, req.Region)
	}
	var params = make(url.Values)
	params.Set("duration", strconv.Itoa(int(req.Duration.Seconds())))
	params.Set("cfg_name", req.InstanceType)
	params.Set("region", req.Region)
	params.Set("start_in", strconv.Itoa(req.StartIn))
	params.Set("wallet", walletAddress)
	params.Set("job_source_uri", req.JobSourceUri)
	if preferredCp != "" {
		params.Add("preferred_cp", preferredCp)
	}

	if err := c.login(); err != nil {
		return nil, err
	}
	if err := c.httpClient.PostForm(apiTask, params, NewResult(&createTaskResp)); err != nil {
		return nil, fmt.Errorf("failed to create task, error: %v", err)
	}

	taskUuid := createTaskResp.Task.UUID
	createTaskResp.TaskUuid = taskUuid
	createTaskResp.InstanceType = req.InstanceType

	estimatePrice, err := c.EstimatePayment(req.InstanceType, req.Duration.Seconds())
	if err != nil {
		return nil, err
	}
	createTaskResp.Price = estimatePrice

	var txHash string
	if req.PrivateKey != "" {
		payment, err := c.PayAndDeployTask(taskUuid, req.PrivateKey, req.Duration, req.InstanceType)
		if err != nil {
			return nil, err
		}
		createTaskResp.ConfigOrder = payment.ConfigOrder
		createTaskResp.TxHash = payment.TxHash
		createTaskResp.ApproveHash = payment.ApproveHash
		log.Printf("Task created successfully, taskUuid=%s, txHash=%s, instanceType=%s", taskUuid, txHash, req.InstanceType)
	}
	return &createTaskResp, nil
}

func (c *APIClient) PayAndDeployTask(taskUuid, privateKey string, duration time.Duration, instanceType string) (*PaymentResult, error) {
	var paymentResult PaymentResult

	if strings.TrimSpace(instanceType) == "" {
		return nil, fmt.Errorf("invalid instanceType")
	}
	if privateKey == "" {
		return nil, fmt.Errorf("no privateKey provided")
	}

	approveTxHash, submitPaymentTx, err := c.submitPayment(taskUuid, privateKey, duration, instanceType)
	if err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)
	validatePaymentResult, err := c.validatePayment(submitPaymentTx, taskUuid)
	if err != nil {
		return nil, err
	}
	paymentResult.ConfigOrder = validatePaymentResult.ConfigOrder
	paymentResult.ApproveHash = approveTxHash
	log.Printf("Payment submitted and validated successfully, taskUuid=%s, approve_hash=%s, tx_hash=%s", taskUuid, approveTxHash, submitPaymentTx)

	return &paymentResult, nil

}

func (c *APIClient) EstimatePayment(instanceType string, duration float64) (float64, error) {
	hardwareBaseInfo, err := c.getInstanceByInstanceType(instanceType)
	if err != nil {
		return 0, err
	}

	priceInt, err := strconv.ParseFloat(hardwareBaseInfo.Price, 64)
	if err != nil {
		return 0, err
	}
	return priceInt * (duration / 3600), nil
}

func (c *APIClient) RenewTask(taskUuid string, duration time.Duration, privateKey string, paidTxHash ...string) (*RenewTaskResp, error) {
	if strings.TrimSpace(taskUuid) == "" {
		return nil, fmt.Errorf("invalid taskUuid")
	}

	if privateKey == "" && len(paidTxHash) == 0 {
		return nil, fmt.Errorf("provide a txHash or privateKey")
	}

	var txHash string
	if len(paidTxHash) == 0 {
		reNewPaymentTxHash, err := c.RenewPayment(taskUuid, duration, privateKey)
		if err != nil {
			return nil, err
		}
		txHash = reNewPaymentTxHash
	} else {
		txHash = paidTxHash[0]
		log.Printf("Using given payment transaction hash, txHash=%s", txHash)
	}

	if txHash != "" && taskUuid != "" {
		var params = make(url.Values)
		params.Set("task_uuid", taskUuid)
		params.Set("duration", strconv.FormatFloat(duration.Seconds(), 'f', 2, 64))
		params.Set("tx_hash", txHash)

		var renewTaskResp RenewTaskResp
		if err := c.httpClient.PostForm(apiReNewTask, params, NewResult(&renewTaskResp)); err != nil {
			return nil, err
		}
		return &renewTaskResp, nil
	} else {
		return nil, fmt.Errorf("txHash or taskUuid invalid")
	}
}

func (c *APIClient) RenewPayment(taskUuid string, duration time.Duration, privateKey string) (string, error) {
	if strings.TrimSpace(taskUuid) == "" {
		return "", fmt.Errorf("invalid taskUuid")
	}
	if privateKey == "" {
		return "", fmt.Errorf("no privateKey provided")
	}

	taskInfo, err := c.TaskInfo(taskUuid)
	if err != nil {
		return "", fmt.Errorf("failed to get task info, taskUuid: %s, error: %v", taskUuid, err)
	}
	instanceType := taskInfo.Task.TaskDetail.Hardware
	hardwareBaseInfo, err := c.getInstanceByInstanceType(instanceType)
	if err != nil {
		return "", err
	}
	var hardwareId = hardwareBaseInfo.ID

	estimatePrice, err := c.EstimatePayment(instanceType, duration.Seconds())
	if err != nil {
		return "", err
	}
	priceBigInt, ok := new(big.Int).SetString(fmt.Sprintf("%.f", estimatePrice), 10)
	if !ok {
		return "", fmt.Errorf("failed to convert float64 to big.Int")
	}

	client, err := ethclient.Dial(c.contractDetail.RpcUrl)
	if err != nil {
		return "", err
	}
	defer client.Close()

	tokenContract, err := contract.NewToken(common.HexToAddress(c.contractDetail.SwanTokenContractAddress), client)
	if err != nil {
		return "", err
	}

	paymentContract, err := contract.NewPaymentContract(common.HexToAddress(c.contractDetail.ClientContractAddress), client)
	if err != nil {
		return "", err
	}

	approveTransactOpts, err := CreateTransactOpts(client, privateKey)
	if err != nil {
		return "", err
	}
	approve, err := tokenContract.Approve(approveTransactOpts, common.HexToAddress(c.contractDetail.ClientContractAddress), priceBigInt)
	if err != nil {
		return "", err
	}

	tokenApproveHash := approve.Hash().String()
	timeout := time.After(1 * time.Minute)
	ticker := time.Tick(3 * time.Second)
	for {
		select {
		case <-timeout:
			return "", fmt.Errorf("timeout waiting for transaction confirmation, tx: %s", tokenApproveHash)
		case <-ticker:
			receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tokenApproveHash))
			if err != nil {
				if errors.Is(err, ethereum.NotFound) {
					continue
				}
				return "", fmt.Errorf("check swan token Approve tx, error: %+v", err)
			}

			if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
				fmt.Printf("swan token approve TX Hash: %s \n", tokenApproveHash)

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
				log.Printf("Payment submitted, task_uuid=%s, duration=%f, hardwareId=%d", taskUuid, duration.Seconds(), hardwareId)
				return transaction.Hash().String(), nil
			} else if receipt != nil && receipt.Status == 0 {
				return "", fmt.Errorf("failed to check swan token approve transaction, tx: %s", tokenApproveHash)
			}
		}
	}
}

func (c *APIClient) TerminateTask(taskUuid string) (TerminateTaskResp, error) {
	var terminateTaskResp TerminateTaskResp

	if strings.TrimSpace(taskUuid) == "" {
		return terminateTaskResp, fmt.Errorf("invalid taskUuid")
	}

	var params = make(url.Values)
	params.Set("task_uuid", taskUuid)
	if err := c.httpClient.PostForm(apiTerminateTask, params, NewResult(&terminateTaskResp)); err != nil {
		return terminateTaskResp, err
	}
	return terminateTaskResp, nil
}

func (c *APIClient) GetRealUrl(taskUuid string) ([]string, error) {
	if strings.TrimSpace(taskUuid) == "" {
		return nil, fmt.Errorf("invalid taskUuid")
	}

	taskInfo, err := c.TaskInfo(taskUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get task info, error: %v", err)
	}

	var deployedUrl []string
	for _, job := range taskInfo.Jobs {
		if strings.TrimSpace(job.JobRealURI) != "" {
			deployedUrl = append(deployedUrl, job.JobRealURI)
		}
	}
	return deployedUrl, nil
}

func (c *APIClient) getSourceUri(repoUri, walletAddress string, instanceType string, repoBranch string) (string, error) {
	var jobSourceUriResult JobSourceUriResult

	hardwareBaseInfo, err := c.getInstanceByInstanceType(instanceType)
	if err != nil {
		return "", err
	}

	if walletAddress == "" {
		return "", fmt.Errorf("no wallet_address provided")
	}
	var reqData = make(url.Values)
	reqData.Set("repo_branch", repoBranch)
	reqData.Set("wallet_address", walletAddress)
	reqData.Set("hardware_id", strconv.FormatInt(hardwareBaseInfo.ID, 10))
	reqData.Set("repo_uri", repoUri)

	if err := c.httpClient.PostForm(apiSourceUri, reqData, NewResult(&jobSourceUriResult)); err != nil {
		return "", err
	}
	return jobSourceUriResult.JobSourceUri, nil
}

func (c *APIClient) verifyHardwareRegion(instanceType, region string) bool {
	for _, hardware := range c.instanceList {
		if hardware.Type == instanceType {
			for _, r := range hardware.Region {
				if region == r || (strings.ToLower(region) == "global" && hardware.Status == "available") {
					return true
				}
			}
		}
	}
	return false
}

func (c *APIClient) submitPayment(taskUuid, privateKey string, duration time.Duration, instanceType string) (string, string, error) {
	hardwareBaseInfo, err := c.getInstanceByInstanceType(instanceType)
	if err != nil {
		return "", "", err
	}

	if privateKey == "" {
		return "", "", fmt.Errorf("no privateKey provided")
	}
	var hardwareId = hardwareBaseInfo.ID
	estimatePrice, err := c.EstimatePayment(instanceType, duration.Seconds())
	if err != nil {
		return "", "", err
	}

	priceBigInt, ok := new(big.Int).SetString(fmt.Sprintf("%.f", estimatePrice), 10)
	if !ok {
		return "", "", fmt.Errorf("failed to convert float64 to big.Int")
	}

	client, err := ethclient.Dial(c.contractDetail.RpcUrl)
	if err != nil {
		return "", "", err
	}
	defer client.Close()

	tokenContract, err := contract.NewToken(common.HexToAddress(c.contractDetail.SwanTokenContractAddress), client)
	if err != nil {
		return "", "", err
	}

	paymentContract, err := contract.NewPaymentContract(common.HexToAddress(c.contractDetail.ClientContractAddress), client)
	if err != nil {
		return "", "", err
	}

	approveTransactOpts, err := CreateTransactOpts(client, privateKey)
	if err != nil {
		return "", "", err
	}

	hardwareIdBigInt := new(big.Int).SetInt64(hardwareId)
	durationBigInt := new(big.Int).SetInt64(int64(duration))
	approve, err := tokenContract.Approve(approveTransactOpts, common.HexToAddress(c.contractDetail.ClientContractAddress), priceBigInt)
	if err != nil {
		return "", "", err
	}

	tokenApproveHash := approve.Hash().String()
	timeout := time.After(1 * time.Minute)
	ticker := time.Tick(3 * time.Second)
	for {
		select {
		case <-timeout:
			return "", "", fmt.Errorf("timeout waiting for transaction confirmation, tx: %s", tokenApproveHash)
		case <-ticker:
			receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tokenApproveHash))
			if err != nil {
				if errors.Is(err, ethereum.NotFound) {
					continue
				}
				return "", "", fmt.Errorf("check swan token Approve tx, error: %+v", err)
			}

			if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
				fmt.Printf("swan token approve TX Hash: %s \n", tokenApproveHash)

				paymentTransactOpts, err := CreateTransactOpts(client, privateKey)
				if err != nil {
					return "", "", err
				}
				transaction, err := paymentContract.SubmitPayment(paymentTransactOpts, taskUuid, hardwareIdBigInt, durationBigInt)
				if err != nil {
					return "", "", fmt.Errorf("failed to submit payment, error: %v", err)
				}
				log.Printf("Payment submitted, task_uuid=%s, duration=%d, hardwareId=%d", taskUuid, duration, hardwareId)
				return tokenApproveHash, transaction.Hash().String(), nil
			} else if receipt != nil && receipt.Status == 0 {
				return "", "", fmt.Errorf("failed to check swan token approve transaction, tx: %s", tokenApproveHash)
			}
		}
	}
}

func (c *APIClient) validatePayment(txHash, taskUuid string) (*ValidatePaymentResult, error) {
	var validatePaymentResult ValidatePaymentResult
	if txHash != "" && taskUuid != "" {
		var params = make(url.Values)
		params.Set("tx_hash", txHash)
		params.Set("task_uuid", taskUuid)

		if err := c.httpClient.PostForm(apiValidatePayment, params, NewResult(&validatePaymentResult)); err != nil {
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

func (c *APIClient) getInstanceByInstanceType(instanceType string) (*InstanceBaseInfo, error) {
	var baseInfo InstanceBaseInfo

	if strings.TrimSpace(instanceType) == "" {
		return nil, fmt.Errorf("invalid instanceType")
	}

	for _, hardware := range c.instanceList {
		if hardware.Type == instanceType {
			baseInfo = hardware.InstanceBaseInfo
			break
		}
	}

	if baseInfo.Description == "" {
		return nil, fmt.Errorf("invalid instanceType: %s", instanceType)
	}
	return &baseInfo, nil
}

func getNetWorkInfo(rpc string) (int64, string, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return 0, "", err
	}
	defer client.Close()

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return 0, "", err
	}

	var network = "Testnet"
	if chainId.Int64() == 254 {
		network = "Mainnet"
	}
	return chainId.Int64(), network, nil
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
