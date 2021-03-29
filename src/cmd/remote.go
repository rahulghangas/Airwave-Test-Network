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
	"awt/remote"
	"awt/test"
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
		topology, err := cmd.Flags().GetString("topology")
		if err != nil {
			panic("Acquiring `topology` flag failed")
		}

		index, err := cmd.Flags().GetInt("index")
		if err != nil {
			panic("Acquiring `index` flag failed")
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

		remote.Run(index, test.Topology(topology), outputFilename, correctness, perf, testOptions)
	},
}

func init() {
	testCmd.AddCommand(remoteCmd)

	remoteCmd.Flags().Int("index", -1, "Node index")
	remoteCmd.Flags().Int("ct", 5, "Timeout for channel client")
	remoteCmd.Flags().Int("tt", 1000, "Running time for transport layer")
	remoteCmd.Flags().Int("got",  2, "Inner timeout defined in gossiper options")
	remoteCmd.Flags().Int("gt", 2, "Timeout for call to Gossip")
	remoteCmd.Flags().Int("st", 2, "Timeout for call to Sync")
	remoteCmd.Flags().Int("swt",  2, "Wait period for a response to a sync request")
	remoteCmd.Flags().Int("opt", 10, "Timeout for a OncePool persistent connection")
}
