# GO SWAN SDK <!-- omit in toc -->

[![Made by FilSwan](https://img.shields.io/badge/made%20by-FilSwan-green.svg)](https://www.filswan.com/)
[![Chat on discord](https://img.shields.io/badge/join%20-discord-brightgreen.svg)](https://discord.com/invite/swanchain)

## Table Of Contents<!-- omit in toc -->

- [Quickstart](#quickstart)
  - [Installation](#installation)
  - [Get Orchestrator API Key](#get-orchestrator-api-key)
  - [Usage](#usage)
    - [New client](#new-client)
    - [Create task](#create-task)
      - [1. Automatic payment and deployment](#1-automatic-payment-and-deployment)
      - [2. Manual payment and deployment](#2-manual-payment-and-deployment)
    - [Get the access url of the application](#get-the-access-url-of-the-application)
    - [Renewal task](#renewal-task)
      - [1. Automatic payment](#1-automatic-payment)
      - [2. Manual payment](#2-manual-payment)
    - [Terminate Task](#terminate-task)
    - [Get Task Detail](#get-task-detail)
    - [Get Task List](#get-task-list)
- [A Sample Tutorial](#a-sample-tutorial)
- [License](#license)


## Quickstart

This guide details the steps needed to install or update the SWAN SDK for Go. The SDK is a comprehensive toolkit designed to facilitate seamless interactions with the SwanChain API.

### Installation

To use Swan SDK, you first need to install it and its dependencies. Before installing Swan SDK, install Go 1.22.3 or later.


Install the latest Swan SDK release via **go**:

```bash
go get github.com/swanchain/go-swan-sdk
```

### Get Orchestrator API Key

To use `swan-sdk`, an Orchestrator API key is required.

Steps to get an API Key:

- Go to [Orchestrator Dashboard](https://orchestrator.swanchain.io/provider-status), switch network to Mainnet.
- Login through MetaMask.
- Click the user icon on the top right.
- Click 'Show API-Key' -> 'New API Key'
- Store your API Key safely, do not share with others.


### Usage

#### [New client]()

```go
import "github.com/swanchain/go-swan-sdk"

client := swan.NewAPIClient("<SWAN_API_KEY>")
```

#### [Create task]()

##### 1. Automatic payment and deployment
```go
var req = CreateTaskReq{
    PrivateKey: "<YOUR_WALLET_ADDRESS_PRIVATE_KEY>",
    AutoPay:    true, 
    RepoUri:    "<Your_RESOURCE_URL>",    
}
createTaskResp, err := client.CreateTask(&createReq)
if err != nil {
    log.Fatalln(err)
}
log.Printf("task result: %v", createTaskResp)
```

##### 2. Manual payment and deployment
```go
var req = CreateTaskReq{
    PrivateKey: "<YOUR_WALLET_ADDRESS_PRIVATE_KEY>",
    AutoPay:      false,
    RepoUri:      "<Your_RESOURCE_URL>",
    Duration:     3600, 
    InstanceType: "C1ae.small", 
}
createTaskResp, err := client.CreateTask(&createReq)
if err != nil {
    log.Fatalln(err)
}
log.Printf("task result: %v", createTaskResp)

taskUuid := "<TASK_UUID>" // taskUuid: returned by create task
payAndDeployTaskResp, err := apiClient.PayAndDeployTask(taskUuid, <YOUR_WALLET_ADDRESS_PRIVATE_KEY>, <Duration>, <InstanceType>)
if err != nil {
    log.Fatalln(err)
}
log.Printf("pay and deploy task response: %v", payAndDeployTaskResp)
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


#### [Renewal task]()

##### 1. Automatic payment
```go
resp, err := apiClient.ReNewTask("<TASK_UUID>", <Duration>, true,"<YOUR_WALLET_ADDRESS_PRIVATE_KEY>", "")
if err != nil {
	log.Fatalln(err)
}
log.Printf("renew task with auto-pay response: %v", resp)
```

##### 2. Manual payment
```go
txHash, err := apiClient.RenewPayment("<TASK_UUID>", <Duration>, "<YOUR_WALLET_ADDRESS_PRIVATE_KEY>")
if err != nil {
	log.Fatalln(err)
}

resp, err := apiClient.ReNewTask("<TASK_UUID>", <Duration>, false, "",txHash)
if err != nil {
	log.Fatalln(err)
}
log.Printf("renew task with manual-pay response: %v", resp)
```

#### [Terminate Task]()

```go
resp, err := apiClient.TerminateTask("<TASK_UUID>")
if err != nil {
    log.Fatalln(err)
}
log.Printf("terminate task response: %v", resp)
```

#### [Get Task Detail]()
```go
resp, err := apiClient.TaskInfo("<TASK_UUID>")
if err != nil {
    log.Fatalln(err)
}
log.Printf("get task info response: %v", resp)
```

#### [Get Task List]()
```go
var req = &TaskQueryReq{
    Wallet: "<YOUR_WALLET_ADDRESS>",
    Page:   0,
    Size:   10,
}

total, resp, err := apiClient.Tasks(req)
if err != nil {
    log.Fatalln(err)
}
log.Printf("get task list response, total: %d, data: %v", total, resp)
```

## A Sample Tutorial

For more detailed samples, consult [SDK Samples](https://github.com/swanchain/go-swan-sdk/blob/main/client_test.go).


## License

The GOLANG SWAN SDK is released under the **MIT** license, details of which can be found in the LICENSE file.
