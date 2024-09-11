package swan

import "time"

type InstanceResult struct {
	Instances []*InstanceResource `json:"hardware"`
}

type InstanceResource struct {
	InstanceBaseInfo
	Type          string                   `json:"hardware_type"`
	Region        []string                 `json:"region"`
	RegionDetails map[string]*RegionDetail `json:"region_detail"`
}

type InstanceBaseInfo struct {
	Description string `json:"hardware_description"`
	ID          int64  `json:"hardware_id"`
	Name        string `json:"hardware_type"`
	Price       string `json:"hardware_price"`
	Status      string `json:"hardware_status"`
	Type        string `json:"hardware_name"`
}

type RegionDetail struct {
	AvailableResource int64      `json:"available_resource"`
	DirectAccessCp    [][]string `json:"direct_access_cp"`
	NoneCollateral    int64      `json:"none_collateral"`
	Whitelist         int64      `json:"whitelist"`
}

type TaskInfo struct {
	Providers []*ComputingProvider `json:"computing_providers"`
	Orders    []*ConfigOrder       `json:"config_orders"`
	Jobs      []*Job               `json:"jobs"`
	Task      Task                 `json:"task"`
}

type ComputingProvider struct {
	Beneficiary      string   `json:"beneficiary"`
	CpAccountAddress string   `json:"cp_account_address"`
	CreatedAt        int64    `json:"created_at"`
	FreezeOnline     any      `json:"freeze_online"`
	ID               int64    `json:"id"`
	Lat              float64  `json:"lat"`
	Lon              float64  `json:"lon"`
	MultiAddress     []string `json:"multi_address"`
	Name             string   `json:"name"`
	NodeID           string   `json:"node_id"`
	Online           int      `json:"online"`
	OwnerAddress     string   `json:"owner_address"`
	Region           string   `json:"region"`
	TaskTypes        string   `json:"task_types"`
	UpdatedAt        int64    `json:"updated_at"`
	Version          string   `json:"version"`
	WorkerAddress    string   `json:"worker_address"`
}

type Job struct {
	BuildLog         string `json:"build_log"`
	Comments         string `json:"comments"`
	ContainerLog     string `json:"container_log"`
	CpAccountAddress string `json:"cp_account_address"`
	CreatedAt        int64  `json:"created_at"`
	Duration         int64  `json:"duration"`
	EndedAt          int64  `json:"ended_at"`
	Hardware         string `json:"hardware"`
	ID               int64  `json:"id"`
	JobRealURI       string `json:"job_real_uri"`
	JobResultURI     string `json:"job_result_uri"`
	JobSourceURI     string `json:"job_source_uri"`
	Name             string `json:"name"`
	NodeID           string `json:"node_id"`
	StartAt          int64  `json:"start_at"`
	Status           string `json:"status"`
	StorageSource    string `json:"storage_source"`
	TaskUUID         string `json:"task_uuid"`
	Type             any    `json:"type"`
	UpdatedAt        int64  `json:"updated_at"`
	UUID             string `json:"uuid"`
}

type Task struct {
	Comments      string      `json:"comments"`
	CreatedAt     int64       `json:"created_at"`
	EndAt         int64       `json:"end_at"`
	ID            int64       `json:"id"`
	LeadingJobID  string      `json:"leading_job_id"`
	Name          string      `json:"name"`
	RefundAmount  string      `json:"refund_amount"`
	RefundWallet  string      `json:"refund_wallet"`
	Source        string      `json:"source"`
	StartAt       int64       `json:"start_at"`
	StartIn       int64       `json:"start_in"`
	Status        string      `json:"status"`
	TaskDetail    *TaskDetail `json:"task_detail"`
	TaskDetailCid string      `json:"task_detail_cid"`
	TxHash        any         `json:"tx_hash"`
	Type          string      `json:"type"`
	UpdatedAt     int64       `json:"updated_at"`
	UserID        int64       `json:"user_id"`
	UUID          string      `json:"uuid"`
}

type TaskDetail struct {
	Amount            float64       `json:"amount"`
	BidderLimit       int64         `json:"bidder_limit"`
	CreatedAt         int64         `json:"created_at"`
	DCCSelectedCpList any           `json:"dcc_selected_cp_list"`
	Duration          int64         `json:"duration"`
	EndAt             int64         `json:"end_at"`
	Hardware          string        `json:"hardware"`
	JobResultURI      string        `json:"job_result_uri"`
	JobSourceURI      string        `json:"job_source_uri"`
	PricePerHour      string        `json:"price_per_hour"`
	Requirements      *Requirements `json:"requirements"`
	Space             *Space        `json:"space"`
	StartAt           int64         `json:"start_at"`
	Status            string        `json:"status"`
	StorageSource     string        `json:"storage_source"`
	Type              string        `json:"type"`
	UpdatedAt         int64         `json:"updated_at"`
}

