package main;

import (
	"github.com/alecthomas/kingpin"
	"fmt"
	"github.com/paocalvi/serial"
	"bufio"
	"os"
"strings"
	"ibus"
)

var (
	device = kingpin.Arg("device", "serial device to communicate with").Required().String()
)

func main() {
	kingpin.Version("0.1")
	kingpin.Parse()

	ttyPath := *device;

	fmt.Println("Writing packets to: " + ttyPath)
	fmt.Println("Packets are formatted all lower case hex, [src] [dest] [message...]")
	fmt.Println("Quote ascii strings for text conversion")

	serialConfig := &serial.Config{Name: ttyPath, Baud: 9600, Parity:serial.PARITY_EVEN}
	port, err := serial.OpenPort(serialConfig)
	if (err != nil) {
		panic(err)
	}

	go func() {
		parser := ibus.NewIbusPacketParser()
		for {
			buffer := make([]byte, 1)
			_, err := port.Read(buffer)
			if (err != nil) {
				panic(err)
			}
			parser.Push(buffer)
			if (parser.HasPacket()) {
				pkt := parser.GetPacket();
				// echo received packets to simulate bus broadcast
				port.Write(pkt.AsBytes());
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
		packet := new(ibus.IbusPacket)
		packet.Src = hexChars[0]
		packet.Dest = hexChars[1]
		packet.Message = hexChars[2:len(hexChars)]
		fmt.Println("==> " + packet.AsString())
		bytes := packet.AsBytes()
		port.Write(bytes)
	}
}
