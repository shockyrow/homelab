package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ActionResult uint8
type ActionFunction func(args []string) ActionResult
type Action struct {
	description string
	usage       string
	action      ActionFunction
}

const RESULT_SUCCESS ActionResult = 0
const RESULT_FAILURE ActionResult = 1
const RESULT_INVALID_USAGE ActionResult = 255

const COLS_GAP int = 5

var actions map[string]Action = map[string]Action{
	"start": {
		description: "Starts given stack",
		usage:       "start <stack_name> [options]",
		action: func(args []string) ActionResult {
			if len(args) < 1 {
				return RESULT_INVALID_USAGE
			}

			return runCommand("./stacks/"+args[0], "docker compose up -d --remove-orphans", nil, nil)
		},
	},
	"restart": {
		description: "Restarts given service of the given stack",
		usage:       "restart <stack_name> <service_name> [anoter_service_names]",
		action: func(args []string) ActionResult {
			if len(args) < 2 {
				return RESULT_INVALID_USAGE
			}

			actionResult := RESULT_SUCCESS

			serviceNames := args[1:]
			longestNameLen := 0

			for _, serviceName := range serviceNames {
				if len(serviceName) > longestNameLen {
					longestNameLen = len(serviceName)
				}
			}

			format := fmt.Sprintf("%%-%ds %%s\n", longestNameLen+COLS_GAP)

			for i, serviceName := range serviceNames {
				commandResult := runCommand("./stacks/"+args[0], fmt.Sprintf("docker compose restart %s", serviceName), nil, nil)
				commandResultString := "Done"

				if commandResult != RESULT_SUCCESS {
					actionResult = RESULT_FAILURE
					commandResultString = "Failed"
				}

				if i == 0 {
					fmt.Printf(format, "Service", "Status")
				}

				fmt.Printf(format, serviceName, commandResultString)
			}

			return actionResult
		},
	},
}

func main() {
	args := os.Args

	if len(args) < 2 {
		showHelp()
	}

	run(args[1], args[2:])
}

func showHelp() {
	rows := [][]string{{"Action", "Description", "Usage"}}

	for name, action := range actions {
		rows = append(rows, []string{name, action.description, action.usage})
	}

	prettyPrintTable(rows, nil)

	os.Exit(0)
}

func run(actionName string, args []string) {
	actionData, exists := actions[actionName]

	if !exists {
		showHelp()
	}

	code := actionData.action(args)

	if code == RESULT_INVALID_USAGE {
		fmt.Println("Invalid usage! Usage: ", actionData.usage)
	} else if code != RESULT_SUCCESS {
		fmt.Println("Something went wrong! Code: ", code)
	}

	os.Exit(int(code))
}

func runCommand(workingDir, command string, outputWriter, errorWriter io.Writer) ActionResult {
	if outputWriter == nil {
		outputWriter = os.Stdout
	}

	if errorWriter == nil {
		errorWriter = os.Stderr
	}

	commandParts := strings.Split(command, " ")

	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	cmd.Dir = filepath.Join(strings.Split(workingDir, "/")...)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Fprintf(errorWriter, "Error creating StdoutPipe: %v\n", err)
		return RESULT_FAILURE
	}

	scanner := bufio.NewScanner(stdout)

	cmd.Start()

	for scanner.Scan() {
		fmt.Fprintln(outputWriter, scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(errorWriter, "Error waiting for command: %v\n", err)
		return RESULT_FAILURE
	}

	return RESULT_SUCCESS
}

func prettyPrintTable(rows [][]string, outputWriter io.Writer) {
	if outputWriter == nil {
		outputWriter = os.Stdout
	}

	colSizes := make([]int, len(rows[0]))

	for _, row := range rows {
		for i, cell := range row {
			colSizes[i] = max(colSizes[i], len(cell))
		}
	}

	for _, row := range rows {
		for i, cell := range row {
			format := fmt.Sprintf("%%-%ds", colSizes[i]+COLS_GAP)

			if i == len(row)-1 {
				format = "%s"
			}

			fmt.Fprintf(outputWriter, format, cell)
		}

		fmt.Println()
	}
}
