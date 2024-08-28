package swan

type HardwareResult struct {
	Hardware []*Hardware `json:"hardware"`
}

type Hardware struct {
	Description string   `json:"hardware_description"`
	ID          int64    `json:"hardware_id"`
	Name        string   `json:"hardware_name"`
	Price       string   `json:"hardware_price"`
	Status      string   `json:"hardware_status"`
	Type        string   `json:"hardware_type"`
	Region      []string `json:"region"`
}

type TaskDetails struct {
	Providers []*ComputingProvider `json:"computing_providers"`
	Jobs      []*Job               `json:"jobs"`
	Task      *Task                `json:"task"`
}

type ComputingProvider struct {
	Beneficiary      string   `json:"beneficiary"`
	CpAccountAddress string   `json:"cp_account_address"`
	CreatedAt        int64    `json:"created_at"`
	ID               int64    `json:"id"`
	Lat              float64  `json:"lat"`
	Lon              float64  `json:"lon"`
	MultiAddress     []string `json:"multi_address"`
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
	BuildLog      string      `json:"build_log"`
	Comments      string      `json:"comments"`
	ContainerLog  string      `json:"container_log"`
	CreatedAt     int64       `json:"created_at"`
	Duration      int64       `json:"duration"`
	EndedAt       int64       `json:"ended_at"`
	Hardware      string      `json:"hardware"`
	ID            int64       `json:"id"`
	JobRealURI    string      `json:"job_real_uri"`
	JobResultURI  string      `json:"job_result_uri"`
	JobSourceURI  string      `json:"job_source_uri"`
	Name          string      `json:"name"`
	NodeID        string      `json:"node_id"`
	StartAt       int64       `json:"start_at"`
	Status        string      `json:"status"`
	StorageSource string      `json:"storage_source"`
	TaskUUID      string      `json:"task_uuid"`
	Type          interface{} `json:"type"`
	UpdatedAt     int64       `json:"updated_at"`
	UUID          string      `json:"uuid"`
}

type Task struct {
	Comments      interface{} `json:"comments"`
	CreatedAt     int64       `json:"created_at"`
	EndAt         int64       `json:"end_at"`
	ID            int64       `json:"id"`
	LeadingJobID  string      `json:"leading_job_id"`
	RefundAmount  float64     `json:"refund_amount"`
	RefundWallet  string      `json:"refund_wallet"`
	Source        string      `json:"source"`
	StartAt       int64       `json:"start_at"`
	StartIn       int64       `json:"start_in"`
	Status        string      `json:"status"`
	TaskDetailCid string      `json:"task_detail_cid"`
	TxHash        string      `json:"tx_hash"`
	Type          interface{} `json:"type"`
	UpdatedAt     int64       `json:"updated_at"`
	UUID          string      `json:"uuid"`
}
