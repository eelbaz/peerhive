package filestreamer

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func StreamAndPlayFromPort(fileName string, port string) {
	var buf []byte = make([]byte, 2048) // buffer for incoming data
	//buf := new(bytes.Buffer)

	addr, _ := net.ResolveUDPAddr("udp", ":"+port)
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()
	hdir, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("Error: ", err)
	}
	x := hdir + "/Downloads" + fileName
	fmt.Println("file: ", x)
	os.Remove(x) // remove the file if it exists
	os.Create(x) // create the file

	f, _ := os.OpenFile(x, syscall.O_WRONLY, 0644) // open the file
	defer f.Close()                                // close the file
	//run ffmpeg from the command line and stream the video from the port while sending it to the file
	go FfmpegCmd(port)          // ffmpeg command
	time.Sleep(5 * time.Second) // wait for ffmpeg to start
	go FfplayCmd(x)             // ffplay command
	//streamBuf := make(chan []byte, 2048)
	for {
		n, err := conn.Read(buf) // read the buffer

		if err != nil {
			fmt.Println("Error in Read: ", err)
		}

		buf := buf[:n] // trim to actual size
		//fmt.Printf("aft-len=%d aft-cap=%d \n", len(buf), cap(buf)) // print size
		f.Write(buf) // write to file
	}

}

/**
* @author ex0bit
* @date 2020/8/5
* @describe - FFMPEG - Peerhive - Hivestream command line - H264/TS
--> ffmpeg -threads 4  -f avfoundation -framerate 30 -video_size 960x540 -i "1:1" -c:v h264 -c:a aac -preset ultrafast -tune zerolatency -f mpegts udp://localhost:12345
**/

func FfmpegCmd(port string) {
	cmd := exec.Command(
		"ffmpeg",             // path to ffmpeg executable
		"-f", "avfoundation", // input format
		"-pixel_format", "yuv420p", // pixel format
		"-framerate", "29.97", // frame rate
		"-video_size", "960x540", // video size
		"-i", "1:0", // input device
		"-c:v", "h264", // video codec
		//"-vf", "fps=29.97", // video filter
		//"-c:a", "he-aac", // audio codec
		"-c:a", "aac", // audio filter
		//"-b:a", "copy", // audio bitrate
		"-preset", "ultrafast", // preset
		"-crf", "29.97", // constant rate factor
		//"-tune", "film", // tune
		"-f", "mpegts", // output format
		"-y",                    // overwrite output file
		"udp://localhost:"+port, // output URL
	)
	fmt.Println("ffmpegCmd: ", strings.Join(cmd.Args, " ")) // print the ffmpeg command line for debugging
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error in ffmpegCmd: ", err)
	}

	defer cmd.Process.Kill()
	defer cmd.Process.Release()

}

func FfplayCmd(fileName string) {
	cmd := exec.Command("mpv", fileName)

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error in ffplayCmd: ", err)
	}

	defer cmd.Process.Kill()
	defer cmd.Process.Release()

}
