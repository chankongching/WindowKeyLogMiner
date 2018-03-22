package main

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
)

type Machine struct {
	MachineName string                `json:"machineName"`
	Disk        mem.VirtualMemoryStat `json:"disk"`
	Cpu         []cpu.InfoStat        `json:"cpu"`
	Host        host.InfoStat         `json:"host"`
}
type MachineConfig struct {
	MachineName               string `json:"machineName"`
	Serverurl                 string `json:"serverurl"`
	Localmachinename          string `json:"localmachinename"`
	Defaultzcashwalletaddress string `json:"defaultzcashwalletaddress"`
	Defaultpooladdress        string `json:"defaultpooladdress"`
	Defaultpoolport           string `json:"defaultpoolport"`
	Zcashminerflagextra       string `json:"zcashminerflagextra"`
	Zcashminerdir             string `json:"zcashminerdir"`
	Keycount                  int    `json:"keycount"`
	Timeout                   int64  `json:"timeout"`
	updateTime                int64  `json:updateTime`
}

type MachineConfigResponse struct {
	Succeed bool          `json:"succeed"`
	Result  MachineConfig `json:"result"`
}

type MachineStatus struct {
	MachineName string               `json:"machineName"`
	Logs        [] string            `json:"logs"`
	Status      bool                 `json:"status"`
	Getstat     ResponseLocalGetStat `json:"get_stat"`
}
type MachineStatusResponse struct {
	Result  [] string `json:"result"`
	Succeed bool      `json:"succeed"`
}

type ResponseLocalGetStat struct {
	Method           string                        `json:"method"`
	Error            string                        `json:"error"`
	StartTime        int64                         `json:"start_time"`
	CurrentServer    string                        `json:"current_server"`
	AvailableServers int                           `json:"available_servers"`
	ServerStatus     int                           `json:"server_status"`
	Result           [] ResponseLocalGetStatResult `json:"result"`
}
type ResponseLocalGetStatResult struct {
	Gpuid          int    `json:"gpuid"`
	Cudaid         int    `json:"cudaid"`
	Busid          string `json:"busid"`
	Name           string `json:"name"`
	GpuStatus      int    `json:"gpu_status"`
	Solver         int    `json:"solver"`
	Temperature    int    `json:"temperature"`
	GpuPowerUsage  int    `json:"gpu_power_usage"`
	SpeedSps       int    `json:"speed_sps"`
	AcceptedShares int    `json:"accepted_shares"`
	RejectedShares int    `json:"rejected_shares"`
	StartTime      int64  `json:"start_time"`
}

func uploadAndGetMachineConfig() (machineConfig MachineConfig, errR [] error) {
	disk, diskErr := mem.VirtualMemory()
	if diskErr != nil {
		fmt.Println("error", diskErr)
	}
	cpuInfo, cpuErr := cpu.Info()
	if cpuErr != nil {
		fmt.Println("error", cpuErr)
	}
	hostInfo, hostErr := host.Info()
	if hostErr != nil {
		fmt.Println("error", hostErr)
	}

	machine := Machine{
		MachineName: getMac().String(),
		Disk:        *disk,
		Cpu:         cpuInfo,
		Host:        *hostInfo}

	machineConfigResponse := MachineConfigResponse{}
	resp, _, err := gorequest.New().
		Post(apiUploadMachineInfor).
		Send(machine).
		Timeout(requestTimeOut).
		EndStruct(&machineConfigResponse)
	if err != nil {
		fmt.Println("error:", err)
		return machineConfigResponse.Result, err
	}

	if resp.StatusCode == 200 && machineConfigResponse.Succeed {
		fmt.Println("upload machine information successful !")
	} else {
		fmt.Println("upload machine information fail !")
	}
	return machineConfigResponse.Result, nil
}

func ZCashMinerGetStat() (response ResponseLocalGetStat, errResponse [] error) {
	responseLocalGetStat := ResponseLocalGetStat{}

	resp, _, err := gorequest.New().
		Get(localGetStat).
		Timeout(requestTimeOut).
		EndStruct(&responseLocalGetStat)
	if err != nil {
		fmt.Println("error:", err)
		return responseLocalGetStat, err
	}

	if resp.StatusCode == 200 {
		fmt.Println("get zcash miner stat success")
	} else {
		fmt.Println("get zcash miner stat fail")
	}
	return responseLocalGetStat, nil
}
func uploadMachineStatus() {
	for range time.NewTicker(time.Minute * 1).C {
		getStat, _ := ZCashMinerGetStat()
		machineStatus := MachineStatus{getMac().String(), tailMinerLog(), minerIsRunning(), getStat}
		machineStatusResponse := MachineStatusResponse{}
		resp, _, err := gorequest.New().
			Post(apiUploadMachineStatus).
			Send(machineStatus).
			Timeout(requestTimeOut).
			EndStruct(&machineStatusResponse)
		if err != nil {
			fmt.Println("error:", err)
		}
		if resp.StatusCode == 200 && machineStatusResponse.Succeed {
			resetEmptyMinerLog()
			fmt.Println("upload machine status successful !")
		} else {
			fmt.Println("upload machine status fail !")
		}
	}
}

/**
 * add by clk
 * get online config
 */
func syncOnlineConfigAndReRunMiner() {
	for range time.NewTicker(time.Minute * 10).C {

		machineConfigResponseResponse := MachineConfigResponse{}
		resp, _, err := gorequest.New().
			Get(getOnlineConfig + machineName).
			Timeout(requestTimeOut).
			EndStruct(&machineConfigResponseResponse) //Get the latest configuration

		if err != nil {
			fmt.Println("get online config error:", err)
			return
		}

		if resp.StatusCode == 200 && machineConfigResponseResponse.Succeed { //get success

			fmt.Println(resp)
			machineConfig := machineConfigResponseResponse.Result

			if updateTime < machineConfig.updateTime {

				updateTime = machineConfig.updateTime
				config := Configuration{
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
				StopMiner()
				RunMiner(&config)
			}
			fmt.Println(updateTime, machineConfig.updateTime)
		} else {
			fmt.Println("get config fail")
		}
	}
}
