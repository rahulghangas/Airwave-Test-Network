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
		topology, err := cmd.Flags().GetString("topology")
		if err != nil {
			panic("Acquiring `topology` flag failed")
		}

		ct, err := cmd.Flags().GetInt("ct")
		if err != nil {
			panic("Acquiring `clientTimeout` flag failed")
		}
		tt, err := cmd.Flags().GetInt("tt")
		if err != nil {
			panic("Acquiring `transportTimeout` flag failed")
		}
		got, err := cmd.Flags().GetInt("got")
		if err != nil {
			panic("Acquiring `gossiperOptsTimeout` flag failed")
		}
		gt, err := cmd.Flags().GetInt("gt")
		if err != nil {
			panic("Acquiring `gossiperTimeout` flag failed")
		}
		st, err := cmd.Flags().GetInt("st")
		if err != nil {
			panic("Acquiring `syncerTimeout` flag failed")
		}
		swt, err := cmd.Flags().GetInt("swt")
		if err != nil {
			panic("Acquiring `syncerWiggletimeout` flag failed")
		}
		ot, err := cmd.Flags().GetInt("opt")
		if err != nil {
			panic("Acquiring `oncePoolTimeout` flag failed")
		}

		testOptions := test.Options{ClientTimeout: ct,
			TransportTimeout:      tt,
			GossiperOptsTimeout:   got,
			GossiperTimeout:       gt,
			SyncerTimeout:         st,
			SyncerWiggleTimeout:   swt,
			OncePoolTimeout:       ot,
		}

		err = local.Run(name, test.Topology(topology), outputFilename, num, correctness, perf, testOptions)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	testCmd.AddCommand(localCmd)

	localCmd.Flags().IntP("num", "n", 10,
		"number of nodes to spin up for performance tests")
	localCmd.Flags().Int("ct", 5, "Timeout for channel client")
	localCmd.Flags().Int("tt", 1000, "Running time for transport layer")
	localCmd.Flags().Int("got",  2, "Inner timeout defined in gossiper options")
	localCmd.Flags().Int("gt", 2, "Timeout for call to Gossip")
	localCmd.Flags().Int("st", 2, "Timeout for call to Sync")
	localCmd.Flags().Int("swt",  2, "Wait period for a response to a sync request")
	localCmd.Flags().Int("opt", 10, "Timeout for a OncePool persistent connection")
}
