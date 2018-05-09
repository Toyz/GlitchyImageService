package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Toyz/GlitchyImageService/engine"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func readFully(conn net.Conn) ([]byte, error) {
	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result.Bytes(), nil
}

func main() {
	if getenv("MODE", "release") == "debug" {
		go func() {
			time.Sleep(2 * time.Second)
			test_Client()
			os.Exit(0)
		}()
	}

	server := engine.NewServer(getenv("LISTEN", "0.0.0.0:1200"))

	server.Listen(ProcessUser)
}

func ProcessUser(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			checkError(err)
			continue
		}

		go func(conn net.Conn) {
			log.Println("Accepted new Client")
			encoder := gob.NewEncoder(conn)
			decoder := gob.NewDecoder(conn)

			var packet engine.Packet
			decoder.Decode(&packet)

			if packet.ID == 0 {
				packet.From = engine.Glitch{
					Name:        packet.To.Name,
					Mime:        packet.To.Mime,
					Expressions: packet.To.Expressions,
					IsGif:       packet.To.IsGif,
				}

				_, buff, bounds := engine.ProcessImage(bytes.NewReader(packet.To.File), strings.ToLower(packet.To.Mime), packet.To.Expressions)
				packet.From.Bounds = bounds
				packet.From.File = buff.Bytes()

				packet.To = engine.Glitch{}
			}
			encoder.Encode(packet)

			conn.Close() // we're finished
			log.Println("Bye Bye Client")
		}(conn)
	}
}

func test_Client() {
	data, err := ioutil.ReadFile("./kanna.jpg")
	checkError(err)

	conn, err := net.Dial("tcp", getenv("LISTEN", "127.0.0.1:1200"))
	checkError(err)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	encoder.Encode(engine.Packet{
		ID: 0,
		To: engine.Glitch{
			Name:  "kanna.jpg",
			Mime:  "image/jpg",
			IsGif: false,
			Expressions: []string{
				"H ^ L",
			},
			File: data,
		},
	})

	var packet engine.Packet
	decoder.Decode(&packet)
	ioutil.WriteFile("kanna_deb.jpg", packet.From.File, 0777)

	conn.Close()
}
