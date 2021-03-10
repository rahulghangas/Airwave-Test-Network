package local

import (
	"awt/test"
	"context"
	"fmt"
	"github.com/renproject/aw/wire"
	"github.com/renproject/id"
	"time"
)

type SyncTest struct {}
var _ test.Test = &SyncTest{}

func (st *SyncTest) Correctness(testOpts test.Options) {
	num := 5
	_, peers, tables, contentResolvers, _, _ := setup(num, testOpts)

	for i := range peers {
		peers[i].Resolve(context.Background(), contentResolvers[i])
		for j := range peers {
			if i != j {
				peers[i].Link(peers[j].ID())
				tables[i].AddPeer(peers[j].ID(),
					wire.NewUnsignedAddress(wire.TCP,
						fmt.Sprintf("%v:%v", "localhost", uint16(3333+j)), uint64(time.Now().UnixNano())))
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6 * time.Second)
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

	gossipCtx, gossipCancel := context.WithTimeout(ctx, 2 * time.Second)
	defer gossipCancel()
	peers[0].Gossip(gossipCtx, contentIDX[:], nil)
	time.Sleep(time.Second)

	for i := range peers {
		contentX , ok := contentResolvers[i].QueryContent(contentIDX[:])
		if !ok {
			println("Failed for ", i)
			panic("initial gossiped content not received")
		}
		if string(contentX) != msgX {
			panic("gossipped content not synced correctly")
		}
	}
	syncCtx, syncCancel := context.WithTimeout(ctx, 2 * time.Second)
	defer syncCancel()
	for i := range peers {
		if i != 0 {
			go func() {
				peerID := peers[i].ID()
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
			}()
		}
	}
	fmt.Println(string("\033[32m"), "Sync test ran successfully")
}

func (st *SyncTest) Perf(num int, outputName string, testOpts test.Options) {

}
