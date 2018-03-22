package main

import (
	"fmt"
	"os/exec"
	"bytes"
	"strings"
	"os"
	"bufio"
	"io"
)

// StopMiner kill minging process
func StopMiner() {
	var killcommand = "taskkill /f /im " + minerprocess
	c := exec.Command("cmd", "/c", killcommand)
	c.Run()
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

func minerIsRunning() bool {
	_, err := getMinerPid()
	return err == nil
}

func resetEmptyMinerLog() {
	err := os.Truncate("miner.log", 0)
	if err != nil {
		fmt.Println("reset miner log fail")
	} else {
		fmt.Println("reset miner log success")
	}
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

	var fullcommand = getCurrentDirectory("/") + "/" + config.ZcashMinerDir + "/" + minerprocess + " --server " + config.DefaultPoolAddress + " --user " + config.DefaultZcashWalletAddress + "." + config.LocalMachineName + " --port " + config.DefaultPoolPort + " " + config.ZcashMinerFlagExtra + " --log 2 --api"
	//fmt.Println(fullcommand)
	// fmt.Print("Process ID = ")
	// fmt.Println(config.ProcessID)
	// fmt.Print("Type of Process ID = ")
	// fmt.Println(reflect.TypeOf(config.ProcessID))
	// Check if process been ran
	if config.ProcessID == 0 {
		c := exec.Command("cmd", "/C", fullcommand)
		if err := c.Run(); err != nil {
			fmt.Println("miner exited ", err)
		}
		config.SetPID(c.Process.Pid)
		// config.ProcessID = c.Process.Pid
	}
}
