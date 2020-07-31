/*
Copyright © 2020 Jonathan Womack <jonathan@jwomack.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"database/sql"
	"fmt"
	"github.com/kpwn243/twelve-factor-app-env-buddy/internal"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfa",
	Short: "",
	Long: ``,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("failed to get home directory")
		os.Exit(1)
	}
	tfaDir := home + "/.tfa"
	dbFile := tfaDir + "/db.sqlite"
	tfaShell := tfaDir + "/tfa.sh"

	if _, err := os.Stat(tfaDir); os.IsNotExist(err) {
		fmt.Println("~/.tfa directory not found. Creating")
		err := os.Mkdir(tfaDir, os.FileMode(0755))
		if err != nil {
			fmt.Println("Failed to create ~/.tfa directory. Exiting")
			os.Exit(1)
		}
		f, err := os.Create(tfaShell)
		if err != nil {
			fmt.Println("Failed to create tfa.sh file. Exiting")
			os.Exit(1)
		}
		f.Close()

		rcFile, err := os.OpenFile(home + "/.zshrc", os.O_APPEND|os.O_WRONLY, 0644)
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

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Println("Failed to open database connection. Exiting")
		os.Exit(1)
	}
	_, err = db.Exec(internal.DbInit)
	if err != nil {
		fmt.Println("Failed to create app database tables. Exiting", err)
		os.Exit(1)
	}


	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tfa.yml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//
//		// Search config in home directory with name ".tfa" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".tfa")
//	}
//
//	viper.AutomaticEnv() // read in environment variables that match
//
//	// If a config file is found, read it in.
//	if err := viper.ReadInConfig(); err == nil {
//		fmt.Println("Using config file:", viper.ConfigFileUsed())
//	}
//}
