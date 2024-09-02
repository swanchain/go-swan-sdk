package swan

import (
	"reflect"
	"testing"
)

func TestAPIClient_CreateTask(t *testing.T) {
	type args struct {
		req *CreateTaskReq
	}
	test := struct {
		name    string
		apiKey  string
		args    args
		want    CreateTaskResp
		wantErr bool
	}{
		name: "createTask",
		args: args{
			req: &CreateTaskReq{
				WalletAddress: "",
				PrivateKey:    "",
				HardwareId:    0,
				Region:        "global",
				Duration:      3600,
				AutoPay:       true,
				JobSourceUri:  "",
				StartIn:       300,
			},
		},
	}
	t.Run(test.name, func(t *testing.T) {
		apiClient := NewAPIClient("")
		got, err := apiClient.CreateTask(test.args.req)
		if (err != nil) != test.wantErr {
			t.Errorf("CreateTask() error = %v, wantErr %v", err, test.wantErr)
			return
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("CreateTask() got = %v, want %v", got, test.want)
		}
	})

}

func TestAPIClient_TaskInfo(t *testing.T) {
	type fields struct {
		apiKey       string
		httpClient   *HttpClient
		contractInfo ContractDetail
	}
	type args struct {
		taskUUID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TaskInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &APIClient{
				apiKey:     tt.fields.apiKey,
				httpClient: tt.fields.httpClient,
			}
			got, err := c.TaskInfo(tt.args.taskUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAPIClient(t *testing.T) {
	type args struct {
		apiKey  string
		testnet []bool
	}
	tests := []struct {
		name string
		args args
		want *APIClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIClient(tt.args.apiKey, tt.args.testnet...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPIClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResult(t *testing.T) {
	type args struct {
		dest any
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResult(tt.args.dest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_Check(t *testing.T) {
	type fields struct {
		Data    any
		Message string
		Status  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{
				Data:    tt.fields.Data,
				Message: tt.fields.Message,
				Status:  tt.fields.Status,
			}
			if err := r.Check(); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getContractInfo(t *testing.T) {
	type args struct {
		validate bool
		apiKey   string
	}
	test := struct {
		name    string
		args    args
		want    ContractDetail
		wantErr bool
	}{
		name: "getContractInfo",
		args: args{
			validate: false,
			apiKey:   "53Qkrwdeyv",
		},
		wantErr: false,
	}
	t.Run(test.name, func(t *testing.T) {
		apiClient := NewAPIClient(test.args.apiKey, true)
		got, err := apiClient.getContractInfo(test.args.validate)
		if (err != nil) != test.wantErr {
			t.Errorf("getcontractInfo() error = %v, wantErr %v", err, test.wantErr)
			return
		}
		t.Logf("getcontractInfo() got = %v", got)
	})
}
