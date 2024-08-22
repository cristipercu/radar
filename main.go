package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Println("Usage: radar <command> [subcommand] [flags]")
    os.Exit(1)
  }

  command := os.Args[1]

  switch command {
  case "sync":
    handleSync(os.Args[2:])
  case "--help":
    handleHelp()
  
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

    Flags:

      -dirname  Directory name for the config. 
                If no config name is provided, I will create .radar for you. 
                Note that if you create your own dir, you will need to always specify it

  --help  Display this help message
`)
}


type SyncConfig struct {
  ServerAddres string `json:"server_address"`
  LocalPath string `json:"local_path"`
  RemotePath string `json:"remote_path"`
  User string `json:"user"`
  KeyPath string `json:"key_path"`
  Exclude []string `json:"exlude"`
}



func handleSync(args []string) {
  fmt.Printf("Args: %v", args)
  subcommand := args[0]

  configDirName := flag.String("dirname", "", "Directory name for the config. If no config name is provided, I will create .radar for you. Note that if you create your own dir, you will need to allways specify it" )
  flag.Parse()

  fmt.Println("path", *configDirName)

  if len(*configDirName) == 0 {
    *configDirName = ".radar"
  }

  configFileName := "conf.json"


  switch subcommand {
  case "create-config":
    err := createBaseConfig(*configDirName, configFileName)
    // TODO: Check if the config file exists and skip or ask for override
    if err != nil {
      fmt.Println("Error creating the config", err)
      os.Exit(1)
    }
  case "push":
    config := readBaseConfig(*configDirName, configFileName)
    command := createRsyncCommand(config, subcommand)
    fmt.Printf("Rsync command: %v \n", command)
    args = strings.Fields(command)
    cmd := exec.Command(args[0], args[1:]...)

    output, err := cmd.Output()
    if err != nil {
      fmt.Printf("Error running the sync command %v", err)
      os.Exit(1)
    }

    fmt.Println(output)
    
  default:
    fmt.Println("Invalid command", subcommand)
    os.Exit(1)
  }

}

func createBaseConfig(configDirName string, configFileName string) error {
  cwd, err := os.Getwd()
  if err != nil {
    return err
  }
  config := SyncConfig{LocalPath: cwd}
  configFilePath := filepath.Join(cwd, configDirName, configFileName)
  err = os.MkdirAll(filepath.Dir(configFilePath), 0755)
  if err != nil {
    return err
  }
 
  jsonData, err := json.MarshalIndent(config, "", " ")
  if err != nil {
    return err
  }
  // TODO: check if file exists and return error if exists
  err = os.WriteFile(configFilePath, jsonData, 0644)
  if err != nil {
    return err
  }
 
  return nil
}

func readBaseConfig(configDirName string, configFileName string) SyncConfig {
  cwd, err := os.Getwd()
  if err != nil {
    fmt.Println("Could not read the current working dir, maybe we do not have access", err)
    os.Exit(1)
  }
  configPath := filepath.Join(cwd, configDirName, configFileName)
  file, err := os.Open(configPath)
  if err != nil {
    fmt.Println("could not read config file, make sure that the config dir (ex: .radar) is in the cwd")
    os.Exit(1)
  }
  defer file.Close()

  var config SyncConfig 
  decoder := json.NewDecoder(file)
  err = decoder.Decode(&config)
  if err != nil {
    fmt.Println("could not decode the json file, maybe the structure is not ok", err)
    os.Exit(1)
  }

  // fmt.Printf("JSON: %v", config)

  return config

}

func createRsyncCommand(config SyncConfig, commandType string) string {
  command := "rsync -avz 'ssh"

  if config.ServerAddres == "" ||
     config.LocalPath == "" ||
     config.RemotePath == "" ||
     config.User == "" {
    fmt.Println("server_address, user, local_path, remote_path are necessary, please update config file")
    os.Exit(1)
    
  }

  if config.KeyPath != "" {
    command += " -i " + config.KeyPath + "'"
  } else {
    command += "' "
  }

  if len(config.Exclude) != 0 {
    for _, value := range config.Exclude {
      command += " --exclude='" + value +"'"
    }
  }

  remoteFullAdress := config.User + "@" + config.ServerAddres + ":" + config.RemotePath

  switch commandType {
  case "push":
    command += " " + config.LocalPath + " " + remoteFullAdress 
  case "pull":
    command += " " + remoteFullAdress + " " + config.LocalPath
  default:
    fmt.Println("500 Internal server error, no avaialble command push|pull")
    os.Exit(1)
  }


  return command 
}
