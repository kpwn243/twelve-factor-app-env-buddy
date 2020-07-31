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
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [app name] [app env] [variable name] [variable value]",
	Short: "A brief description of your command",
	Long: ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 4 {
			return errors.New("app name, environment, variable name, and variable value are required")
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

		appName := strings.ToUpper(args[0])
		appEnv := strings.ToUpper(args[1])
		varName := strings.ToUpper(args[2])
		varValue := args[3]

		var appEnvRecordId int
		appEnvExistsStmt := db.QueryRow("SELECT e.id FROM applications a JOIN environments e ON a.id = e.APP_KEY WHERE a.APP_NAME = ? AND e.APP_ENV = ?", appName, appEnv)
		err = appEnvExistsStmt.Scan(&appEnvRecordId)
		if err != nil {
			if err == sql.ErrNoRows {
				appEnvRecordId = 0
			} else {
				fmt.Println("Failed to check if the application/env combo was present when inserting variable value. Exiting")
				os.Exit(1)
			}
		}
		if appEnvRecordId == 0 {
			fmt.Println("Unable to create variable name/value combo as the app/env combo does not exist. Please run tfa create. Exiting")
			os.Exit(1)
		}
		_, err = db.Exec("INSERT INTO variables (APP_ENV_RECORD, VAR_NAME, VAR_VALUE) VALUES (?, ?, ?)", appEnvRecordId, varName, varValue)
		if err != nil {
			fmt.Println("Failed to insert variable name/value into the database. Exiting.")
			os.Exit(1)
		}
		fmt.Println("Added variable to the database. Don't forget to commit your changes via tfa commit")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
