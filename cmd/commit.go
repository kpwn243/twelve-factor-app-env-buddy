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
	"fmt"
	"github.com/kpwn243/twelve-factor-app-env-buddy/internal"
	"github.com/spf13/cobra"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Write tfa database records to the tfa shell file.",
	Run: func(cmd *cobra.Command, args []string) {
		config := internal.InitConfiguration()
		db, err := internal.InitDbConnection()
		if err != nil {
			fmt.Println("Failed to create database connection. Exiting")
			os.Exit(1)
		}
		var osVars strings.Builder

		rows, _ := db.Query(internal.SelectAppValues)
		var (
			appHeaderWritten bool
			currentAppName   string
			currentAppEnv    string
			appName          string
			appEnv           string
			varName          string
			varValue         string
		)
		for rows.Next() {
			err = rows.Scan(&appName, &appEnv, &varName, &varValue)
			if err != nil {
				fmt.Println("Failed to read data from the tfa database. Exiting")
				os.Exit(1)
			}
			if currentAppName != appName && currentAppEnv != appEnv {
				if appHeaderWritten {
					_, err = fmt.Fprintf(&osVars, "\n")
					if err != nil {
						fmt.Println("Failed writing os variables header. Exiting")
						os.Exit(1)
					}
				}
				_, err = fmt.Fprintf(&osVars, "## %s environment variables for %s\n", appEnv, appName)
				if err != nil {
					fmt.Println("Failed writing os variables header. Exiting")
					os.Exit(1)
				}
				appHeaderWritten = true
			}
			_, err := fmt.Fprintf(&osVars, "export %s_%s_%s=%s\n", strings.ToUpper(appName), strings.ToUpper(appEnv), strings.ToUpper(varName), varValue)
			if err != nil {
				fmt.Println("Failed writing os variables lines. Exiting")
				os.Exit(1)
			}
		}

		err = os.Truncate(config.TfaShellFileLocation, 0)
		if err != nil {
			fmt.Println("Failed to truncate the tfa shell file. Exiting")
			os.Exit(1)
		}

		fileContent := osVars.String()
		_, err = config.TfaShellFile.WriteAt([]byte(fileContent), 0)
		if err != nil {
			fmt.Println("Failed writing os variables file. Exiting")
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
