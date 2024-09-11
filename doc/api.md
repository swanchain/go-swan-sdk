# APIs

- [APIs](#apis)
  - [NewClient](#newclient)
  - [Instances](#instances)
  - [Create task](#create-task)
  - [PayAndDeployTask](#payanddeploytask)
  - [EstimatePayment](#estimatepayment)
  - [RenewTask](#renewtask)
  - [RenewPayment](#renewpayment)
  - [TerminateTask](#terminatetask)
  - [GetRealUrl](#getrealurl)
  - [Tasks](#tasks)
  - [TaskInfo](#taskinfo)
- [Models](#models)
  - [InstanceBaseInfo](#instancebaseinfo)
  - [RegionDetail](#regiondetail)
  - [Task](#task)
  - [TaskDetail](#taskdetail)
  - [Requirements](#requirements)
  - [Space](#space)
  - [ActiveOrder](#activeorder)
  - [Config](#config)
  - [ConfigOrder](#configorder)
  - [TaskInfo](#taskinfo-1)
  - [ComputingProvider](#computingprovider)
  - [Job](#job)

## NewClient

Definition:
Creates a Swan Client instance

```go
func NewClient(apiKey, isTestnet) *APIClient
```
Inputs:

| Field Name | type   | description                                   |
| ---------- | ------ | --------------------------------------------- |
| apiKey     | string | Swan API key                                  |
| isTestnet  | bool   | If set to true use testnet, otherwise maninet |

Outputs:

```shell
*APIClient            # Created swan Client instance.
```

## Instances

`Instances` Fetch instance resources

```go
(c *APIClient) InstanceResources(available ...bool) ([]*InstanceResource, error)
```

Inputs:

| Field Name | type | description                                                       |
| ---------- | ---- | ----------------------------------------------------------------- |
| available  | bool | If set to true to get only available instances, otherwise get all |


Outputs:

| Field Name       | Type                                      | Description                                                                                                                                            |
| ---------------- | ----------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| InstanceBaseInfo | [InstanceBaseInfo](#InstanceBaseInfo)     | Contains the basic details of the instance.                                                                                                            |
| Region           | `[]string`                                | A slice of strings representing the regions where this instance is available.                                                                          |
| RegionDetails    | map[string]*[RegionDetail](#regiondetail) | A map where the key is the region name, and the value is a pointer to `RegionDetail` struct, containing detailed information for that specific region. |

## Create task
```go
func (c *APIClient) CreateTask(req *CreateTaskReq) (CreateTaskResp, error)
```
Inputs:

| Field Name      | Type            | Description                                                                                                                            |
| --------------- | --------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| PrivateKey      | `string`        | The private key of the user's wallet.                                                                                                  |
| WalletAddress   | `string`        | The user's wallet address.                                                                                                             |
| InstanceType    | `string`        | instance type of instance config. Defaults to 'C1ae.small' (Free tier). All Instance type list can be found [here](./instance.md)  |
| Region          | `string`        | The region where the task will be executed.                                                                                            |
| Duration        | `time.Duration` | The duration for which the task will run.                                                                                              |
| RepoUri         | `string`        | The The URI of the repo to be deployed. The repository must contain a `Dockerfile` or `deploy.yaml`. Please see [RepoUri](repo_uri.md) |
| RepoBranch      | `string`        | The branch of the repo to be deployed.                                                                                                 |
| StartIn         | `int`           | The delay (in seconds) before the task starts. (Default: 300)                                                                          |
| PreferredCpList | `[]string`      | A list of preferred cp account addresses.                                                                                              |

**Note:**
- `RepoUri` must contain a `Dockerfile` or `deploy.yaml` 
- `deploy.yaml` needs to follow the following [standards](https://docs.lagrangedao.org/spaces/intro/lagrange-definition-language-ldl) 


Outputs:

| Field Name   | Type                        | Description                                                          |
| ------------ | --------------------------- | -------------------------------------------------------------------- |
| Task         | [Task](#task)               | The `Task` struct containing details about the created task.         |
| ConfigOrder  | [ConfigOrder](#configorder) | The `ConfigOrder` struct containing the configuration order details. |
| TxHash       | `string`                    | Transaction hash for the task creation.                              |
| ApproveHash  | `string`                    | Transaction hash for the token approve.                              |
| TaskUuid     | `string`                    | Universally unique identifier for the task.                          |
| InstanceType | `string`                    | Type of instance created for the task.                               |
| Price        | `float64`                   | Price for the task.                                                  |



## PayAndDeployTask
```go
func (c *APIClient) PayAndDeployTask(taskUuid, privateKey string, duration time.Duration, instanceType string) (PaymentResult, error)
```
Inputs:

| Field Name   | Type            | Description                                                          |
| ------------ | --------------- | -------------------------------------------------------------------- |
| taskUuid     | `string`        | The universally unique identifier (UUID) of the task to be deployed. |
| privateKey   | `string`        | The private key used for payment authorization.                      |
| duration     | `time.Duration` | The duration for which the task will run.                            |
| instanceType | `string`        | The type of instance to be used for the task.                        |

Outputs:

| Field Name  | Type                                  | Description                               |
| ----------- | ------------------------------------- | ----------------------------------------- |
| ConfigOrder | [ConfigOrder](#configorder)(embedded) | Information about the configuration order |



## EstimatePayment
```go
func (c *APIClient) EstimatePayment(instanceType string, duration time.Duration) (float64, error)
```
Inputs:
| Field Name   | Type            | Description                                   |
| ------------ | --------------- | --------------------------------------------- |
| duration     | `time.Duration` | The duration for which the task will run.     |
| instanceType | `string`        | The type of instance to be used for the task. |

Outputs:
| Field Name | Type      | Description                 |
| ---------- | --------- | --------------------------- |
| -          | `float64` | The estimate payment price. |

## RenewTask
```go
(c *APIClient) RenewTask(taskUuid string, duration time.Duration, privateKey string, paidTxHash ...string) (*RenewTaskResp, error) 
```
In:
| Field Name | Type            | Description                                                          |
| ---------- | --------------- | -------------------------------------------------------------------- |
| taskUuid   | `string`        | The universally unique identifier (UUID) of the task to be deployed. |
| duration   | `time.Duration` | The duration for which the task will run.                            |
| privateKey | `string`        | The private key used for payment authorization.                      |
| paidTxHash | `string`        | The paid tx_hash.                                                    |

Outputs:

| Field Name  | Type                                  | Description                                                                |
| ----------- | ------------------------------------- | -------------------------------------------------------------------------- |
| ConfigOrder | [ConfigOrder](#configorder)(embedded) | The `ConfigOrder` struct containing details about the configuration order. |
| Task        | [Task](#task)                         | The `Task` struct containing details about the task.                       |


## RenewPayment
```go
func (c *APIClient) RenewPayment(taskUuid string, duration time.Duration, privateKey string) (string, error)
```
Inputs:
| Field Name | Type            | Description                                                          |
| ---------- | --------------- | -------------------------------------------------------------------- |
| taskUuid   | `string`        | The universally unique identifier (UUID) of the task to be deployed. |
| duration   | `time.Duration` | The duration for which the task will run.                            |
| privateKey | `string`        | The private key used for payment authorization.                      |

Outputs:
| Field Name | Type     | Description       |
| ---------- | -------- | ----------------- |
| -          | `string` | The paid tx_hash. |

## TerminateTask
```go
func (c *APIClient) TerminateTask(taskUuid string) (TerminateTaskResp, error)
```
Inputs:
| Field Name | Type     | Description                                                          |
| ---------- | -------- | -------------------------------------------------------------------- |
| taskUuid   | `string` | The universally unique identifier (UUID) of the task to be deployed. |

Outputs:
| Field Name | Type     | Description                                                   |
| ---------- | -------- | ------------------------------------------------------------- |
| Retryable  | `bool`   | Indicates whether the task termination is retryable.          |
| TaskStatus | `string` | The current status of the task after the termination attempt. |


## GetRealUrl
```go
func (c *APIClient) GetRealUrl(taskUuid string) ([]string, error) 
```
Inputs:
| Field Name | Type     | Description                                          |
| ---------- | -------- | ---------------------------------------------------- |
| taskUuid   | `string` | The universally unique identifier (UUID) of the task |

Outputs:
| Field Name | Type       | Description                  |
| ---------- | ---------- | ---------------------------- |
| -          | `[]string` | The application access urls. |

## Tasks
```go
func (c *APIClient) Tasks(req *TaskQueryReq) (total int64, list []*TaskInfo, err error)
```
Input:
| Field Name | Type     | Description                                 |
| ---------- | -------- | ------------------------------------------- |
| Wallet     | `string` | The wallet address used for querying tasks. |
| Page       | `uint`   | The page number for pagination.             |
| Size       | `uint`   | The number of tasks per page.               |

Outputs:
| Field Name | Type                       | Description                       |
| ---------- | -------------------------- | --------------------------------- |
| total      | int64                      | The total of data.                |
| list       | []*[TaskInfo](#taskinfo-1) | The containing task list details. |




## TaskInfo
```go
func (c *APIClient) TaskInfo(taskUUID string) (*TaskInfo, error) 
```
Input:
| Field Name | Type     | Description                                          |
| ---------- | -------- | ---------------------------------------------------- |
| taskUUID   | `string` | The universally unique identifier (UUID) of the task |

Outputs:
| Field Name | Type                     | Description                  |
| ---------- | ------------------------ | ---------------------------- |
| TaskInfo   | *[TaskInfo](#taskinfo-1) | The containing task details. |



# Models

## InstanceBaseInfo
The `InstanceBaseInfo` struct holds the fundamental details about the instance.

| Field Name  | Type     | Description                           |
| ----------- | -------- | ------------------------------------- |
| Description | `string` | A description of the instance.        |
| ID          | `int64`  | A unique identifier for the instance. |
| Name        | `string` | The name of the instance.             |
| Price       | `string` | The price of the instance.            |
| Status      | `string` | The status of the instance.           |
| Type        | `string` | The type of the instance.             |

## RegionDetail
The `RegionDetail` struct, containing detailed information for that specific region.

| Field Name        | Type         | Description                                                                |
| ----------------- | ------------ | -------------------------------------------------------------------------- |
| AvailableResource | `int64`      | The amount of available resources for the instance in this region.         |
| DirectAccessCp    | `[][]string` | A 2D slice of strings representing the direct access CPs for the instance. |
| NoneCollateral    | `int64`      | The amount of resources that have no collateral in this region.            |
| Whitelist         | `int64`      | The number of whitelisted items related to the instance in this region.    |

## Task
The `Task` struct contains information about a task, including its lifecycle and associated details.

| Field Name    | Type                       | Description                                                              |
| ------------- | -------------------------- | ------------------------------------------------------------------------ |
| Comments      | `string`                   | Comments or notes about the task.                                        |
| CreatedAt     | `int64`                    | Timestamp when the task was created.                                     |
| EndAt         | `int64`                    | Timestamp when the task ended.                                           |
| ID            | `int64`                    | Unique identifier for the task.                                          |
| LeadingJobID  | `string`                   | Identifier for the leading job associated with this task.                |
| Name          | `string`                   | Name of the task.                                                        |
| RefundAmount  | `string`                   | Amount to be refunded for the task.                                      |
| RefundWallet  | `string`                   | Wallet address where the refund should be sent.                          |
| Source        | `string`                   | Source of the task.                                                      |
| StartAt       | `int64`                    | Timestamp when the task started.                                         |
| StartIn       | `int64`                    | Delay (in seconds) before the task starts.                               |
| Status        | `string`                   | Current status of the task.                                              |
| TaskDetail    | *[TaskDetail](#taskdetail) | Pointer to the `TaskDetail` struct containing detailed task information. |
| TaskDetailCid | `string`                   | CID (Content Identifier) for the task details.                           |
| TxHash        | `any`                      | Transaction hash associated with the task.                               |
| Type          | `string`                   | Type of the task.                                                        |
| UpdatedAt     | `int64`                    | Timestamp when the task was last updated.                                |
| UserID        | `int64`                    | Unique identifier for the user associated with the task.                 |
| UUID          | `string`                   | Universally unique identifier for the task.                              |

## TaskDetail
The `TaskDetail` struct contains detailed information about the task, including resource requirements and status.

| Field Name        | Type                           | Description                                                            |
| ----------------- | ------------------------------ | ---------------------------------------------------------------------- |
| Amount            | `float64`                      | Amount associated with the task.                                       |
| BidderLimit       | `int64`                        | Limit on the number of bidders.                                        |
| CreatedAt         | `int64`                        | Timestamp when the task detail was created.                            |
| DCCSelectedCpList | `any`                          | List of selected CPs.                                                  |
| Duration          | `int64`                        | Duration (seconds) for which the task will run.                        |
| EndAt             | `int64`                        | Timestamp when the task ends.                                          |
| Hardware          | `string`                       | Instance configuration for the task.                                   |
| JobResultURI      | `string`                       | URI where the job results can be accessed.                             |
| JobSourceURI      | `string`                       | URI where the job source is located.                                   |
| PricePerHour      | `string`                       | Price per hour for the task.                                           |
| Requirements      | *[Requirements](#requirements) | Pointer to the `Requirements` struct containing resource requirements. |
| Space             | *[Space](#space)               | Pointer to the `Space` struct for storage space details.               |
| StartAt           | `int64`                        | Timestamp when the task starts.                                        |
| Status            | `string`                       | Current status of the task.                                            |
| StorageSource     | `string`                       | Source of the storage used by the task.                                |
| Type              | `string`                       | Type of the task.                                                      |
| UpdatedAt         | `int64`                        | Timestamp when the task detail was last updated.                       |

## Requirements

The `Requirements` struct defines the instance and resource requirements for the task.

| Field Name      | Type     | Description                                  |
| --------------- | -------- | -------------------------------------------- |
| Hardware        | `string` | Instance required for the task.              |
| HardwareType    | `string` | Type of instance required.                   |
| Memory          | `string` | Memory required for the task.                |
| PreferredCpList | `any`    | List of preferred CPs.                       |
| Region          | `string` | Region where the instance should be located. |
| Storage         | `string` | Storage requirements for the task.           |
| UpdateMaxLag    | `any`    | Maximum allowable lag for updates.           |
| Vcpu            | `string` | Number of virtual CPUs required.             |

## Space
The `Space` struct represents a space entity that includes information about its active order, name, and universally unique identifier (UUID).

| Field Name  | Type                         | Description                                                                |
| ----------- | ---------------------------- | -------------------------------------------------------------------------- |
| ActiveOrder | *[ActiveOrder](#activeorder) | Pointer to the `ActiveOrder` struct representing the current active order. |
| Name        | `string`                     | The name of the space.                                                     |
| UUID        | `string`                     | The universally unique identifier (UUID) of the space.                     |


## ActiveOrder

The `ActiveOrder` struct contains configuration details for an active order.

| Field Name | Type              | Description                                               |
| ---------- | ----------------- | --------------------------------------------------------- |
| Config     | [Config](#config) | The `Config` struct containing the configuration details. |

## Config

The `Config` struct defines the configuration of a instance order, including pricing and resource allocation.

| Field Name   | Type      | Description                                 |
| ------------ | --------- | ------------------------------------------- |
| Description  | `string`  | Description of the instance configuration.  |
| Hardware     | `string`  | Instance associated with the configuration. |
| HardwareID   | `int64`   | Unique identifier for the instance.         |
| HardwareType | `string`  | Type of instance used.                      |
| Memory       | `int64`   | Amount of memory allocated.                 |
| Name         | `string`  | Name of the configuration.                  |
| PricePerHour | `float64` | Price per hour for the instance.            |
| Vcpu         | `int64`   | Number of virtual CPUs allocated.           |

## ConfigOrder

The `ConfigOrder` struct provides detailed information about the configuration order associated with the task.

| Field Name      | Type     | Description                                                           |
| --------------- | -------- | --------------------------------------------------------------------- |
| ConfigID        | `int64`  | Unique identifier for the configuration.                              |
| CreatedAt       | `int64`  | Timestamp when the configuration order was created.                   |
| Duration        | `int64`  | Duration(seconds) of the configuration order.                         |
| EndedAt         | `int`    | Timestamp when the configuration order ended.                         |
| ErrorCode       | `int`    | Error code associated with the configuration order, if any.           |
| ID              | `int64`  | Unique identifier for the order.                                      |
| OrderType       | `string` | Type of the order.                                                    |
| PreferredCpList | `any`    | List of preferred CPs for the configuration order.                    |
| RefundTxHash    | `string` | Transaction hash for the refund associated with the order.            |
| Region          | `string` | Region where the configuration is applied.                            |
| SpaceID         | `string` | Identifier for the space associated with the configuration order.     |
| StartIn         | `int64`  | Delay (in seconds) before the configuration order starts.             |
| StartedAt       | `int64`  | Timestamp when the configuration order started.                       |
| Status          | `string` | Current status of the configuration order.                            |
| TaskUUID        | `string` | Universally unique identifier for the task associated with the order. |
| TxHash          | `string` | Transaction hash associated with the order.                           |
| UpdatedAt       | `int64`  | Timestamp when the configuration order was last updated.              |
| UUID            | `string` | Universally unique identifier for the configuration order.            |

## TaskInfo

| Field Name | Type                                       | Description                                           |
| ---------- | ------------------------------------------ | ----------------------------------------------------- |
| Providers  | []*[ComputingProvider](#computingprovider) | List of computing providers associated with the task. |
| Orders     | []*[ConfigOrder](#configorder)             | List of configuration orders related to the task.     |
| Jobs       | []*[Job](#job)                             | List of jobs related to the task.                     |
| Task       | [Task](#task)                              | The task itself.                                      |

## ComputingProvider

The `ComputingProvider` struct contains information about computing providers, including their location, status, and other relevant details.

| Field Name       | Type       | Description                                                                      |
| ---------------- | ---------- | -------------------------------------------------------------------------------- |
| Beneficiary      | `string`   | The beneficiary of the computing provider.                                       |
| CpAccountAddress | `string`   | The account address of the computing provider.                                   |
| CreatedAt        | `int64`    | Timestamp when the computing provider was created.                               |
| FreezeOnline     | `any`      | Information about whether the provider's service is frozen online.               |
| ID               | `int64`    | Unique identifier for the computing provider.                                    |
| Lat              | `float64`  | Latitude of the computing provider's location.                                   |
| Lon              | `float64`  | Longitude of the computing provider's location.                                  |
| MultiAddress     | `[]string` | List of multiple addresses associated with the provider.                         |
| Name             | `string`   | Name of the computing provider.                                                  |
| NodeID           | `string`   | Identifier for the node.                                                         |
| Online           | `int`      | Status indicating if the provider is online (e.g., 1 for online, 0 for offline). |
| OwnerAddress     | `string`   | Address of the owner of the computing provider.                                  |
| Region           | `string`   | Region where the computing provider is located.                                  |
| TaskTypes        | `string`   | Types of tasks the provider can handle.                                          |
| UpdatedAt        | `int64`    | Timestamp when the computing provider was last updated.                          |
| Version          | `string`   | Version of the computing provider's software.                                    |
| WorkerAddress    | `string`   | Address of the worker associated with the provider.                              |

## Job

The `Job` struct contains details about individual jobs, including logs, status, and associated metadata.

| Field Name       | Type     | Description                                                 |
| ---------------- | -------- | ----------------------------------------------------------- |
| BuildLog         | `string` | Log of the build process.                                   |
| Comments         | `string` | Comments related to the job.                                |
| ContainerLog     | `string` | Log from the container execution.                           |
| CpAccountAddress | `string` | Account address of the computing provider handling the job. |
| CreatedAt        | `int64`  | Timestamp when the job was created.                         |
| Duration         | `int64`  | Duration of the job (in seconds).                           |
| EndedAt          | `int64`  | Timestamp when the job ended.                               |
| Hardware         | `string` | Description of the instance used for the job.               |
| ID               | `int64`  | Unique identifier for the job.                              |
| JobRealURI       | `string` | URI for the actual job resource.                            |
| JobResultURI     | `string` | URI for the job result.                                     |
| JobSourceURI     | `string` | URI for the source of the job.                              |
| Name             | `string` | Name of the job.                                            |
| NodeID           | `string` | Identifier for the node where the job was executed.         |
| StartAt          | `int64`  | Timestamp when the job started.                             |
| Status           | `string` | Current status of the job.                                  |
| StorageSource    | `string` | Source of storage used for the job.                         |
| TaskUUID         | `string` | UUID of the task associated with the job.                   |
| Type             | `any`    | Type of the job.                                            |
| UpdatedAt        | `int64`  | Timestamp when the job was last updated.                    |
| UUID             | `string` | Universally unique identifier for the job.                  |
