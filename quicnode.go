package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	tpt "github.com/libp2p/go-libp2p-core/transport"
	quic "github.com/libp2p/go-libp2p-quic-transport"
)

func quicknode() {

}

func getTransport() tpt.Transport {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
		return (nil)
	}

	key, err := ic.UnmarshalRsaPrivateKey(x509.MarshalPKCS1PrivateKey(rsaKey))
	tr, err := quic.NewTransport(key, nil, nil, nil)
	return tr
}
