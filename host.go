package main

import (
	"context"
	"crypto/rand"
	"io"
	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"
)

func NewHost(ctx context.Context, seed int64, port int) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if seed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(seed))
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 4096, r)
	if err != nil {
		return nil, err
	}

	//add bootstrap
	multiAddr, err := multiaddr.NewMultiaddr("/ip4/172.105.135.138/udp/46351/quic/p2p/12D3KooWD7PXURAfkktsGKAVtwNnvgJKUzVksXCqJcsX7pKHxQKX")
	if err != nil {
		panic(err)
	}

	opts := []libp2p.Option{
		libp2p.EnableHolePunching(),
		libp2p.Identity(priv),
		libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
		libp2p.EnableRelayService(),
		libp2p.ListenAddrs(multiAddr),
	}

	return libp2p.New(opts...)
}
