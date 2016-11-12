package cmdface

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	inputBufSize = 1024
)

var inputBuf []byte
var strChan chan string

var cmdLinkOfUser *commandNode
var cmdLinkOfSys *commandNode
var cmdlTail *commandNode

func init() {
	inputBuf = make([]byte, inputBufSize)

	strChan = make(chan string)
	shower(strChan)

	cmdLinkOfSys = createSysCommandLink()
	cmdlTail = nil
}

// shower open a goroutine to capture string from strChan, then print them.
func shower(strChan <-chan string) {
	go func() {
		for {
			str, ok := <-strChan
			if !ok {
				break
			}
			fmt.Print(str)
		}
		fmt.Println("testClient shower ready.")
	}()
}

// commandNode is the node of command link list of user(sys).
type commandNode struct {
	command string
	summary string
	cmdFunc func(params []string)
	next    *commandNode
}

func helpFunc(params []string) {
	Show("=================All Commands==================\n")
	// show user command information
	cmdn := cmdLinkOfUser
	for cmdn != nil {
		Show(fmt.Sprintf("%s\t%s\n", cmdn.command, cmdn.summary))
		cmdn = cmdn.next
	}

	Show("\n")
	// show sys command information
	cmdn = cmdLinkOfSys
	for cmdn != nil {
		Show(fmt.Sprintf("%s\t%s\n", cmdn.command, cmdn.summary))
		cmdn = cmdn.next
	}
}

func exitFunc(params []string) {
	os.Exit(0)
}

func createSysCommandLink() *commandNode {

	cmdn := &commandNode{
		command: "help",
		summary: "show all commands.",
		cmdFunc: helpFunc,
	}
	cmdn.next = &commandNode{
		command: "exit",
		summary: "exit command line.",
		cmdFunc: exitFunc,
	}

	return cmdn
}

// AddCommand create a command and push it to then tail of the command link list of user.
func AddCommand(command, summary string, cmdFunc func(paras []string)) {
	cmdn := &commandNode{
		command: command,
		summary: summary,
		cmdFunc: cmdFunc,
	}

	if cmdlTail == nil {
		cmdLinkOfUser = cmdn
		cmdlTail = cmdLinkOfUser
	} else {
		cmdlTail.next = cmdn
		cmdlTail = cmdlTail.next
	}
}

func findCmdNode(link *commandNode, cmd string) (resultCmd *commandNode) {
	cmdn := link
	for resultCmd == nil && cmdn != nil {
		if cmdn.command == cmd {
			resultCmd = cmdn
		}
		cmdn = cmdn.next
	}

	return resultCmd
}

// InputAndRunCommand get command line from user inputing, then run command,
// If command is not find in comand link list of user and sys, return not find error.
func InputAndRunCommand(msg string) error {
	inputStr := strings.TrimSpace(Input(msg))
	cmdlist := strings.Split(inputStr, " ")

	if cmd := findCmdNode(cmdLinkOfSys, cmdlist[0]); cmd != nil {
		cmd.cmdFunc(cmdlist[1:])
		return nil
	}

	if cmd := findCmdNode(cmdLinkOfUser, cmdlist[0]); cmd != nil {
		cmd.cmdFunc(cmdlist[1:])
		return nil
	}

	return errors.New("Command Not Find.")
}

// Show send string to strChan to print out string synchronously.
func Show(str string) {
	strChan <- str
}

// Input show msg on the terminal then get string from stdin.
func Input(msg string) string {
	Show(msg)

	n, err := os.Stdin.Read(inputBuf)
	if err != nil {
		if err == io.EOF {
			exitFunc(nil)
		}
		fmt.Printf("Gets Error: %s.", err.Error())
		os.Exit(1)
	}

	return string(inputBuf[:n])
}