type Requirements struct {
	Hardware        string `json:"hardware"`
	HardwareType    string `json:"hardware_type"`
	Memory          string `json:"memory"`
	PreferredCpList any    `json:"preferred_cp_list"`
	Region          string `json:"region"`
	Storage         string `json:"storage"`
	UpdateMaxLag    any    `json:"update_max_lag"`
	Vcpu            string `json:"vcpu"`
}

type ActiveOrder struct {
	Config Config `json:"config"`
}

type Config struct {
	Description  string  `json:"description"`
	Hardware     string  `json:"hardware"`
	HardwareID   int64   `json:"hardware_id"`
	HardwareType string  `json:"hardware_type"`
	Memory       int64   `json:"memory"`
	Name         string  `json:"name"`
	PricePerHour float64 `json:"price_per_hour"`
	Vcpu         int64   `json:"vcpu"`
}

type Space struct {
	ActiveOrder *ActiveOrder `json:"activeOrder"`
	Name        string       `json:"name"`
	UUID        string       `json:"uuid"`
}

type PageResult struct {
	List      any   `json:"list"`
	Page      int64 `json:"page"`
	Size      int64 `json:"size"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
}

type ConfigOrder struct {
	ConfigID        int64  `json:"config_id"`
	CreatedAt       int64  `json:"created_at"`
	Duration        int64  `json:"duration"`
	EndedAt         int    `json:"ended_at"`
	ErrorCode       int    `json:"error_code"`
	ID              int64  `json:"id"`
	OrderType       string `json:"order_type"`
	PreferredCpList any    `json:"preferred_cp_list"`
	RefundTxHash    string `json:"refund_tx_hash"`
	Region          string `json:"region"`
	SpaceID         string `json:"space_id"`
	StartIn         int64  `json:"start_in"`
	StartedAt       int64  `json:"started_at"`
	Status          string `json:"status"`
	TaskUUID        string `json:"task_uuid"`
	TxHash          string `json:"tx_hash"`
	ApproveHash     string `json:"approve_hash"`
	UpdatedAt       int64  `json:"updated_at"`
	UUID            string `json:"uuid"`
}

type TaskQueryReq struct {
	Wallet string `form:"wallet"`
	Page   uint   `form:"page"`
	Size   uint   `form:"size"`
}

/*
CreateTaskReq

PrivateKey:   The wallet's private key
WalletAddress:The wallet's address
InstanceType: The type(name) of the hardware. (Default = `C1ae.small`)
Region:       The region of the hardware. (Default: global)
Duration:     The duration of the service runtime in seconds. (Default = 3600)
JobSourceUri: Optional. The job source URI to be deployed. If this is provided, app_repo_image and repo_uri are ignored.
RepoUri:      Optional. The URI of the repo to be deployed. If job_source_uri and app_repo_image are not provided, this is required.
RepoBranch:   Optional. The branch of the repo to be deployed.
StartIn:      Optional. The starting time (expected time for the app to be deployed, not mandatory). (Default = 300)
PreferredCpList:  Optional. A list of preferred cp account address(es).
*/
type CreateTaskReq struct {
	PrivateKey      string        `json:"private_key,omitempty"`
	WalletAddress   string        `json:"wallet_address"`
	InstanceType    string        `json:"instance_type"`
	Region          string        `json:"region"`
	Duration        time.Duration `json:"duration"`
	JobSourceUri    string        `json:"job_source_uri"`
	RepoUri         string        `json:"repo_uri"`
	RepoBranch      string        `json:"repo_branch"`
	StartIn         int           `json:"start_in"`
	PreferredCpList []string      `json:"preferred_cp_list"`
}

type CreateTaskResp struct {
	Task         Task        `json:"task"`
	ConfigOrder  ConfigOrder `json:"config_order"`
	TxHash       string      `json:"tx_hash"`
	ApproveHash  string      `json:"approve_hash"`
	TaskUuid     string      `json:"task_uuid"`
	InstanceType string      `json:"instance_type"`
	Price        float64     `json:"price"`
}

func (task *CreateTaskReq) WithPrivateKey(privateKey string) {
	task.PrivateKey = privateKey
}

type RepoImageResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type JobSourceUriResult struct {
	JobSourceUri string `json:"job_source_uri"`
}

type ValidatePaymentResult struct {
	ConfigOrder
}

type PaymentResult struct {
	ConfigOrder
}

type RenewTaskResp struct {
	ConfigOrder ConfigOrder `json:"config_order"`
	Task        Task        `json:"task"`
}

type TerminateTaskResp struct {
	Retryable  bool   `json:"retryable"`
	TaskStatus string `json:"task_status"`
}

type ContractResult struct {
	ContractInfo ContractInfo `json:"contract_info"`
	Signature    string       `json:"signature"`
}

type ContractInfo struct {
	ContractDetail ContractDetail `json:"contract_detail"`
	Time           int            `json:"time"`
}

type ContractDetail struct {
	ClientContractAddress    string `json:"client_contract_address"`
	PaymentContractAddress   string `json:"payment_contract_address"`
	RpcUrl                   string `json:"rpc_url"`
	SwanTokenContractAddress string `json:"swan_token_contract_address"`
}
