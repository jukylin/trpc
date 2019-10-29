// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/jukylin/trpc/rpc"
	"github.com/spf13/cobra"
	"time"
)

//var cfgFile string

var url string
var fn string
var fm bool
var args []string
var bench bool
var Nrun int
var Ncon int
var Dur time.Duration
var Type string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "trpc",
	Short: "RPC 调试工具",
	Long: `RPC 调试工具，用于调试远程RPC接口，暂只支持yar和Hprose，HTTP协议
trpc -u URL -f function -a param1 -a param2
`,
// Uncomment the following line if your bare application
// has an action associated with it:
	Run: func(cmd *cobra.Command, para []string) {
		if len(url) == 0 || len(fn) == 0{
			cmd.Help()
			return
		}

		args := rpc.RpcArgs{
			Type:Type,
			Url:url,
			Fn:fn,
			Format:fm,
			Bench:bench,
			Nrun:Nrun,
			Ncon:Ncon,
			Dur: Dur,
			Args:args,
		}
		rpc.DebugStart(&args)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.derpc.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().StringVarP(&Type, "Type", "t", "yar", "rpc类型 yar和hprose")
	RootCmd.Flags().StringVarP(&url, "url", "u", "", "请求地址")
	RootCmd.Flags().StringVarP(&fn, "func", "f", "", "调用的函数")
	RootCmd.Flags().BoolVarP(&fm, "format", "m", false, "是否格式化结果，主要针对map，便于使用者查看。")
	RootCmd.Flags().BoolVarP(&bench, "bench", "b", false, "进行压力测试，测试工具使用" +
		"【https://github.com/rakyll/hey，使用go重写的ab压力测试工具】。")
	RootCmd.Flags().IntVarP(&Nrun, "Nrun", "n", 200, "总的请求数，默认200")
	RootCmd.Flags().IntVarP(&Ncon, "Ncon", "c", 50, "总的压力数，不能低于总的请求数，默认50")
	RootCmd.Flags().DurationVarP(&Dur, "Dur", "z", time.Duration(10 * time.Second), "")

	RootCmd.Flags().StringArrayVarP(&args, "args", "a", []string{}, "函数参数，按顺序传")
}

// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//	if cfgFile != "" { // enable ability to specify config file via flag
//		viper.SetConfigFile(cfgFile)
//	}

//	viper.SetConfigName(".derpc") // name of config file (without extension)
//	viper.AddConfigPath("$HOME")  // adding home directory as first search path
//	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
//	if err := viper.ReadInConfig(); err == nil {
//		fmt.Println("Using config file:", viper.ConfigFileUsed())
//	}
//}
