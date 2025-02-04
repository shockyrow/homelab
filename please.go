package main

import (
	"fmt"
	"os"
)

type Action func(args []string) int

const ACTION_HELP string = "help"

var actions map[string]Action = map[string]Action{
	ACTION_HELP: func(args []string) int {
		fmt.Println("Coming soon...")
		return 0
	},
	"start": func(args []string) int {
		fmt.Println("Not implemented!")
		return 1
	},
}

func main() {
	args := os.Args

	if len(args) < 2 {
		run(ACTION_HELP, []string{})
	}

	run(args[1], args[2:])
}

func run(actionName string, args []string) {
	action, exists := actions[actionName]

	if !exists {
		action = actions[ACTION_HELP]
	}

	os.Exit(action(args))
}
