package main

import (
	"bufio"
	"fmt"
	"github.com/r00ky/goserial"
	"os"
	"strings"
)

func main() {

	ttyPath := os.Args[1]
	fmt.Println("Writing packets to: " + ttyPath)
	fmt.Println("Packets are formatted all lower case hex, [src] [dest] [message...] (skip checksum)")
	fmt.Println("Quote ascii strings for text conversion")

	port, err := goserial.OpenPort(&goserial.Config{Name: ttyPath, Parity: goserial.ParityEven, Baud: 9600})
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		parser := NewIbusPacketParser()
		for {
			nextByte := make([]byte, 1)
			port.Read(nextByte)
			parser.Push(nextByte)
			if parser.HasPacket() {
				pkt := parser.GetPacket()
				fmt.Println("\n<== " + pkt.AsString())
				fmt.Print("Enter IBUS packet: ")
			}
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter IBUS packet: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		hexChars := strings.Split(text, " ")
		packet := new(IbusPacket)
		packet.Src = hexChars[0]
		packet.Dest = hexChars[1]
		packet.Message = hexChars[2:len(hexChars)]
		fmt.Println("==> " + packet.AsString())
		bytes := packet.AsBytes()
		port.Write(bytes)
	}
}
