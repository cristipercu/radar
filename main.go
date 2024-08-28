package main

import (
	"fmt"
	"os"

	"github.com/cristipercu/radar/mm"
	"github.com/cristipercu/radar/sync"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Println("Usage: radar <command> [subcommand] [flags]")
    os.Exit(1)
  }

  command := os.Args[1]

  switch command {
  case "--help":
    handleHelp()
  case "sync":
    sync.HandleSync(os.Args[2:])
  case "mm":
    //TODO: change the duration to be a flag
    mm.MoveMouse(30)
  
  default:
    fmt.Println("Invalid command", command)
    os.Exit(1)
  }

} 

func handleHelp() {
    fmt.Println(`
Usage: radar <command> [subcommand] [flags]

Commands:

  sync    Synchronize files and directories

    Subcommands:

      create-config  Create a basic configuration file
      push           Push local changes to the remote server
      mm             Start the WIP... tool that does nothing, but keeps your laptop awake

    Flags:

      -dirname  Directory name for the config. 
                If no config name is provided, I will create .radar for you. 
                Note that if you create your own dir, you will need to always specify it

  --help  Display this help message`)
}


