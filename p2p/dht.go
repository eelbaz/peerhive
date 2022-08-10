package p2p

import (
	"context"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

func DHT(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr) (*dht.IpfsDHT, error) {

	// Define Bootstrap Nodes.
	peers := []multiaddr.Multiaddr{
		multiaddr.StringCast("/ip4/172.105.135.138/udp/7654/quic/p2p/12D3KooWCg73QYCExwBaXgBp9cs4CadXF44fktQdz7tM9jnzUEX5"),
		multiaddr.StringCast("/ip4/172.105.135.138/tcp/7654/p2p/12D3KooWCg73QYCExwBaXgBp9cs4CadXF44fktQdz7tM9jnzUEX5"),
	}
	bootstrapPeers = append(bootstrapPeers, peers...)

	var options []dht.Option

	// if no bootstrap node are available make this node act as a bootstraping node
	// so other peers may use this node's ipfs address for peer discovery via dht.
	options = append(options, dht.Mode(dht.ModeAuto))

	//
	dht, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}

	if err = dht.Bootstrap(ctx); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Printf("Error connecting to node %q: %-v", peerinfo, err)
			} else {
				log.Printf("Connected to bootstrap node: %q", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return dht, nil
}
