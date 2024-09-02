package main

import (
	"github.com/swanchain/go-swan-sdk"
	"log"
)

func main() {
	// init sdk client
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

	// Get application instances URL
	appUrls, err := client.GetRealUrl(createTaskResp.TaskUuid)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("app urls: %v", appUrls)

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

	terminateTask, err := client.TerminateTask(taskUuid)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("terminate task result: %v", terminateTask)

}
