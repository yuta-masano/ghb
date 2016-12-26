package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const windows = "windows"

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "ghb",
	Short: "a utility tool to operate GitHub repository",
}

func init() {
	// cobra.OnInitialize() は cmd.Execute() の RootCmd.Execute() 内で
	// 実行される。
	// ただし、存在するサブコマンドがエラーを返さないように実行した場合しか
	// 呼び出されない。
	// サブコマンドが正常に呼び出されると、cobra.OnInitialize() を実行した後に
	// サブコマンドの内容が実行される。
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgPath, "config",
		func() (defaultCfgPath string) {
			var cfgDir string
			if runtime.GOOS == windows {
				cfgDir = filepath.Join(os.Getenv("APPDATA"), "ghb")
			} else {
				cfgDir = filepath.Join(os.Getenv("HOME"), ".config", "ghb")
			}
			return filepath.Join(cfgDir, "config.yml")
		}(),
		"path to the config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "failed in reading config file: %s\n", err)
		os.Exit(-1)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed in setting config parameters: %s\n", err)
		os.Exit(-1)
	}
	for _, param := range []string{"Apikey"} {
		if !viper.IsSet(param) {
			fmt.Fprintf(os.Stderr, "failed in reading config parameter: %s must be specified\n", param)
			os.Exit(-1)
		}
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
