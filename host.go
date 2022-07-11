package main

import (
	"context"
	"crypto/rand"
	"io"
	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
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

	opts := []libp2p.Option{
		libp2p.EnableHolePunching(),
		libp2p.Identity(priv),
		libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
		libp2p.EnableRelayService(),
	}

	return libp2p.New(opts...)
}
