package swan

import (
	"fmt"
	"testing"
	"time"
)

//# wallet_address = "0x3E364F3Ea3599329Cd9b6444Ac7947f0D73bAe75"
//# private_key = "543ed71922c1013ebaa25dc434e8470324a8ed8885f43578763fca79b1a02ccf"
//# api_key="HKScSBCAcj"

// ApiKey       = "53Qkrwdeyv"
// Wallet       = "0xFbc1d38a2127D81BFe3EA347bec7310a1cfa2373"
// PrivateKey   = "c3e47d07d520fd3022a4b61764cfcb831cdafc3352e97c21acb0138684c5d703"
const (
	ApiKey       = "HKScSBCAcj"
	Wallet       = "0x3E364F3Ea3599329Cd9b6444Ac7947f0D73bAe75"
	PrivateKey   = "543ed71922c1013ebaa25dc434e8470324a8ed8885f43578763fca79b1a02ccf"
	JobSourceUri = "https://test-api.lagrangedao.org/spaces/143a526d-0cfc-41d6-b95c-53a4018829c8"
)

func TestAPIClient_CreateTaskWithAutoPay(t *testing.T) {
	var req = CreateTaskReq{
		PrivateKey:   PrivateKey,
		JobSourceUri: JobSourceUri,
	}
	apiClient := NewAPIClient(ApiKey, true)
	resp, err := apiClient.CreateTask(&req)
	if err != nil {
		t.Errorf("CreateTaskWithAutoPay() error = %v", err)
	}
	t.Logf("create task with auto-pay response: %v", resp)
}

func TestAPIClient_CreateTask(t *testing.T) {
	var req = CreateTaskReq{
		JobSourceUri: JobSourceUri,
	}

	apiClient := NewAPIClient(ApiKey, true)
	resp, err := apiClient.CreateTask(&req)
	if err != nil {
		t.Errorf("CreateTask() error = %v", err)
	}
	t.Logf("create task response: %v", resp)
}

func TestAPIClient_PayAndDeployTask(t *testing.T) {
	apiClient := NewAPIClient(ApiKey, true)

	// taskUuid:  returned by create task
	taskUuid := "a9d2f2ca-8819-43f7-9347-7ccf0ea11822"
	payAndDeployTaskResp, err := apiClient.PayAndDeployTask(taskUuid, PrivateKey, time.Duration(3600), "C1ae.small")
	if err != nil {
		t.Errorf("PayAndDeployTask() error = %v", err)
	}
	t.Logf("pay and deploy task response: %v", payAndDeployTaskResp)
}

func TestAPIClient_TaskInfo(t *testing.T) {
	apiClient := NewAPIClient(ApiKey, true)
	resp, err := apiClient.TaskInfo("a9d2f2ca-8819-43f7-9347-7ccf0ea11822")
	if err != nil {
		t.Errorf("TaskInfo() error = %v", err)
	}
	t.Logf("get task info response: %v", resp)
}

func TestAPIClient_Tasks(t *testing.T) {
	var req = &TaskQueryReq{
		Wallet: Wallet,
		Page:   0,
		Size:   10,
	}
	apiClient := NewAPIClient(ApiKey, true)
	total, resp, err := apiClient.Tasks(req)
	if err != nil {
		t.Errorf("Tasks() error = %v", err)
	}
	t.Logf("get task list response, total: %d, data: %v", total, resp)
}

func TestAPIClient_Hardwares(t *testing.T) {
	apiClient := NewAPIClient(ApiKey)
	resp, err := apiClient.Hardwares()
	if err != nil {
		t.Errorf("Hardwares() error = %v", err)
	}

	for _, hardware := range resp {
		fmt.Printf("%s,%s,%s,%s\n", hardware.Name, hardware.Type, hardware.Description, hardware.Price)
	}
}

func TestAPIClient_GetRealUrl(t *testing.T) {
	apiClient := NewAPIClient(ApiKey, true)
	resp, err := apiClient.GetRealUrl("01c0f29a-2304-45fa-bb03-976348e714e4")
	if err != nil {
		t.Errorf("CreateTask() error = %v", err)
	}
	t.Logf("create task response: %v", resp)
}

func TestAPIClient_TerminateTask(t *testing.T) {
	apiClient := NewAPIClient(ApiKey, true)
	resp, err := apiClient.TerminateTask("01c0f29a-2304-45fa-bb03-976348e714e4")
	if err != nil {
		t.Errorf("TerminateTask() error = %v", err)
	}
	t.Logf("terminate task response: %v", resp)
}

func Test_getContractInfo(t *testing.T) {

	apiClient := NewAPIClient(ApiKey, true)
	resp, err := apiClient.getContractInfo(false)
	if err != nil {
		t.Errorf("getContractInfo() error = %v", err)
	}
	t.Logf("get contract response: %v", resp)
}

func TestAPIClient_EstimatePayment(t *testing.T) {
	apiClient := NewAPIClient(ApiKey, true)
	resp, err := apiClient.EstimatePayment("P1ae.medium", 3600)
	if err != nil {
		t.Errorf("EstimatePayment() error = %v", err)
	}
	// 32
	t.Logf("estimate Payment response: %v", resp)
}

func TestAPIClient_ReNewTask(t *testing.T) {
	apiClient := NewAPIClient(ApiKey, true)
	resp, err := apiClient.RenewTask("a9d2f2ca-8819-43f7-9347-7ccf0ea11822", time.Duration(3600), PrivateKey, "")
	if err != nil {
		t.Errorf("ReNewTask() error = %v", err)
	}
	t.Logf("renew task with auto-pay response: %v", resp)
}

func TestAPIClient_RenewPayment(t *testing.T) {
	apiClient := NewAPIClient(ApiKey, true)
	txHash, err := apiClient.RenewPayment("a9d2f2ca-8819-43f7-9347-7ccf0ea11822", time.Duration(3600), PrivateKey)
	if err != nil {
		t.Errorf("RenewPayment() error = %v", err)
	}
	t.Logf("renew payment response: %v", txHash)
}
