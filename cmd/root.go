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
	"github.com/spf13/cobra"
	"os"
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
	config := internal.InitConfiguration()
	db, err := internal.InitDbConnection()
	if err != nil {
		fmt.Println("Failed to create database connection. Exiting")
		os.Exit(1)
	}

	db, err = sql.Open("sqlite3", config.DbFileLocation)
	if err != nil {
		fmt.Println("Failed to open database connection. Exiting")
		os.Exit(1)
	}
	_, err = db.Exec(internal.DbInit)
	if err != nil {
		fmt.Println("Failed to create app database tables. Exiting", err)
		os.Exit(1)
	}
}
