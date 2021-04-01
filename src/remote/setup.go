package remote

import (
	"awt/test"
	"context"

	"github.com/renproject/aw/channel"
	"github.com/renproject/aw/dht"
	"github.com/renproject/aw/handshake"
	"github.com/renproject/aw/peer"
	"github.com/renproject/aw/transport"
	"github.com/renproject/id"
	"go.uber.org/zap"
	"time"
)

func duration(num int) time.Duration {
	return time.Duration(num) * time.Second
}

func setup(key *id.PrivKey, testOpts test.Options) (peer.Options, *peer.Peer, dht.Table, dht.ContentResolver, *channel.Client, *transport.Transport) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	// Init options for all peers.
	opts := peer.DefaultOptions().WithPrivKey(key).WithLogger(logger).
		WithGossiperOptions(peer.DefaultGossiperOptions().WithLogger(logger)).
		WithSyncerOptions(peer.DefaultSyncerOptions().WithLogger(logger))

	opts.GossiperOptions.Timeout = duration(2)
	opts.SyncerOptions.WiggleTimeout = duration(2)

	self := opts.PrivKey.Signatory()
	h := handshake.Filter(func(id.Signatory) error { return nil }, handshake.ECIES(opts.PrivKey))
	client := channel.NewClient(
		channel.DefaultOptions().
			WithLogger(logger),
		self)
	table := dht.NewInMemTable(self)
	contentResolver := dht.NewDoubleCacheContentResolver(dht.DefaultDoubleCacheContentResolverOptions(), nil)
	t := transport.New(
		transport.DefaultOptions().
			WithLogger(logger).
			WithClientTimeout(5*time.Second).
			WithOncePoolOptions(handshake.DefaultOncePoolOptions().WithMinimumExpiryAge(10*time.Second)).
			WithPort(uint16(8080)),
		self,
		client,
		h,
		table)
	p := peer.New(
		opts,
		t)
	p.Resolve(context.Background(), contentResolver)

	return opts, p, table, contentResolver, client, t
}
