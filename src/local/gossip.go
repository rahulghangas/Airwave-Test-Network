package local

import (
	"awt/test"
	"context"
	"fmt"
	"github.com/renproject/aw/dht/dhtutil"
	"github.com/renproject/aw/peer"
	"github.com/renproject/aw/wire"
	"github.com/renproject/id"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"
)

const name = "Gossip"

type GossipTest struct{}

var _ test.Test = &GossipTest{}

func (gt *GossipTest) Correctness(testOpts test.Options) {
	panic(fmt.Sprintf("No corretness tests for %v test", name))
}

func (gt *GossipTest) Perf(numPeers int, outputFilePath string, testOpts test.Options) {
	if outputFilePath == "" {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		outputFilePath = filepath.Join(basepath, "../../output/output.txt")

	}
	fo, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	opts, peers, tables, contentResolvers, _, transports := setup(numPeers, testOpts)

	for i := range opts {
		peers[i].Resolve(context.Background(), contentResolvers[i])
	}

	for i := range peers {
		ctx, cancel := context.WithTimeout(context.Background(), duration(testOpts.TransportTimeout))
		defer cancel()

		index := i
		go func() {
			transports[index].Receive(context.Background(), func() func(from id.Signatory, msg wire.Msg) error {
				var x int64 = 0
				go func() {
					//var seconds int64 = 0
					ticker := time.NewTicker(time.Second)
					defer ticker.Stop()
					for {
						select {
						case <-ticker.C:
							//seconds += 1
							//fmt.Printf("Average throughput for peer %v: %v/second\n", index, x/seconds)
							_, err = fo.WriteString(fmt.Sprintf("%d %d\n", index, x))
							atomic.StoreInt64(&x, 0)
							if err != nil {
								fmt.Printf("error writing to file: %v", err)
							}
						}
					}

				}()
				return func(from id.Signatory, msg wire.Msg) error {
					atomic.AddInt64(&x, 1)
					if err := peers[index].Syncer().DidReceiveMessage(from, msg); err != nil {
						return err
					}
					if err := peers[index].Gossiper().DidReceiveMessage(from, msg); err != nil {
						return err
					}
					return nil
				}
			}())
			transports[index].Run(ctx)
		}()
		for j := range peers {
			if i != j {
				transports[i].Link(peers[j].ID())
			}
		}

		tables[i].AddPeer(opts[(i+1)%numPeers].PrivKey.Signatory(),
			wire.NewUnsignedAddress(wire.TCP,
				fmt.Sprintf("%v:%v", "localhost", uint16(3333+i+1)), uint64(time.Now().UnixNano())))
		tables[(i+1)%numPeers].AddPeer(opts[i].PrivKey.Signatory(),
			wire.NewUnsignedAddress(wire.TCP,
				fmt.Sprintf("%v:%v", "localhost", uint16(3333+i)), uint64(time.Now().UnixNano())))
	}

	timeout := duration(testOpts.GossiperTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	for iter := 0; iter < 10000; iter++ {
		select {
		case <-ctx.Done():
			cancel()
			ctx, cancel = context.WithTimeout(context.Background(), timeout)
		default:
		}
		for i := range peers {
			msgHello := fmt.Sprintf(string(dhtutil.RandomContent()), peers[i].ID().String())
			contentID := id.NewHash([]byte(msgHello))
			contentResolvers[i].InsertContent(contentID[:], []byte(msgHello))
			peers[i].Gossip(ctx, contentID[:], &peer.DefaultSubnet)
		}
		println("Round", iter)
	}
	cancel()

	time.Sleep(5 * time.Second)
	fmt.Printf("%v\n", testOpts)
}
