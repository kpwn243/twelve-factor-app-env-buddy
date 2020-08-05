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
	"errors"
	"fmt"
	"github.com/kpwn243/twelve-factor-app-env-buddy/internal"
	"github.com/spf13/cobra"
	"os"
	"strings"
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
		db := internal.GetDbConnection()

		appName := strings.ToUpper(args[0])
		appEnv := strings.ToUpper(args[1])

		var app internal.Application

		if db.Where("name=?", appName).First(&app).RecordNotFound() {
			db.Save(&internal.Application{
				Name:         appName,
				Active:       true,
				Environments: []internal.Environment{{
					EnvName: appEnv,
					WriteWithoutPrefix: false,
					Active: true,
				}},
			})
			fmt.Println(fmt.Sprintf("Created %s environment for application %s.", appEnv, appName))
			os.Exit(0)
		}
		
		var env internal.Environment
		
		if db.Where("env_name=?", appEnv).First(&env).RecordNotFound() {
			db.Create(&internal.Environment{
				ApplicationId:      app.ID,
				EnvName:            appEnv,
				WriteWithoutPrefix: false,
				Active:             true,
			})
			fmt.Println(fmt.Sprintf("Created %s environment for application %s.", appEnv, appName))
			os.Exit(0)
		}

		fmt.Println(fmt.Sprintf("%s %s already exists in the database", appName, appEnv))
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
