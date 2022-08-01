package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	var buf []byte = make([]byte, 2048)

	addr, _ := net.ResolveUDPAddr("udp", ":8100")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()
	os.Create("/Users/ex0bit/Downloads/stream.ts")
	f, _ := os.OpenFile("/Users/ex0bit/Downloads/stream.ts", syscall.O_WRONLY, 0644)
	defer f.Close()

	go ffmpegCmd()
	time.Sleep(5 * time.Second)
	go ffplayCmd()

	for {
		n, _ := conn.Read(buf)
		if n > 1480 {
			fmt.Printf("buf size: %v\n", n)
		}
		f.Write(buf[0:n])
		buf = make([]byte, 2048)
	}

}

/**
* @author ex0bit
* @date 2020/8/5
* @describe - FFMPEG - Peerhive - Hivestream command line - H264/TS
--> ffmpeg -threads 4  -f avfoundation -framerate 30 -video_size 960x540 -i "1:1" -c:v h264 -c:a aac -preset ultrafast -tune zerolatency -f mpegts udp://localhost:12345
**/

func ffmpegCmd() {
	cmd := exec.Command(
		"ffmpeg",        // path to ffmpeg executable
		"-threads", "4", // number of threads
		"-f", "avfoundation", // input format
		"-framerate", "30", // frame rate
		"-video_size", "960x540", // video size
		"-i", "1:0", // input device
		"-c:v", "h264", // video codec
		"-c:a", "aac", // audio codec
		"-preset", "ultrafast", // preset
		"-tune", "film", // tune
		"-f", "mpegts", // output format
		"udp://localhost:8100", // output URL
	)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error in ffmpegCmd: ", err)
	}

	defer cmd.Process.Kill()
	defer cmd.Process.Release()

}

func ffplayCmd() {
	cmd := exec.Command("/usr/local/bin/mpv", "/Users/ex0bit/Downloads/stream.ts")

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error in ffplayCmd: ", err)
	}

	defer cmd.Process.Kill()
	defer cmd.Process.Release()

}
