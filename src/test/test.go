package test

type Test interface {
	Correctness(opts Options)
	Perf(num int, outputFilename string, opts Options)
}
