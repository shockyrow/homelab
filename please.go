package main

import (
	"fmt"
	"os"
)

type ActionResult uint8

const RESULT_SUCCESS ActionResult = 0
const RESULT_INVALID_USAGE ActionResult = 255

type ActionFunction func(args []string) ActionResult
type Action struct {
	description string
	usage       string
	action      ActionFunction
}

var actions map[string]Action = map[string]Action{
	"start": {
		description: "Starts given stack",
		usage:       "start <stack_name> [options]",
		action: func(args []string) ActionResult {
			if len(args) == 0 {
				return RESULT_INVALID_USAGE
			}

			return RESULT_SUCCESS
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
	fmt.Printf("%-20s\t%-40s\t%-20s\n", "Action", "Description", "Usage")

	for name, action := range actions {
		fmt.Printf("%-20s\t%-40s\t%-20s\n", name, action.description, action.usage)
	}

	os.Exit(0)
}

func run(actionName string, args []string) {
	actionData, exists := actions[actionName]

	if !exists {
		showHelp()
	}

	code := actionData.action(args)

	if code == RESULT_INVALID_USAGE {
		fmt.Println("Usage: ", actionData.usage)
	}

	os.Exit(int(code))
}
