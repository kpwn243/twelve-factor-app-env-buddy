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

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [app name] [app env] [variable name] [variable value]",
	Short: "Set application variable name and value to existing application/environment combo",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 4 {
			return errors.New("app name, environment, variable name, and variable value are required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		db := internal.GetDbConnection()

		appName := strings.ToUpper(args[0])
		appEnv := strings.ToUpper(args[1])
		varName := strings.ToUpper(args[2])
		varValue := args[3]

		var env internal.Environment
		envQuery := db.Joins("JOIN applications ON environments.application_id=applications.id").
			Where("applications.name=? AND environments.env_name=?", appName, appEnv).
			Find(&env)

		if envQuery.RecordNotFound() {
			fmt.Println("Unable to create variable name/value combo as the app/env combo does not exist. Please run tfa create. Exiting")
			os.Exit(1)
		}

		var variable internal.Variable

		db.Where("environment_id=? AND var_name=?", env.ID, varName).Find(&variable)
		variable.EnvironmentId = env.ID
		variable.VarName = varName
		variable.VarValue = varValue
		variable.Active = true

		db.Save(&variable)
		fmt.Println(fmt.Sprintf("Variable %s set to the database for %s(%s). Don't forget to commit your changes via tfa commit!", varName, appName, appEnv))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
