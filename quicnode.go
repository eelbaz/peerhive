package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	tp "github.com/libp2p/go-libp2p-core/transport"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/multiformats/go-multiaddr"
)

func quicknode(ctx context.Context, raddr multiaddr.Multiaddr, p peer.ID) {
	fmt.Println("in quicknode")
	tp := getTransport()
	conn, err := tp.Dial(ctx, raddr, p)
	if err != nil {
		log.Println(err)
	}
	ms, _ := conn.OpenStream(ctx)
	defer ms.Close()
	as, _ := conn.AcceptStream()
	defer as.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(as), bufio.NewWriter(ms))
	fmt.Fprintf(rw, "Hello World")

	for {
		rw.Write([]byte("hello" + "Hello:" + p.String() + ":" + raddr.String() + "\n"))
		rw.Flush()
	}
}

func getTransport() tp.Transport {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println("error in getTransport generateRSAKey: ", err)
	}

	key, err := crypto.UnmarshalRsaPrivateKey(x509.MarshalPKCS1PrivateKey(rsaKey))
	if err != nil {
		log.Println(err)
	}
	tr, err := quic.NewTransport(key, nil, nil, nil)
	if err != nil {
		log.Println(err)
	}
	return tr
}
