package internal

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
)

type Configuration struct {
	HomeDirectory        string
	TfaDirectory         string
	TfaShellFileLocation string
	DbFileLocation       string
	TfaShellFile         *os.File
}

var config *Configuration

func InitConfiguration() *Configuration {
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println("failed to get home directory")
		os.Exit(1)
	}
	tfaDir := homeDir + "/.tfa"
	tfaShellPath := tfaDir + "/tfa.sh"
	dbFilePath := tfaDir + "/db.sqlite"

	initFiles(tfaDir, tfaShellPath, homeDir)

	tfaShellFile, err := os.OpenFile(tfaShellPath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Failed to open tfa.sh file. Exiting")
		os.Exit(1)
	}

	config = &Configuration{
		HomeDirectory:        homeDir,
		TfaDirectory:         tfaDir,
		TfaShellFileLocation: tfaShellPath,
		DbFileLocation:       dbFilePath,
		TfaShellFile:         tfaShellFile,
	}

	return config
}

func GetConfiguration() *Configuration {
	return config
}

func initFiles(tfaDirectory string, tfaShellFileLocation string, homeDir string) {
	if _, err := os.Stat(tfaDirectory); os.IsNotExist(err) {
		fmt.Println("~/.tfa directory not found. Creating")
		err := os.Mkdir(tfaDirectory, os.FileMode(0755))
		if err != nil {
			fmt.Println("Failed to create ~/.tfa directory. Exiting")
			os.Exit(1)
		}
		f, err := os.Create(tfaShellFileLocation)
		if err != nil {
			fmt.Println("Failed to create tfa.sh file. Exiting")
			os.Exit(1)
		}
		f.Close()

		rcFile, err := os.OpenFile(homeDir + "/.zshrc", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Failed to open .zshrc file. Exiting")
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		_, err = rcFile.WriteString("\n## Twelve Factor App Shell File\nsource ~/.tfa/tfa.sh")
		if err != nil {
			fmt.Println("Failed to append to .zshrc file. Exiting")
			fmt.Println(err)
			os.Exit(1)
		}
		rcFile.Sync()
		fmt.Println("Appended to .zshrc")
	}
}
