package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
)

func NewNode(ctx context.Context, seed int64, port int) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if seed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(seed))
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.ECDSA, 2048, r)
	if err != nil {
		return nil, err
	}

	ip6quic := fmt.Sprintf("/ip6/::/udp/%d/quic", port)
	ip4quic := fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", port)

	ip6tcp := fmt.Sprintf("/ip6/::/tcp/%d", port)
	ip4tcp := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)

	//add bootstrap
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(ip6quic, ip4quic, ip6tcp, ip4tcp),
		libp2p.Identity(priv),
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.DefaultMuxers,
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.FallbackDefaults,
		libp2p.EnableAutoRelay(),
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
	}

	return libp2p.New(opts...)
}
