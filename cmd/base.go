/*
 * queenie - a spelling bee helper
 * Copyright (C) 2022 Michael D Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var globalBase struct {
	TestFlag    bool
	VerboseFlag bool
	ConfigFile  string // configuration file from command line flag

	envPrefix  string // value to prepend when converting flags to env variables
	cfgName    string // default configuration file name
	homeFolder string // derived path to home directory
}

// cmdBase represents the base command when called without any subcommands
var cmdBase = &cobra.Command{
	Use:   "queenie",
	Short: "Spelling Bee helper",
	Long:  `queenie provides services to help solve the Spelling Bee.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// now bind viper and cobra configuration since this hook always runs early
		return bindConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("%-20s == %q\n", "HOME", globalBase.homeFolder)
		log.Printf("%-20s == %q\n", "QUEENIE_RC", viper.ConfigFileUsed())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(cmdBase.Execute())
}

func init() {
	// set the env and config
	globalBase.envPrefix, globalBase.cfgName = "QUEENIE", ".queenie"
	// find home directory
	if home, err := homedir.Dir(); err == nil {
		globalBase.homeFolder = home
	}

	cmdBase.PersistentFlags().StringVar(&globalBase.ConfigFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s.json)", globalBase.cfgName))
	cmdBase.PersistentFlags().BoolVar(&globalBase.TestFlag, "test", false, "test mode")
	cmdBase.PersistentFlags().BoolVar(&globalBase.VerboseFlag, "verbose", false, "verbose mode")

	//// Cobra also supports local flags, which will only run when this action is called directly.
	//cmdBase.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
