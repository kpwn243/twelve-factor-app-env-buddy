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
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [app name] [app environment]",
	Short: "Create application/environment combo for setting variables",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("app name and environment are required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("failed to get home directory")
			os.Exit(1)
		}
		tfaDir := home + "/.tfa"
		dbFile := tfaDir + "/db.sqlite"

		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			fmt.Println("Failed to open database connection. Exiting")
			os.Exit(1)
		}
		var appCount int
		appName := strings.ToUpper(args[0])
		appEnv := strings.ToUpper(args[1])
		appExistsStatement:= db.QueryRow("SELECT COUNT(1) FROM applications WHERE APP_NAME = ?", appName)
		err = appExistsStatement.Scan(&appCount)
		if err != nil {
			fmt.Println("Failed to check for app existence when creating env var. Exiting")
			os.Exit(1)
		}
		if appCount > 0 {
			var appKey int
			appKeyStatement := db.QueryRow("SELECT id FROM applications WHERE APP_NAME = ?", appName)
			err = appKeyStatement.Scan(&appKey)
			if err != nil {
				fmt.Println("Failed to get key of existing application. Exiting")
				os.Exit(1)
			}
			var envCount int
			envExistsStmt := db.QueryRow("SELECT COUNT(1) FROM environments WHERE APP_KEY = ? and APP_ENV = ?", appKey, appEnv)
			err = envExistsStmt.Scan(&envCount)
			if err != nil {
				fmt.Println("Failed to check for env existence when adding env to app")
				os.Exit(1)
			}
			if envCount == 0 {
				_, err := db.Exec("INSERT INTO environments (APP_KEY, APP_ENV) VALUES (?, ?)", appKey, appEnv)
				if err != nil {
					fmt.Println("Failed to create insert query for new environment with existing app. Exiting")
					os.Exit(1)
				}
			} else {
				fmt.Println(fmt.Sprintf("%s %s already exists in the database", appName, appEnv))
			}
		} else {
			_, err := db.Exec("INSERT INTO applications (APP_NAME, ACTIVE) VALUES (?, 1); INSERT INTO environments (APP_KEY, APP_ENV) VALUES (last_insert_rowid(), ?);", appName, appEnv)
			if err != nil {
				fmt.Println("Failed to insert new application/env combo into the db. Exiting")
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
