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
	"fmt"
	"github.com/spf13/cobra"

)

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Execute remote test(s)",
	Long: `Spins up a single node for remote execution in the cloud. 

Any metrics calculated are collated for a single node and are output
to a file called 'output.txt'' or a filename defined by usage of the flag 
'-o' or '--output'

Execute correctness/performance tests using the flags --correctness and --perf`,

	Run: func(cmd *cobra.Command, args []string) {
		s, _ := cmd.Flags().GetString("name")
		fmt.Println("remote called", s)
	},
}

func init() {
	testCmd.AddCommand(remoteCmd)
}
