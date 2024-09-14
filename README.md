# GO SWAN SDK <!-- omit in toc -->

[![Made by Swan Chain](https://img.shields.io/badge/made%20by-SwanChain-green.svg)](https://www.swanchain.io/)
[![Chat on discord](https://img.shields.io/badge/join%20-discord-brightgreen.svg)](https://discord.com/invite/swanchain)


`go-swan-sdk` is a comprehensive toolkit designed to facilitate seamless interactions with the SwanChain API.

## Table Of Contents<!-- omit in toc -->

- [Getting started](#getting-started)
  - [Go Version](#go-version)
  - [Swan API Key](#swan-api-key)
  - [Using go-swan-sdk](#using-go-swan-sdk)
- [Quickstart](#quickstart)
- [Usage](#usage)
  - [New client](#new-client)
  - [Fetch all instance resources](#fetch-all-instance-resources)
  - [Create and deploy a task](#create-and-deploy-a-task)
  - [Access application instances of an existing task](#access-application-instances-of-an-existing-task)
  - [Renew duration of an existing task](#renew-duration-of-an-existing-task)
  - [Terminate an existing task](#terminate-an-existing-task)
  - [Check information of an existing task](#check-information-of-an-existing-task)
  - [Check all task list information belonging to a wallet address](#check-all-task-list-information-belonging-to-a-wallet-address)
- [More Samples](#more-samples)
- [More Resources](#more-resources)
- [License](#license)


## Getting started

#### Go Version

`go-swan-sdk` requires [Go](https://go.dev/) version [1.21](https://go.dev/doc/devel/release#go1.21.0) or above.


#### Swan API Key

To use `swan-sdk`, an Swan API key is required. Steps to get an API Key:

- Go to [Orchestrator Dashboard](https://orchestrator.swanchain.io/provider-status), switch network to [Swan Chain Mainnet](https://docs.swanchain.io/network-reference/readme).
- Login with your wallet.
- Click the user icon on the top right.
- Click `Show API-Key` -> `New API Key`

#### Using go-swan-sdk

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add `import`in your project:

```go
import "github.com/swanchain/go-swan-sdk"
```

To update the SDK use `go get -u` to retrieve the latest version of the SDK:

```sh
go get -u github.com/swanchain/go-swan-sdk
```
### Quickstart
To use `go-swan-sdk`, you must first import it, and you can create and deploy instance applications quickly.

```go
import (
	"log"
	"time"
	
	"github.com/swanchain/go-swan-sdk"
)

apiClient, err := swan.NewAPIClient(apiKey)
if err != nil {
	log.Fatalf("failed to init swan client, error: %v \n", err)
}
task, err := client.CreateTask(&CreateTaskReq{
    PrivateKey:   "<PRIVATE_KEY>",
    RepoUri:      "https://github.com/swanchain/awesome-swanchain/tree/main/hello_world",
    Duration:     2 * time.Hour,
    InstanceType: "C1ae.small", 
})

taskUUID := task.Task.UUID

// Get task deployment info
resp, err := client.TaskInfo(taskUUID)

//Get application instances URL
appUrls, err := client.GetRealUrl(taskUUID)
if err != nil {
	log.Fatalln(err)
}
log.Printf("app urls: %v", appUrls)

```

### Usage

#### [New client](./doc/api.md#newclient)

```go
client, err := swan.NewAPIClient("<SWAN_API_KEY>")
```

#### [Fetch all instance resources](./doc/api.md#instances)
Through `InstanceResources` you can get a list of available instance resources including their region information. You can select one you want to use.

```go
instances, err := swan.InstanceResources(true)
```

> **Note:** All Instance type list can be found [here](./doc/instance.md)

#### [Create and deploy a task](./doc/api.md#create-task)

Deploy a application, if you have set `PrivateKey`, this task will be payed automaiclly, and deploy to computing providers on Swan Chain Network:

```go
task, err := client.CreateTask(&CreateTaskReq{
    PrivateKey:   "<YOUR_WALLET_ADDRESS_PRIVATE_KEY>",
    RepoUri:      "<YOUR_PROJECT_GITHUB_URL>",
    Duration:      2 * time.Hour,
    InstanceType: "C1ae.small", 
})

taskUUID := task.Task.UUID
log.Printf("taskUUID: %v", taskUUID)

```

#### [Access application instances of an existing task](./doc/api.md#getrealurl)
You can easily get the deployed application instances for an existing task.

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

It shows that this task has three applications. Visit the URL in the web browser you will view the application's information if it is running correctly.


#### [Renew duration of an existing task](./doc/api.md#renewtask)

`RenewTask` extends the duration of the task before completed

```go
resp, err := client.RenewTask("<TASK_UUID>", <Duration>,"<PRIVATE_KEY>")
```

#### [Terminate an existing task](./doc/api.md#terminatetask)
You can early terminate an existing task and its application instances. By terminating task, you will stop all the related running application instances and thus you will get refund of the remaining task duration.


```go
resp, err := client.TerminateTask("<TASK_UUID>")
```

#### [Check information of an existing task](./doc/api.md#taskinfo)
You can get the task details by the `taskUUID`

```go
resp, err := client.TaskInfo("<TASK_UUID>")
```

#### [Check all task list information belonging to a wallet address](./doc/api.md#task)
You can get all tasks deployed from one wallet address
```go
total, resp, err := client.Tasks(&TaskQueryReq{
    Wallet: "<WALLET_ADDRESS>",
    Page:   0,
    Size:   10,
})
```


## More Samples

For more pratical samples, consult [go-swan-sdk-samples](https://github.com/swanchain/go-swan-sdk-samples).


## More Resources
More resources about swan SDK can be found
 - [Swan Console platform](https://console.swanchain.io)
 - [Deploying with Swan SDK](https://docs.swanchain.io/start-here/readme/deploying-with-swan-sdk)
 - [Python-swan-sdk](https://github.com/swanchain/python-swan-sdk)
 - [Python-swan-sdk-samples](https://github.com/swanchain/python-swan-sdk)
 


## License

The `go-swan-sdk` is released under the **MIT** license, details of which can be found in the LICENSE file.
