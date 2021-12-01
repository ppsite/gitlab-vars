/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pyfs/gitlab-vars/gitlab"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// gitlab server info
type GitlabConfig struct {
  Server string         `mapstructure:"server"`
  PrivateToken string   `mapstructure:"private_token"`
}


type ConfigType struct {
  Gitlab GitlabConfig `json:"gitlab"`
  Vars []gitlab.ProjectVariable
}


var cfgFile string
var config ConfigType
var projectId string    // gitlab project id
var action string       // action


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitlab-vars",
	Short: "add gitlab project level vars for gitops",
	Long: "add gitlab project level vars for gitops",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
    // read config from viper
    err := viper.Unmarshal(&config)
    if err != nil {
      fmt.Println("read config error:", err)
      return
    }
    // init gitlab client
    gl := gitlab.NewGitlab(config.Gitlab.Server, "/api/v4", config.Gitlab.PrivateToken, false)

    // do action
    switch action {
      case "list":
        vars, _ := gl.ListProjectVaribles(projectId)
        for _, item := range vars.Items {
          projectVar, err := json.Marshal(item)
          if err != nil {
            fmt.Println("json marshal error: ", err)
            return
          }
          fmt.Println(string(projectVar))
        }
      case "create":
        for _, item := range config.Vars {
          createdVar, err := gl.CreateProjectVariable(projectId, &item)
          if err != nil {
            fmt.Println(err)
            return
          }
          if createdVar.Key != item.Key {
            fmt.Println("[", item.Key, "]", "may exsit")
          }
        }
      default:    // list
        vars, _ := gl.ListProjectVaribles(projectId)
        for _, item := range vars.Items {
          fmt.Println(item)
        }
    }
  },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitlab-vars.yaml)")
  rootCmd.Flags().StringVarP(&projectId, "project", "p", "", "gitlab project id, ex: 1234 (required)")
  rootCmd.Flags().StringVarP(&action, "action", "a", "list", "[list|create] support")

  rootCmd.MarkFlagRequired("project")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gitlab-vars" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gitlab-vars")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
