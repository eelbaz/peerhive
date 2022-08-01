package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	tpt "github.com/libp2p/go-libp2p-core/transport"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	run("7654")
}

func run(port string) error {
	addr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/udp/%s/quic", port))
	check(err)

	priv, _, err := ic.GenerateECDSAKeyPair(rand.New(rand.NewSource(int64(7654))))
	check(err)

	peerID, err := peer.IDFromPrivateKey(priv)
	check(err)

	t, err := libp2pquic.NewTransport(priv, nil, nil, nil)
	check(err)

	ln, err := t.Listen(addr)
	check(err)

	fmt.Printf("Listening. Now run: go run main.go %s %s\n", ln.Multiaddr(), peerID)
	for {
		conn, err := ln.Accept()
		check(err)
		log.Printf("Accepted new connection from %s (%s)\n", conn.RemotePeer(), conn.RemoteMultiaddr())
		go func() {
			if err := handleConn(conn); err != nil {
				log.Printf("handling conn failed: %s", err.Error())
				if conn.IsClosed() {
					fmt.Println(err)
				}
			}
		}()
	}
}

func handleConn(conn tpt.CapableConn) error {
	str, err := conn.AcceptStream()
	if err != nil {
		return err
	}
	defer str.Close()
	msg := bufio.NewReadWriter(bufio.NewReader(str), bufio.NewWriter(str))
	defer msg.Flush()
	p := make([]byte, 2048)

	for {
		n, err := msg.Read(p)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", p[:n])
	}

	return err
}

func check(err error) {
	if err != nil {
		log.Printf("[*] %v", err)
	}
}
