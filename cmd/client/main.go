package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	/** if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <multiaddr> <peer id>", os.Args[0])
		return
	}**/
	//if err := dialQuic(os.Args[1], os.Args[2]); err != nil {
	//	log.Fatalf(err.Error())
	//}
	dialQuic("/ip4/0.0.0.0/udp/7654/quic", "QmXwFNChSPjF1KB4RcrimiPf1gADKYBLyP8yKubUsMtz9t")
}

func dialQuic(raddr string, p string) {
	peerID, err := peer.Decode(p)
	check(err)
	addr, err := ma.NewMultiaddr(raddr)
	check(err)
	priv, _, err := ic.GenerateECDSAKeyPair(rand.Reader)
	check(err)

	t, err := libp2pquic.NewTransport(priv, nil, nil, nil)
	check(err)
	log.Printf("Dialing %s\n", addr.String())
	conn, err := t.Dial(context.Background(), addr, peerID)
	check(err)
	defer conn.Close()
	str, err := conn.OpenStream(context.Background())
	check(err)
	defer str.Close()
	msg := bufio.NewReadWriter(bufio.NewReader(str), bufio.NewWriter(str))
	go writeData(msg)
	go readData(msg)
	select {}

}

func writeData(rw *bufio.ReadWriter) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str := scanner.Text()

		if str != "\n" {

			if fileExists(str) {
				fmt.Printf("[*] File %s exists ", str)
				rw.Write([]byte("\x1b[32m%s\x1b[0m> wrote data"))

				file, _ := os.Open("test.jpg")
				b, _ := ioutil.ReadAll(file)
				rw.Write(b)
				rw.Write([]byte("\x1b[32m%s\x1b[0m> wrote data"))
				rw.WriteString("\n")
				//for {rw.Write([]byte("\x1b[32m%s\x1b[0m> wrote data\n"))}
				continue
			}
		}

	}

}

func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			rw.Write([]byte("\x1b[32m%s\x1b[0m> write data\n"))
			// Green console colour:    \x1b[32m
			// Reset console colour:    \x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[32m> ", str)
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return err == nil
}
func check(err error) {
	if err != nil {
		log.Printf("[*] %e", err)
	}
}
