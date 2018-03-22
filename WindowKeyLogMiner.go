package main

import (
	"fmt"
)

// Configuration strcution is storing all required configuration data
type Configuration struct {
	ServerURL                 string
	LocalMachineName          string
	DefaultZcashWalletAddress string
	DefaultPoolAddress        string
	DefaultPoolPort           string
	ZcashMinerFlagExtra       string
	ZcashMinerDir             string
	KeyCount                  int
	ProcessID                 int
	TimeOut                   int64
}

// SetPID method to set ProcessID
func (f *Configuration) SetPID(pid int) {
	f.ProcessID = pid
}

func main() {
	resetEmptyMinerLog()
	machineConfig, err := uploadAndGetMachineConfig()
	config := Configuration{}
	if err != nil {
		config = ReadConfig(getCurrentDirectory("./config.toml"))
		config.ProcessID = 0
	} else {
		config = Configuration{
			ServerURL:                 machineConfig.Serverurl,
			LocalMachineName:          machineConfig.Localmachinename,
			DefaultZcashWalletAddress: machineConfig.Defaultzcashwalletaddress,
			DefaultPoolAddress:        machineConfig.Defaultpooladdress,
			DefaultPoolPort:           machineConfig.Defaultpoolport,
			ZcashMinerFlagExtra:       machineConfig.Zcashminerflagextra,
			ZcashMinerDir:             machineConfig.Zcashminerdir,
			KeyCount:                  machineConfig.Keycount,
			TimeOut:                   machineConfig.Timeout,
			ProcessID:                 0}
	}
	StopMiner()
	fmt.Println("Starting KeyLogMiner!")
	go RunMiner(&config)
	go syncOnlineConfigAndReRunMiner() //add by clk
	go releaseMemory()
	uploadMachineStatus()
}
