package local

import (
	"awt/test"
	"bytes"
	"context"
	"fmt"
	"github.com/renproject/aw/dht/dhtutil"
	"github.com/renproject/id"
	"time"
)

type SyncTest struct {}
var _ test.Test = &SyncTest{}

func syncTest1(testOpts test.Options) {
	num := 5
	opts, peers, tables, contentResolvers, _, _ := setup(num, testOpts)

	for i := range peers {
		peers[i].Resolve(context.Background(), contentResolvers[i])
	}

	createBidiFullyConnectedTopology(num, opts, peers, tables)

	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()
	for i := range peers {
		go peers[i].Run(ctx)
	}

	msgX := "Hello Nodes! I have a secret for you"
	msgY := "My name's Node 0"
	msgZ := "My secret is - I'm Batman"
	contentIDX := id.NewHash([]byte(msgX))
	contentIDY := id.NewHash([]byte(msgY))
	contentIDZ := id.NewHash([]byte(msgZ))
	contentResolvers[0].InsertContent(contentIDX[:], []byte(msgX))
	contentResolvers[0].InsertContent(contentIDY[:], []byte(msgY))
	contentResolvers[0].InsertContent(contentIDZ[:], []byte(msgZ))

	gossipCtx, gossipCancel := context.WithTimeout(ctx, 5 * time.Second)
	defer gossipCancel()
	peers[0].Gossip(gossipCtx, contentIDX[:], nil)
	time.Sleep(time.Second)

	for i := range peers {
		contentX , ok := contentResolvers[i].QueryContent(contentIDX[:])
		if !ok {
			fmt.Println(string(colorRed), "[ ] - Follow up information sync test failed")
			panic("initial gossiped content not received\n\n")
		}
		if string(contentX) != msgX {
			panic("gossipped content not synced correctly")
		}
	}
	syncCtx, syncCancel := context.WithTimeout(ctx, 20 * time.Second)
	endChan := make(chan struct{}, num-1)
	defer syncCancel()
	for i := range peers {
		if i != 0 {
			go func() {
				peerID := peers[0].ID()
				syncedMsgY, err := peers[i].Sync(syncCtx, contentIDY[:], &peerID)
				if err != nil {
					panic(fmt.Sprintf("error when syncing 2nd content: %v", err))
				}
				if string(syncedMsgY) != msgY {
					panic("synced content is incorrect")
				}
				syncedMsgZ, err := peers[i].Sync(syncCtx, contentIDZ[:], &peerID)
				if err != nil {
					panic(fmt.Sprintf("error when syncing 2nd content: %v", err))
				}
				if string(syncedMsgZ) != msgZ {
					panic("synced content is incorrect")
				}

				endChan <- struct{}{}
			}()
		}
	}

	channelCount := 0

	for channelCount != num -1 {
		select {
		case <-endChan:
			channelCount++
		case <-syncCtx.Done():
			fmt.Println(string(colorRed), "[ ] - Follow up information sync test failed")
			panic("Context cancelled before successful syncing\n\n")
		}
	}

	fmt.Println(string(colorGreen), "[x] - Follow up information sync test ran successfully\n\n")
}

func syncTest2(testOpts test.Options) {
	num := 100
	opts, peers, tables, contentResolvers, _, _ := setup(num, testOpts)

	for i := range peers {
		peers[i].Resolve(context.Background(), contentResolvers[i])
	}

	createBidiFullyConnectedTopology(num, opts, peers, tables)

	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()
	for i := range peers {
		go peers[i].Run(ctx)
	}

	syncCount := 0

	startTime := time.Now()
	for round := 0; round < 1; round++ {
		for i := range peers {
			newMsg := dhtutil.RandomContent()
			hash := id.NewHash(newMsg)
			contentResolvers[i].InsertContent(hash[:], newMsg)

			for j := range peers {
				if i != j {
					syncCtx, syncCancel := context.WithTimeout(ctx, 5 * time.Second)
					hint := peers[i].ID()
					syncedMsg, err := peers[j].Sync(syncCtx, hash[:], &hint)
					if err == nil && bytes.Equal(syncedMsg, newMsg) {
						syncCount++
					}
					syncCancel()
				}
			}
		}
	}
	endTime := time.Now()
	totalMsgs := 1 * (num - 1) * num
	totalTime := endTime.Unix() - startTime.Unix()

	if syncCount != totalMsgs {
		fmt.Println(string(colorRed), syncCount,"/",totalMsgs, " synced in ", totalTime," seconds\n\n")
	} else {
		fmt.Println(string(colorRed), "All messages (", totalMsgs, ") synced in ", totalTime," seconds\n\n")
	}
}

func (st *SyncTest) Correctness(testOpts test.Options) {
	syncTest1(testOpts)
	syncTest2(testOpts)
}

func (st *SyncTest) Perf(num int, topology test.Topology, outputName string, testOpts test.Options) {
	panic(fmt.Sprintf("No perf tests for %v test", name))
}
