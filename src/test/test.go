package test

type Test interface {
	Correctness(opts Options)
	Perf(num int, topology Topology, outputFilename string, opts Options)
}
