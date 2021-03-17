package local

import (
	"awt/test"
	"fmt"
)

type PeerDiscoveryTest struct {}
var _ test.Test = &PeerDiscoveryTest{}

func (st *PeerDiscoveryTest) Correctness(testOpts test.Options) {
	panic(fmt.Sprintf("No corretness tests for %v test", name))
}

func (st *PeerDiscoveryTest) Perf(num int, topology test.Topology, outputName string, testOpts test.Options) {
	panic(fmt.Sprintf("No perf tests for %v test", name))
}
