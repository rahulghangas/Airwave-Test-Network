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
	"awt/local"
	"awt/test"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Execute remote test",
	Long: `Spin up a cluster of local nodes

Any metrics calculated are collated for a single node and are output
to a file called 'output.txt'' or a filename defined by usage of the flag 
'-o' or '--output'

Execute correctness/performance tests using the flags --correctness and --perf

Define number of nodes in cluster for performance test using the flag --num or -n`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			panic("Acquiring `name` flag failed")
		}
		correctness, err := cmd.Flags().GetBool("correctness")
		if err != nil {
			panic("Acquiring `correctness` flag failed")
		}
		perf, err := cmd.Flags().GetBool("perf")
		if err != nil {
			panic("Acquiring `perf` flag failed")
		}
		outputFilename, err := cmd.Flags().GetString("output")
		if err != nil {
			panic("Acquiring `output` flag failed")
		}
		num, err := cmd.Flags().GetInt("num")
		if err != nil {
			panic("Acquiring `num` flag failed")
		}
		ct, err := cmd.Flags().GetInt("clientTimeout")
		if err != nil {
			panic("Acquiring `clientTimeout` flag failed")
		}
		tt, err := cmd.Flags().GetInt("transportTimeout")
		if err != nil {
			panic("Acquiring `transportTimeout` flag failed")
		}
		got, err := cmd.Flags().GetInt("gossiperOptsTimeout")
		if err != nil {
			panic("Acquiring `gossiperOptsTimeout` flag failed")
		}
		gt, err := cmd.Flags().GetInt("gossiperTimeout")
		if err != nil {
			panic("Acquiring `gossiperTimeout` flag failed")
		}
		gwt, err := cmd.Flags().GetInt("gossiperWiggleTimeout")
		if err != nil {
			panic("Acquiring `gossiperWiggleTimeout` flag failed")
		}
		st, err := cmd.Flags().GetInt("syncerTimeout")
		if err != nil {
			panic("Acquiring `syncerTimeout` flag failed")
		}
		swt, err := cmd.Flags().GetInt("syncerWiggleTimeout")
		if err != nil {
			panic("Acquiring `syncerWiggletimeout` flag failed")
		}
		ot, err := cmd.Flags().GetInt("oncePoolTimeout")
		if err != nil {
			panic("Acquiring `oncePoolTimeout` flag failed")
		}

		testOptions := test.Options{ClientTimeout: ct,
			TransportTimeout:      tt,
			GossiperOptsTimeout:   got,
			GossiperTimeout:       gt,
			GossiperWiggleTimeout: gwt,
			SyncerTimeout:         st,
			SyncerWiggleTimeout:   swt,
			OncePoolTimeout:       ot}

		err = local.Run(name, outputFilename, num, correctness, perf, testOptions)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	testCmd.AddCommand(localCmd)

	localCmd.Flags().IntP("num", "n", 10,
		"number of nodes to spin up for performance tests")
	localCmd.Flags().Int("clientTimeout", 5, "Timeout for channel client")
	localCmd.Flags().Int("transportTimeout", 1000, "Running time for transport layer")
	localCmd.Flags().Int("gossiperOptsTimeout",  2, "Inner timeout defined in gossiper options")
	localCmd.Flags().Int("gossiperTimeout", 2, "Timeout for call to Gossip")
	localCmd.Flags().Int("gossiperWiggleTimeout", 2, "Wait period for a response to a sync request")
	localCmd.Flags().Int("syncerTimeout", 2, "Timeout for call to Sync")
	localCmd.Flags().Int("syncerWiggleTimeout",  2, "Wait period for a response to a sync request")
	localCmd.Flags().Int("oncePoolTimeout", 10, "Timeout for a OncePool persistent connection")
}
