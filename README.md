# GO SWAN SDK <!-- omit in toc -->

[![Made by FilSwan](https://img.shields.io/badge/made%20by-FilSwan-green.svg)](https://www.filswan.com/)
[![Chat on discord](https://img.shields.io/badge/join%20-discord-brightgreen.svg)](https://discord.com/invite/swanchain)

`go-swan-sdk` is a comprehensive toolkit designed to facilitate seamless interactions with the SwanChain API.

## Table Of Contents<!-- omit in toc -->

- [Getting started](#getting-started)
  - [Prerequisites](#prerequisites)
    - [Go version](#go-version)
    - [Swan API Key](#swan-api-key)
  - [Getting go-swan-sdk](#getting-go-swan-sdk)
  - [Usage](#usage)
    - [New client](#new-client)
    - [Select a hardware form hardwares](#select-a-hardware-form-hardwares)
    - [Create a task to deploy an application](#create-a-task-to-deploy-an-application)
    - [Get the access url of the application](#get-the-access-url-of-the-application)
    - [Renew task duration](#renew-task-duration)
    - [Terminate Task](#terminate-task)
    - [Get Task Detail](#get-task-detail)
    - [Get Task List](#get-task-list)
- [A Sample Tutorial](#a-sample-tutorial)
- [License](#license)


## Getting started

### Prerequisites

#### Go version

`go-swan-sdk` requires [Go](https://go.dev/) version [1.21](https://go.dev/doc/devel/release#go1.21.0) or above.


#### Swan API Key

To use `swan-sdk`, an Orchestrator API key is required.

Steps to get an API Key:

- Go to [Orchestrator Dashboard](https://orchestratorswanchain.io/provider-status), switch network toMainnet.
- Login through MetaMask.
- Click the user icon on the top right.
- Click 'Show API-Key' -> 'New API Key'
- Store your API Key safely, do not share with others.

### Getting go-swan-sdk

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```sh
import "github.com/swanchain/go-swan-sdk"
```

Alternatively, use `go get`:

```sh
go get -u github.com/swanchain/go-swan-sdk
```

### Usage

#### [New client]()

```go
client, err := swan.NewAPIClient("<SWAN_API_KEY>")
```

#### [Select a hardware form hardwares]()

`Hardwares` lists all hardwares, you can select a  hardware you want.

```go
hardwares, err := swan.Hardwares()
```

#### [Create a task to deploy an application]()

create, pay and deploy a task

```go
createTaskResp, err := client.CreateTask(&CreateTaskReq{
    PrivateKey:   "<YOUR_WALLET_ADDRESS_PRIVATE_KEY>",
    RepoUri:      "<Your_RESOURCE_URL>",
    Duration:     time.Duration(3600),
    InstanceType: "C1ae.small", // hardware_type
})

taskUUID := createTaskResp.Task.UUID
```

#### [Get the access url of the application]()
```go
// Get application instances URL
appUrls, err := client.GetRealUrl("<TASK_UUID>")
if err != nil {
	log.Fatalln(err)
}
log.Printf("app urls: %v", appUrls)
```
A sample output:

```
['https://krfswstf2g.anlu.loveismoney.fun', 'https://l2s5o476wf.cp162.bmysec.xyz', 'https://e2uw19k9uq.cp5.node.study']
```

It shows that this task has three applications. Open the URL in the web browser you will view the application's information if it is running correctly.


#### [Renew task duration]()

`RenewTask` extends the duration of the task before completed

```go
resp, err := apiClient.RenewTask("<TASK_UUID>", <Duration>,"<YOUR_WALLET_ADDRESS_PRIVATE_KEY>", "")
```

#### [Terminate Task]()

`TerminateTask`  terminates the task

```go
resp, err := apiClient.TerminateTask("<TASK_UUID>")
```

#### [Get Task Detail]()
```go
resp, err := apiClient.TaskInfo("<TASK_UUID>")
```

#### [Get Task List]()
```go
total, resp, err := apiClient.Tasks(&TaskQueryReq{
    Wallet: "<PAY_WALLET_ADDRESS>",
    Page:   0,
    Size:   10,
})
```

## A Sample Tutorial

For more detailed samples, consult [SDK Samples](https://github.com/swanchain/go-swan-sdk/blob/main/client_test.go).


## License

The `go-swan-sdk` is released under the **MIT** license, details of which can be found in the LICENSE file.
