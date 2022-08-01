package utils

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func SendFile(rw *bufio.ReadWriter, path string) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// Send filename
	filename := filepath.Base(path)
	_, err = rw.Write([]byte(fmt.Sprintf("%s\n", filename)))
	if err != nil {
		log.Println(err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println(err)
	}

	// Send file size
	fileStats, err := file.Stat()
	if err != nil {
		log.Println(err)
	}
	fileSize := fileStats.Size()
	fmt.Println(fileSize)
	_, err = rw.Write([]byte(fmt.Sprintf("%d\n", fileSize)))
	if err != nil {
		log.Println(err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println(err)
	}

	// Send file chunks
	hashFunc := sha1.New()
	const chunkSize = 2048
	buffer := make([]byte, chunkSize)
	for {
		_, err = file.Read(buffer)
		if err != nil {
			// if err is not EOF, error during reading, else stop
			if err != io.EOF {
				log.Println(err)
			} else {
				break
			}
		}

		_, err = rw.Write(buffer)
		if err != nil {
			log.Println(err)
		}
		err = rw.Flush()
		if err != nil {
			log.Println(err)
		}

		hashFunc.Write(buffer)
	}
	fmt.Println("File sent", buffer, len(buffer))

	// Send file hash
	hash := hashFunc.Sum(nil)
	fmt.Println(hash)

	_, err = rw.Write(append(hash, '\n'))
	if err != nil {
		log.Println(err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println(err)
	}
}

func RecvFile(rw *bufio.ReadWriter) {
	// Recv filename
	filename, err := rw.ReadString('\n')
	filename = strings.Replace(filename, "\n", "", -1)
	if err != nil {
		log.Println(err)
	}
	file, err := os.Create(filename + ".tmp")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// Recv file size
	str, err := rw.ReadString('\n')
	str = strings.Replace(str, "\n", "", -1)
	if err != nil {
		log.Println(err)
	}
	fileSize, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
	}

	// Recv file chunks
	hashFunc := sha1.New()
	buf := make([]byte, 2048)
	for bytesRead := 0; bytesRead < fileSize; {
		// Receive bytes
		n, err := rw.Read(buf)
		if err != nil {
			log.Println(err)
		}

		bytesRead += n

		// Write bytes to the file
		file.Write(buf)

		// Add to hash
		hashFunc.Write(buf)
	}
	fmt.Println("File received")

	// Recv file hash
	hashRecv, err := rw.ReadBytes('\n')
	if err != nil {
		log.Println(err)
	}
	hashRecv = bytes.Trim(hashRecv, "\x00") // remove leading or trailing 0x00
	hashRecv = hashRecv[:len(hashRecv)-1]   // remove trailing \n
	fmt.Println("Hash received", hashRecv)

	hash := hashFunc.Sum(nil)
	if string(hash) != string(hashRecv) {
		log.Println("File NOT received correctly: hashes not equal")
		log.Println("Computed hash:", hash)
	} else {
		log.Println("File received correctly: hashes equal")
	}
}
