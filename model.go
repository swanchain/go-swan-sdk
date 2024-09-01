package swan

type HardwareResult struct {
	Hardware []*Hardware `json:"hardware"`
}

type Hardware struct {
	Description   string                   `json:"hardware_description"`
	ID            int64                    `json:"hardware_id"`
	Name          string                   `json:"hardware_name"`
	Price         string                   `json:"hardware_price"`
	Status        string                   `json:"hardware_status"`
	Type          string                   `json:"hardware_type"`
	Region        []string                 `json:"region"`
	RegionDetails map[string]*RegionDetail `json:"region_detail"`
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

type Space struct {
	ActiveOrder *ActiveOrder `json:"activeOrder"`
	Name        string       `json:"name"`
	UUID        string       `json:"uuid"`
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
	UpdatedAt       int64  `json:"updated_at"`
	UUID            string `json:"uuid"`
}

type TaskQueryReq struct {
	Wallet string `form:"wallet"`
	Page   uint   `form:"page"`
	Size   uint   `form:"size"`
}

type CreateTaskReq struct {
	WalletAddress   string   `json:"wallet_address"`
	PrivateKey      string   `json:"private_key,omitempty"`
	HardwareId      int64    `json:"hardware_id"`
	Region          string   `json:"region"`
	Duration        int      `json:"duration"`
	AppRepoImage    string   `json:"app_repo_image"`
	AutoPay         bool     `json:"auto_pay"`
	JobSourceUri    string   `json:"job_source_uri"`
	RepoUri         string   `json:"repo_uri"`
	RepoBranch      string   `json:"repo_branch"`
	RepoOwner       string   `json:"repo_owner"`
	RepoName        string   `json:"repo_name"`
	StartIn         int      `json:"start_in"`
	PreferredCpList []string `json:"preferred_cp_list"`
}

type CreateTaskResp struct {
	TaskUuid    string      `json:"task_uuid"`
	TxHash      string      `json:"tx_hash"`
	Id          string      `json:"id"`
	HardwareId  int64       `json:"hardware_id"`
	ConfigOrder ConfigOrder `json:"config_order"`
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
}

type PaymentResult struct {
	ConfigOrder ConfigOrder `json:"config_order"`
	TxHash      string      `json:"tx_hash"`
}
