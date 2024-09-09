package swan

import (
	"testing"
	"time"
)

const (
	ApiKey        = "<API_KEY>"
	WalletAddress = "<WALLET_ADDRESS>"
	PrivateKey    = "<WALLET_ADDRESS_PRIVATE_KEY>"
	JobSourceUri  = "https://test-api.lagrangedao.org/spaces/143a526d-0cfc-41d6-b95c-53a4018829c8"
)

func TestAPIClient_CreateTaskWithAutoPay(t *testing.T) {
	var req = CreateTaskReq{
		Duration:     time.Duration(3600),
		PrivateKey:   PrivateKey,
		JobSourceUri: JobSourceUri,
	}
	apiClient, _ := NewAPIClient(ApiKey, true)
	resp, err := apiClient.CreateTask(&req)
	if err != nil {
		t.Errorf("CreateTaskWithAutoPay() error = %v", err)
	}
	t.Logf("create task with auto-pay response: %v", resp)
}

func TestAPIClient_CreateTask(t *testing.T) {
	var req = CreateTaskReq{
		JobSourceUri:  JobSourceUri,
		WalletAddress: WalletAddress,
	}

	apiClient, _ := NewAPIClient(ApiKey, true)
	resp, err := apiClient.CreateTask(&req)
	if err != nil {
		t.Errorf("CreateTask() error = %v", err)
	}
	t.Logf("create task response: %v", resp)
}

func TestAPIClient_PayAndDeployTask(t *testing.T) {
	apiClient, _ := NewAPIClient(ApiKey, true)

	// taskUuid:  returned by create task
	taskUuid := "a9d2f2ca-8819-43f7-9347-7ccf0ea11822"
	payAndDeployTaskResp, err := apiClient.PayAndDeployTask(taskUuid, PrivateKey, time.Duration(3600), "C1ae.small")
	if err != nil {
		t.Errorf("PayAndDeployTask() error = %v", err)
	}
	t.Logf("pay and deploy task response: %v", payAndDeployTaskResp)
}

func TestAPIClient_TaskInfo(t *testing.T) {
	apiClient, _ := NewAPIClient(ApiKey, true)
	resp, err := apiClient.TaskInfo("a9d2f2ca-8819-43f7-9347-7ccf0ea11822")
	if err != nil {
		t.Errorf("TaskInfo() error = %v", err)
	}
	t.Logf("get task info response: %v", resp)
}

func TestAPIClient_Tasks(t *testing.T) {
	var req = &TaskQueryReq{
		Wallet: WalletAddress,
		Page:   0,
		Size:   10,
	}
	apiClient, _ := NewAPIClient(ApiKey, true)
	total, resp, err := apiClient.Tasks(req)
	if err != nil {
		t.Errorf("Tasks() error = %v", err)
	}
	t.Logf("get task list response, total: %d, data: %v", total, resp)
}

func TestAPIClient_Hardwares(t *testing.T) {
	apiClient, _ := NewAPIClient(ApiKey, true)
	resp, err := apiClient.Hardwares()
	if err != nil {
		t.Errorf("Hardwares() error = %v", err)
	}
	t.Logf("get hardware list response: %v", resp)
}

func TestAPIClient_GetRealUrl(t *testing.T) {
	apiClient, _ := NewAPIClient(ApiKey, true)
	resp, err := apiClient.GetRealUrl("01c0f29a-2304-45fa-bb03-976348e714e4")
	if err != nil {
		t.Errorf("CreateTask() error = %v", err)
	}
	t.Logf("create task response: %v", resp)
}

func TestAPIClient_TerminateTask(t *testing.T) {
	apiClient, _ := NewAPIClient(ApiKey, true)
	resp, err := apiClient.TerminateTask("01c0f29a-2304-45fa-bb03-976348e714e4")
	if err != nil {
		t.Errorf("TerminateTask() error = %v", err)
	}
	t.Logf("terminate task response: %v", resp)
}

func Test_getContractInfo(t *testing.T) {

	apiClient, _ := NewAPIClient(ApiKey, true)
	resp, err := apiClient.getContractInfo(false)
	if err != nil {
		t.Errorf("getContractInfo() error = %v", err)
	}
	t.Logf("get contract response: %v", resp)
}

func TestAPIClient_EstimatePayment(t *testing.T) {
	apiClient, _ := NewAPIClient(ApiKey, true)
	resp, err := apiClient.EstimatePayment("P1ae.medium", 3600)
	if err != nil {
		t.Errorf("EstimatePayment() error = %v", err)
	}
	// 32
	t.Logf("estimate Payment response: %v", resp)
}

func TestAPIClient_ReNewTask(t *testing.T) {
	apiClient, _ := NewAPIClient(ApiKey, true)
	resp, err := apiClient.RenewTask("a9d2f2ca-8819-43f7-9347-7ccf0ea11822", time.Duration(3600), PrivateKey, "")
	if err != nil {
		t.Errorf("ReNewTask() error = %v", err)
	}
	t.Logf("renew task with auto-pay response: %v", resp)
}

func TestAPIClient_RenewPayment(t *testing.T) {
	apiClient, _ := NewAPIClient(ApiKey, true)
	txHash, err := apiClient.RenewPayment("a9d2f2ca-8819-43f7-9347-7ccf0ea11822", time.Duration(3600), PrivateKey)
	if err != nil {
		t.Errorf("RenewPayment() error = %v", err)
	}
	t.Logf("renew payment response: %v", txHash)
}
