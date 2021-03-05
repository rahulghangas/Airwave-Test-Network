package network

import (
	"context"
	"fmt"
	"github.com/renproject/aw/dht/dhtutil"
	"github.com/renproject/aw/peer"
	"github.com/renproject/aw/wire"
	"github.com/renproject/id"
	"os"
	"sync/atomic"
	"time"
)

func Run() {
	pwd, _ := os.Getwd()
	fo, err := os.Create(pwd + "/../output/output.txt")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	opts, p, _, contentResolver, _, t := setup()

	opts.SyncerOptions = opts.SyncerOptions.WithTimeout(10000 * time.Second)
	opts.GossiperOptions = opts.GossiperOptions.WithTimeout(10000 * time.Second)
	p = peer.New(
		opts,
		t)
	p.Resolve(context.Background(), contentResolver)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var x int64 = 0
	go func() {
		t.Receive(context.Background(), func() func(from id.Signatory, msg wire.Msg) error {
			go func() {
				var seconds int64 = 0
				ticker := time.NewTicker(time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						atomic.AddInt64(&seconds, 1)
						//fmt.Printf("Average throughput for peer %v: %v/second\n", index, x/seconds)
						_, err = fo.WriteString(fmt.Sprintf("%d\n", x/seconds))
						if err != nil {
							fmt.Printf("error writing to file: %v", err)
						}
					}
				}
			}()
			return func(from id.Signatory, msg wire.Msg) error {
				atomic.AddInt64(&x, 1)
				if err := p.Syncer().DidReceiveMessage(from, msg); err != nil {
					return err
				}
				if err := p.Gossiper().DidReceiveMessage(from, msg); err != nil {
					return err
				}
				return nil
			}
		}())
		t.Run(ctx)
	}()

	ctxGossip, cancelGossip := context.WithTimeout(context.Background(), time.Second*3)
	for iter := 0; iter < 1000; iter++ {
		select {
		case <-ctx.Done():
			cancel()
			ctxGossip, cancelGossip = context.WithTimeout(context.Background(), time.Second*3)
		default:
		}

		msgHello := fmt.Sprintf(string(dhtutil.RandomContent()), p.ID().String())
		contentID := id.NewHash([]byte(msgHello))
		contentResolver.InsertContent(contentID[:], []byte(msgHello))
		p.Gossip(ctxGossip, contentID[:], &peer.DefaultSubnet)
	}
	cancelGossip()
	println("Node up and running")
	time.Sleep(5*time.Second)

	var prev int64 = -1
	for {
		if x := atomic.LoadInt64(&x); x == prev {
			break
		}
		prev = x
		time.Sleep(time.Second)
	}
}
