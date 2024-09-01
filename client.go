package swan

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/swanchain/go-swan-sdk/contract"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type APIClient struct {
	apiKey     string
	httpClient *HttpClient
}

func NewAPIClient(apiKey string, testnet ...bool) *APIClient {
	host := gatewayMainnet
	if len(testnet) > 0 && testnet[0] {
		host = gatewayTestnet
	}

	header := make(http.Header)
	header.Add("Authorization", "Bearer "+apiKey)

	return &APIClient{
		apiKey:     apiKey,
		httpClient: NewHttpClient(host, header),
	}
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

	hardwares, err := c.Hardwares()
	if err != nil {
		return createTaskResp, err
	}

	var cfgName string
	for _, hardware := range hardwares {
		if hardware.ID == req.HardwareId {
			cfgName = hardware.Name
			break
		}
	}

	if cfgName == "" {
		return createTaskResp, fmt.Errorf("invalid HardwareId selected")
	}

	log.Printf("Using %s machine, %d %s %d (seconds) \n", cfgName, req.HardwareId, req.Region, req.Duration)

	if req.JobSourceUri == "" {
		if req.AppRepoImage != "" {
			if !req.AutoPay && req.PrivateKey != "" {
				req.AutoPay = true
			}
			appRepoImage, err := c.GetAppRepoImage(req.AppRepoImage)
			if err != nil {
				return createTaskResp, err
			}
			if appRepoImage.Url == "" {
				return createTaskResp, fmt.Errorf("Invalid appRepoImage url ")
			}
			req.RepoUri = appRepoImage.Url
		}
		if req.RepoUri != "" {
			sourceUri, err := c.GetSourceUri(req.RepoUri, req.WalletAddress, &req.HardwareId, req.RepoBranch, req.RepoOwner, req.RepoName)
			if err != nil {
				return createTaskResp, fmt.Errorf("please provide AppRepoImage, or JobSourceUri, or RepoUri, error: %v", err)
			}
			req.JobSourceUri = sourceUri
		}
	}

	if req.JobSourceUri == "" {
		return createTaskResp, fmt.Errorf("cannot get JobSourceUri. make sure `AppRepoImage` or `RepoUri` or `JobSourceUri` is correct")
	}

	var preferredCp string
	if len(req.PreferredCpList) > 0 {
		preferredCp = strings.Join(req.PreferredCpList, ",")
	}

	var params map[string]interface{}
	if c.VerifyHardwareRegion(cfgName, req.Region) {
		params["duration"] = req.Duration
		params["cfg_name"] = cfgName
		params["region"] = req.Region
		params["start_in"] = req.StartIn
		params["wallet"] = req.WalletAddress
		params["job_source_uri"] = req.JobSourceUri
	} else {
		return createTaskResp, fmt.Errorf("no %s machine in %s", cfgName, req.Region)
	}

	if preferredCp != "" {
		params["preferred_cp"] = preferredCp
	}

	var task Task
	if err := c.httpClient.PostJSON(apiTask, params, NewResult(&task)); err != nil {
		return createTaskResp, fmt.Errorf("failed to create task, error: %v", err)
	}

	taskUuid := task.UUID
	var txHash string

	if req.AutoPay {
		payment, err := c.MakePayment(taskUuid, req.PrivateKey, req.Duration, &req.HardwareId)
		if err != nil {
			return createTaskResp, err
		}
		createTaskResp.ConfigOrder = payment.ConfigOrder
		createTaskResp.TxHash = payment.TxHash
		createTaskResp.TaskUuid = taskUuid
		createTaskResp.HardwareId = req.HardwareId
		createTaskResp.Id = taskUuid
		log.Printf("Task created successfully, taskUuid=%s, txHash=%s, hardwareId=%d", taskUuid, txHash, req.HardwareId)
		return createTaskResp, nil
	}
	return createTaskResp, nil
}

func (c *APIClient) MakePayment(taskUuid, privateKey string, duration int, hardwareId *int64) (PaymentResult, error) {
	var paymentResult PaymentResult
	if hardwareId == nil {
		return paymentResult, fmt.Errorf("invalid hardwareId")
	}
	if privateKey == "" {
		return paymentResult, fmt.Errorf("no privateKey provided")
	}

	submitPaymentTx, err := c.SubmitPayment(taskUuid, privateKey, duration, hardwareId)
	if err != nil {
		return paymentResult, err
	}
	time.Sleep(3 * time.Second)
	_, err = c.ValidatePayment(submitPaymentTx, taskUuid)
	if err != nil {
		return paymentResult, err
	}
	log.Printf("Payment submitted and validated successfully, taskUuid=%s, tx_hash=%s", taskUuid, submitPaymentTx)

	return paymentResult, nil

}

