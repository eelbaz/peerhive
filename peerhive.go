package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"

	p2p "peerhive/p2p"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
)

type Config struct {
	Port           int
	ProtocolID     string
	Rendezvous     string
	Seed           int64
	DiscoveryPeers addrList
	BootstrapRelay bool
}

func main() {

	// 1: commandline config
	config := Config{}

	flag.StringVar(&config.Rendezvous, "rendezvous", "/peerhive", "")
	flag.Int64Var(&config.Seed, "seed", 0, "Seed value for generating a PeerID, 0 is random")
	flag.Var(&config.DiscoveryPeers, "peer", "Peer multiaddress for peer discovery")
	flag.StringVar(&config.ProtocolID, "protocolid", "/pubsub", "")
	flag.IntVar(&config.Port, "port", 0, "")
	flag.BoolVar(&config.BootstrapRelay, "bootstraprelay", false, "bootstraprelay flag sets node to background relay and bootstrap mode only")
	flag.Parse()

	//Create context with Cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//2: Create new libp2p Host with host options from command line
	h, err := p2p.NewNode(ctx, config.Seed, config.Port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Node assigned addresses:\n")
	for i, addr := range h.Addrs() {
		fmt.Printf("%v: %s/p2p/%s\n", i+1, addr.String(), h.ID())
	}
	fmt.Printf("\n")

	//3: create DHT
	dht, err := p2p.DHT(ctx, h, config.DiscoveryPeers)
	if err != nil {
		panic(err)
	}

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	//Public DHT Discovery
	go p2p.Discover(ctx, h, dht, config.Rendezvous)

	// setup local mDNS discovery
	if err := setupMDNSDiscovery(h, config.Rendezvous); err != nil {
		panic(err)
	}

	// join the pubsub topic
	topic, err := ps.Join(config.Rendezvous)
	if err != nil {
		panic(err)
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	// start the publisher and subscriber
	go subscribe(sub, ctx, h.ID())
	if !config.BootstrapRelay {
		go publish(ctx, topic, config.BootstrapRelay)
	}

	go ExecuteFFMpegCmd("12345")                          //execute ffmpeg command
	r := HandleUDPConnection(FFMpegUDPConnected("12345")) //handle UDP connection
	fmt.Println("r is: ", r)
	// wait for ctrl-c
	select {} //hang forever to allow publish to run and program to background
}

// start subsriber to topic
func subscribe(subscriber *pubsub.Subscription, ctx context.Context, hostID peer.ID) {
	for {
		msg, err := subscriber.Next(ctx)
		if err != nil {
			panic(err)
		}

		// only consider messages delivered by other peers
		if msg.ReceivedFrom == hostID {
			continue
		}
		fmt.Printf("%s -> %s\n", msg.ReceivedFrom.Pretty(), string(msg.Data))
	}
}

// start publisher to topic
func publish(ctx context.Context, topic *pubsub.Topic, bootsrapRelay bool) {
	//we dont want a bootstrap publish to send messages, just relay them
	if bootsrapRelay {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		for scanner.Scan() {

			fmt.Printf("(msg):\n")
			msg := scanner.Text()
			if len(msg) != 0 {
				// publish message to topic
				bytes := []byte(msg)
				topic.Publish(ctx, bytes)
			}
		}
	}
}

type addrList []multiaddr.Multiaddr

func (al *addrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
	addr, err := multiaddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

/** shortID returns the last 8 chars of a base58-encoded peer id.
func shortID(p peer.ID) string {
	pretty := p.Pretty()
	return pretty[len(pretty)-8:]
}**/

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if n.h.ID() != pi.ID {
		fmt.Printf("(new peer discovered): %s\n", pi.ID.Pretty())
		err := n.h.Connect(context.Background(), pi)
		if err != nil {
			fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
		}
	}

}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupMDNSDiscovery(h host.Host, ns string) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, ns, &discoveryNotifee{h: h})
	return s.Start()
}

//given a port number, execute ffmpeg command line process to stream video from port
func ExecuteFFMpegCmd(port string) {
	cmd := exec.Command(
		"ffmpeg",             // path to ffmpeg executable
		"-f", "avfoundation", // input format
		"-pixel_format", "yuv420p", // pixel format
		"-framerate", "29.97", // frame rate
		"-video_size", "960x540", // video size
		"-i", "1:0", // input device
		"-c:v", "h264", // video codec
		"-c:a", "aac", // audio filter
		"-preset", "ultrafast", // preset
		"-crf", "29.97", // constant rate factor
		"-f", "mpegts", // output format
		"udp://localhost:"+port, // output URL
	)
	fmt.Println(cmd.Args)
	//cmd.Stdout = os.Stdout   // redirect stdout to terminal
	//cmd.Stderr = os.Stderr   // redirect stderr to terminal
	cmd.Run()                // run the command
	defer cmd.Process.Kill() // kill the command
}

func HandleUDPConnection(udpReader io.Reader) io.Reader {
	buf := make([]byte, 1500)
	go func() {
		for {
			n, err := udpReader.Read(buf)
			if err != nil {
				fmt.Println(err)
			}
			buf = buf[:n]
			fmt.Println(buf)
		}
	}()
	return udpReader
}

//Given a port number return a UDP connection to the port to receive data from
func FFMpegUDPConnected(port string) *net.UDPConn {
	addr, _ := net.ResolveUDPAddr("udp", ":"+port)
	conn, _ := net.ListenUDP("udp", addr)
	return conn
}
