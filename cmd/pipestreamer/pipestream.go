package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func main() {
	println("pipestreamer")
	go ExecuteFFMpegCmd("12345")
	conn := FFMpegUDPReader("12345")
	HandleUDPConnection(conn)

}

//Given a port number return a UDP connection to the port to receive data from
func FFMpegUDPReader(port string) *net.UDPConn {
	addr, _ := net.ResolveUDPAddr("udp", ":"+port)
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()
	return conn
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
		//"|", "ffplay",           // ffplay command
	)
	fmt.Println(cmd.Args)
	//cmd.Stdout = os.Stdout   // redirect stdout to terminal
	cmd.Stderr = os.Stderr   // redirect stderr to terminal
	cmd.Run()                // run the command
	defer cmd.Process.Kill() // kill the command
}

func HandleUDPConnection(u io.Reader) {
	// create the pipe and tee reader
	pr, pw := io.Pipe()
	fmt.Println(pr.Read(make([]byte, 1024)))
	tr := io.TeeReader(u, pw)

	io.TeeReader(tr, os.Stdout) // redirect to terminal

	// Everything read from r will be copied to stdout.
	if _, err := io.ReadAll(pr); err != nil {
		log.Fatal(err)
	}
}
