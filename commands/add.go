package commands

import (
	"log"
)

// Add subcommand
var Add = &Command{
	Run:       runAdd,
	UsageLine: "add",
	Short:     "add new work unit",
	Long:      `Add a new work unit to the brewery.`,
}

func runAdd(cmd *Command, args []string) {
	log.Println("delete")
}
