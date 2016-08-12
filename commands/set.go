package commands

import (
	"log"
)

// Set subcommand
var Set = &Command{
	Run:       runSet,
	UsageLine: "set",
	Short:     "set the properties of a work unit",
	Long:      `Set the properties of a work unit`,
}

func runSet(cmd *Command, args []string) {
	log.Println("set")
}