// SubmitPayment Submit payment for a task
//
// Args:
// task_uuid: unique id returned by `swan_api.create_task`
// hardware_id: id of cp/hardware configuration set
// duration: duration of service runtime (seconds).
//
// Returns:
// tx_hash
func (c *APIClient) SubmitPayment(taskUuid, privateKey string, duration int, hardwareId *int64) (string, error) {
	if hardwareId == nil {
		return "", fmt.Errorf("invalid hardwareId")
	}
	if privateKey == "" {
		return "", fmt.Errorf("no privateKey provided")
	}

	client, err := ethclient.Dial("")
	if err != nil {
		return "", err
	}
	defer client.Close()

	hardwareIdBigInt := new(big.Int).SetInt64(*hardwareId)
	durationBigInt := new(big.Int).SetInt64(int64(duration))

	// call token contract approve
	tokenContract, err := contract.NewToken(common.HexToAddress(""), client)
	if err != nil {
		return "", err
	}

	// call payment contract submit_payment
	paymentContract, err := contract.NewPaymentContract(common.HexToAddress(""), client)
	if err != nil {
		return "", err
	}

	hardwareInfo, err := paymentContract.HardwareInfo(&bind.CallOpts{}, hardwareIdBigInt)
	if err != nil {
		return "", err
	}
	price := hardwareInfo.PricePerHour.Int64() * int64(duration/3600)
	priceBigInt := new(big.Int).SetInt64(price)

	// client contract address
	approveTransactOpts, err := CreateTransactOpts(client, privateKey)
	if err != nil {
		return "", err
	}
	approve, err := tokenContract.Approve(approveTransactOpts, common.HexToAddress("client_contract_address"), priceBigInt)
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

func (c *APIClient) ValidatePayment(txHash, taskUuid string) (*ValidatePaymentResult, error) {
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

func (c *APIClient) VerifyHardwareRegion(cfgName, region string) bool {
	hardwares, err := c.Hardwares()
	if err != nil {
		log.Printf("failed to get hardware, error: %v", err)
		return false
	}

	for _, hardware := range hardwares {
		if hardware.Name == cfgName {
			for _, r := range hardware.Region {
				if region == r || (strings.ToLower(region) == "global" && hardware.Status == "available") {
					return true
				}
			}
		}
	}
	return false
}

func (c *APIClient) GetAppRepoImage(name ...string) (*RepoImageResult, error) {
	var repoImageResult RepoImageResult

	if len(name) > 0 {
		v := url.Values{}
		v.Add("name", name[0])
		if err := c.httpClient.Get(apiPremadeImage, v, NewResult(&repoImageResult)); err != nil {
			return nil, err
		}
		return &repoImageResult, nil
	}
	if err := c.httpClient.Get(apiPremadeImage, nil, NewResult(&repoImageResult)); err != nil {
		return nil, err
	}
	return &repoImageResult, nil
}

func (c *APIClient) GetSourceUri(repoUri, walletAddress string, hardwareId *int64, repoBranch, repoOwner, repoName string) (string, error) {
	var jobSourceUriResult JobSourceUriResult
	if hardwareId == nil {
		return "", fmt.Errorf("no hardware_id provided")
	}
	if walletAddress == "" {
		return "", fmt.Errorf("no wallet_address provided")
	}
	var reqData = map[string]interface{}{
		"repo_owner":     repoOwner,
		"repo_name":      repoName,
		"repo_branch":    repoBranch,
		"wallet_address": walletAddress,
		"hardware_id":    hardwareId,
		"repo_uri":       repoUri,
	}

	if err := c.httpClient.PostJSON(apiSourceUri, reqData, NewResult(&jobSourceUriResult)); err != nil {
		return "", err
	}
	return jobSourceUriResult.JobSourceUri, nil
}

func (c *APIClient) ReNewTask(taskUuid, txHash, privateKey string, hardwareId *int64, duration int, autoPay bool) (*ReNewTaskResp, error) {
	if hardwareId == nil {
		return nil, fmt.Errorf("invalid hardwareId")
	}

	if !autoPay && privateKey == "" && txHash == "" {
		return nil, fmt.Errorf("auto_pay off or tx_hash not provided, please provide a tx_hash or set auto_pay to True and provide private_key")
	}
	if txHash != "" {
		// renew_payment
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

func (c *APIClient) RenewPayment(taskUuid, privateKey string, hardwareId *int64, duration int) (*ReNewTaskResp, error) {
	if hardwareId == nil {
		return nil, fmt.Errorf("invalid hardwareId")
	}
	if privateKey == "" {
		return nil, fmt.Errorf("no privateKey provided")
	}

	client, err := ethclient.Dial("")
	if err != nil {
		return nil, err
	}
	defer client.Close()

	hardwareIdBigInt := new(big.Int).SetInt64(*hardwareId)
	durationBigInt := new(big.Int).SetInt64(int64(duration))

	// call token contract approve
	tokenContract, err := contract.NewToken(common.HexToAddress(""), client)
	if err != nil {
		return nil, err
	}

	// call payment contract submit_payment
	paymentContract, err := contract.NewPaymentContract(common.HexToAddress(""), client)
	if err != nil {
		return nil, err
	}

	hardwareInfo, err := paymentContract.HardwareInfo(&bind.CallOpts{}, hardwareIdBigInt)
	if err != nil {
		return nil, err
	}
	price := hardwareInfo.PricePerHour.Int64() * int64(duration/3600)
	priceBigInt := new(big.Int).SetInt64(price)

	// client contract address
	approveTransactOpts, err := CreateTransactOpts(client, privateKey)
	if err != nil {
		return nil, err
	}
	approve, err := tokenContract.Approve(approveTransactOpts, common.HexToAddress("client_contract_address"), priceBigInt)
	if err != nil {
		return nil, err
	}
	log.Printf("token approve tx: %v \n", approve.Hash().String())

	paymentTransactOpts, err := CreateTransactOpts(client, privateKey)
	if err != nil {
		return nil, err
	}
	transaction, err := paymentContract.SubmitPayment(paymentTransactOpts, taskUuid, hardwareIdBigInt, durationBigInt)
	if err != nil {
		return nil, fmt.Errorf("failed to submit payment, error: %v", err)
	}
	log.Printf("txHash: %v", transaction.Hash().String())
	return nil, nil
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
