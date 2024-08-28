package sync

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SyncConfig struct {
  ServerAddres string `json:"server_address"`
  LocalPath string `json:"local_path"`
  RemotePath string `json:"remote_path"`
  User string `json:"user"`
  Exclude []string `json:"exlude"`
}


func HandleSync(args []string) {
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
    if err != nil {
      fmt.Println("Error creating the config", err)
      os.Exit(1)
    }
    fmt.Println(`Config file created. You can update the json config file with the necessary info.
      Note that rsync is using the openssh protocol, and if you use a private key, you need to update your ssh config file with the server info. 
      For more info check the readme from https://github.com/cristipercu/radar
      `)
  case "push":
    config := readBaseConfig(*configDirName, configFileName)
    command := createRsyncCommand(config, subcommand)
    fmt.Printf("Rsync command: %v \n", command)
    args := strings.Split(command, " ")

    cmd := exec.Command(args[0], args[1:]...)
    
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatalf("Error creating stdout pipe: %v", err)
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        log.Fatalf("Error creating stderr pipe: %v", err)
    }

    if err := cmd.Start(); err != nil {
        log.Fatalf("Error starting command: %v", err)
    }

    out, _ := io.ReadAll(stdout)
    errOut, _ := io.ReadAll(stderr)

    err = cmd.Wait()
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            log.Fatalf("Command failed with exit code %d. Stderr: %s", exitErr.ExitCode(), errOut)

        } else {
            log.Fatalf("Command failed: %v", err)
        }
    }

    if len(out) > 0 {
        fmt.Println("Output:", string(out))
    }

    fmt.Println("Command completed successfully")   

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
  if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
    err = os.WriteFile(configFilePath, jsonData, 0644)
    if err != nil {
      return err
    }
  } else {
    fmt.Printf("Config file %v already exists", configFilePath)
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
  command := "rsync -avz -e ssh"

  if config.ServerAddres == "" ||
     config.LocalPath == "" ||
     config.RemotePath == "" ||
     config.User == "" {
    fmt.Println("server_address, user, local_path, remote_path are necessary, please update config file")
    os.Exit(1)
    
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
