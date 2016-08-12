package commands

import (
	"log"
)

// Edit subcommand
var Edit = &Command{
	Run:       runEdit,
	UsageLine: "edit",
	Short:     "edit a work unit",
	Long:      `Edit a work unit to the brewery.`,
}

func runEdit(cmd *Command, args []string) {
	log.Println("edit")
}
