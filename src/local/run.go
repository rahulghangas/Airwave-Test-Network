package local

import (
	"awt/test"
	"fmt"
)

var tests = map[string]test.Test {
	"gossip" : &GossipTest{},
	"sync"   : &SyncTest{},
	"peerdiscovery" : &PeerDiscoveryTest{},
}

func Run(testName string, topology test.Topology, outputFilename string, num int, correctness bool, perf bool, opts test.Options) error {
	t, ok := tests[testName]
	if !ok {
		return fmt.Errorf("supplied test name doesn't have an existent test associated with it")
	}
	if correctness {
		t.Correctness(opts)
	}
	if perf {
		t.Perf(num, topology, outputFilename, opts)
	}
	return nil
}
