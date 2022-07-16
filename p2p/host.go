package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
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

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	//add bootstrap

	opts := []libp2p.Option{
		libp2p.EnableHolePunching(),
		libp2p.Identity(priv),
		libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
		libp2p.EnableRelayService(),
	}
	opts = append(opts, libp2p.ListenAddrStrings("/ip4/0.0.0.0/udp/"+fmt.Sprint(port)+"/quic"), libp2p.ListenAddrStrings("/ip6/::/udp/"+fmt.Sprint(port)+"/quic"), libp2p.ListenAddrStrings("/ip6/::/tcp/"+fmt.Sprint(port)+""), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/"+fmt.Sprint(port)+""))

	return libp2p.New(opts...)
}
