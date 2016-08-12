package commands

import (
	"log"
)

// Delete subcommand
var Delete = &Command{
	Run:       runDelete,
	UsageLine: "delete",
	Short:     "delete a work unit",
	Long:      `Delete a work unit to the brewery.`,
}

func runDelete(cmd *Command, args []string) {
	log.Println("delete")
}
