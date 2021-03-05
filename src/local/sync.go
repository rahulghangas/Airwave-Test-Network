package local

import "awt/test"

type SyncTest struct {}
var _ test.Test = &SyncTest{}

func (st *SyncTest) Correctness(testOpts test.Options) {

}

func (st *SyncTest) Perf(num int, outputName string, testOpts test.Options) {

}
