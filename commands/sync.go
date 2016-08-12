package commands

import (
	"log"
)

// Sync subcommand
var Sync = &Command{
	Run:       runSync,
	UsageLine: "sync",
	Short:     "synchronize the brewery",
	Long:      `Synchronize the brewery.`,
}

func runSync(cmd *Command, args []string) {
	log.Println("syncing")
}
