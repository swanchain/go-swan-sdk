# PYTHON SWAN SDK <!-- omit in toc -->

[![Made by FilSwan](https://img.shields.io/badge/made%20by-FilSwan-green.svg)](https://www.filswan.com/)
[![Chat on discord](https://img.shields.io/badge/join%20-discord-brightgreen.svg)](https://discord.com/invite/swanchain)

## Table Of Contents<!-- omit in toc -->

- [Quickstart](#quickstart)
    - [Installation](#installation)
    - [Get Orchestrator API Key](#get-orchestrator-api-key)
    - [Using Swan](#using-swan)
- [A Sample Tutorial](#a-sample-tutorial)
    - [Orchestrator](#orchestrator)
        - [Create and deploy a task](#create-and-deploy-a-task)
        - [Check information of an existing task](#check-information-of-an-existing-task)
        - [Access application instances of an existing task](#access-application-instances-of-an-existing-task)
        - [Renew an existing task](#renew-an-existing-task)
        - [Terminate an existing task](#terminate-an-existing-task)
- [License](#license)


## Quickstart

This guide details the steps needed to install or update the SWAN SDK for Golang. The SDK is a comprehensive toolkit designed to facilitate seamless interactions with the SwanChain API.

### Installation

To use Swan SDK, you first need to install it and its dependencies. Before installing Swan SDK, install Golang 1.22.3 or later.


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

To use Swan SDK, you must first import it and indicate which service you're going to use:

```go
import "github.com/swanchain/go-swan-sdk"

client := swan.NewAPIClient("<SWAN_API_KEY>")
```

Now that you have an `Orchestrator` service, you can create and deploy instance applications as an Orchestrator task with the service.

```go
var createReq = swan.CreateTaskReq{
    WalletAddress: "",
    PrivateKey:    "",
    RepoUri:       "",
    AutoPay:       true,
}
createTaskResp, err := client.CreateTask(&createReq)
if err != nil {
log.Fatalln(err)
}
log.Printf("task result: %v", createTaskResp)
```

Then you can follow up task deployment information and the URL for running applications.

```go
// Get task deployment info
taskInfo, err := client.TaskInfo(createTaskResp.TaskUuid)
if err != nil {
log.Fatalln(err)
}
log.Printf("task info: %v", taskInfo)

// Get application instances URL
appUrls, err := client.GetRealUrl(createTaskResp.TaskUuid)
if err != nil {
log.Fatalln(err)
}
log.Printf("app urls: %v", appUrls)
```

## A Sample Tutorial

For more detailed samples, consult [SDK Samples](https://github.com/swanchain/github.com/swanchain/go-swan-sdk/sample).

### Orchestrator

Orchestrator allows you to create task to run application instances to the powerful distributed computing providers network.

#### Create and deploy a task

Deploy a simple application with Swan SDK:

```go
import (
    "github.com/swanchain/go-swan-sdk"
    "log"
)

client := swan.NewAPIClient("<SWAN_API_KEY>")

var createReq = swan.CreateTaskReq{
    WalletAddress: "",
    PrivateKey:    "",
    RepoUri:       "",
    AutoPay:       true,
}

// create task
createTaskResp, err := client.CreateTask(&createReq)
if err != nil {
    log.Fatalln(err)
}
log.Printf("task result: %v", createTaskResp)

// get task info
taskInfo, err := client.TaskInfo(createTaskResp.TaskUuid)
if err != nil {
   log.Fatalln(err)
}
log.Printf("task info: %v", taskInfo)

```

It may take up to 5 minutes to get the deployment result:

```go
// Get application instances URL
appUrls, err := client.GetRealUrl(createTaskResp.TaskUuid)
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

#### Check information of an existing task

With Orchestrator, you can check information for an existing task to follow up or view task deployment.

```go
import (
    "github.com/swanchain/go-swan-sdk"
    "log"
)

client := swan.NewAPIClient("<SWAN_API_KEY>")

// Get an existing task deployment info
taskInfo, err := client.TaskInfo(createTaskResp.TaskUuid)
if err != nil {
    log.Fatalln(err)
}
log.Printf("task info: %v", taskInfo)
```

#### Access application instances of an existing task

With Orchestrator, you can easily get the deployed application instances for an existing task.

```go
import (
"github.com/swanchain/go-swan-sdk"
"log"
)

client := swan.NewAPIClient("<SWAN_API_KEY>")

// Get application instances URL
appUrls, err := client.GetRealUrl(createTaskResp.TaskUuid)
if err != nil {
  log.Fatalln(err)
}
log.Printf("app urls: %v", appUrls)
```

#### Renew an existing task

If you have already submitted payment for the renewal of a task, you can use the `tx_hash` with `renew_task` to extend the task.

```go
import (
    "github.com/swanchain/go-swan-sdk"
    "log"
)

client := swan.NewAPIClient("<SWAN_API_KEY>")

taskUuid := "taskUuid"
txHash := "txHash"
privateKey := "privateKey"
instanceType := "instanceType"
duration := 3600
reNewTask, err := client.ReNewTask(taskUuid, txHash, privateKey, instanceType, duration, true)
if err != nil {
  log.Fatalln(err)
}
log.Printf("renew task result: %v", reNewTask)

```

#### Terminate an existing task

You can also early terminate an existing task and its application instances. By terminating task, you will stop all the related running application instances and thus you will get refund of the remaining task duration.

```go
import (
    "github.com/swanchain/go-swan-sdk"
    "log"
)

client := swan.NewAPIClient("<SWAN_API_KEY>")

// Terminate an existing task (and its application instances)
terminateTask, err := client.TerminateTask(taskUuid)
if err != nil {
    log.Fatalln(err)
}
log.Printf("terminate task result: %v", terminateTask)
```


## License

The GOLANG SWAN SDK is released under the **MIT** license, details of which can be found in the LICENSE file.
