package ibus

import (
	"fmt"
	"strconv"
)

type IbusPacketParser struct {
	buffer    []byte
	packet    *IbusPacket
	hasPacket bool
}

func NewIbusPacketParser() (*IbusPacketParser) {
	p := new(IbusPacketParser)
	p.hasPacket = false
	p.buffer = make([]byte, 0)
	return p
}

func (p *IbusPacketParser) Push(bytes []byte) {
	for _, el := range bytes {
		p.buffer = append(p.buffer, el)
	}
	p.parse()
}

func (p *IbusPacketParser) HasPacket() bool {
	return p.hasPacket
}

func (p *IbusPacketParser) GetPacket() (*IbusPacket) {
	p.hasPacket = false
	return p.packet
}

func (p *IbusPacketParser) parse() {
	//p.debug()
	p.hasPacket = false

	if (len(p.buffer)) < 5 {
		//fmt.Println("packet buffer < 5 bytes")
		return // all packets will be at least 5 bytes long
	}

	tmpLength, _ := strconv.ParseInt(getHexStringFromByte(p.buffer[1]), 16, 0)
	length := int(tmpLength)

	if (length < 3) {
		//fmt.Println("packet length field not long enough (" + strconv.Itoa(length) + ")")
		p.shiftBuffer()
		return // not possible to have a length value less than 3 (dest + data + checksum)
	} else if (length > 64) {
		//fmt.Println("packet length field too long (" + strconv.Itoa(length) + ")")
		p.shiftBuffer()
		return
	}

	if (len(p.buffer) < (2 + length)) {
		//fmt.Println("Assuming length of " + strconv.Itoa(length) + ", we need more bytes")
		return // current byte slice is not long enough assuming length value at position 2
	} else {
		packet := new(IbusPacket)
		packet.Src = getHexStringFromByte(p.buffer[0])
		packet.Dest = getHexStringFromByte(p.buffer[2])
		packet.Message = getHexStringSliceFromByteSlice(p.buffer[3:(3 + length - 2)])
		packet.Checksum = getHexStringFromByte(p.buffer[2 + length - 1])
		if (packet.IsValid()) {
			p.hasPacket = true
			p.packet = packet
			p.buffer = p.buffer[2 + length : len(p.buffer)]
		} else {
			// we have enough bytes to theoretically be a valid packet, but either the checksum failed for unknown
			// reasons (unlikely), or buffer[0] is not actually the beggining of a packet. In this case, we will shift
			// off the first byte, as we know it is useless. This then means we are too short, but the new buffer[0]
			// might be a correct length byte, such that the next byte pushed onto the buffer will create a valid packet
			//fmt.Println("packet validation error.")
			p.shiftBuffer()
		}
	}
}

func (p *IbusPacketParser) shiftBuffer() {
	//fmt.Println("Shifting buffer")
	p.buffer = p.buffer[1:len(p.buffer)]
	//p.debug()
}

func (p *IbusPacketParser) debug() {
	fmt.Println("")
	fmt.Println("")
	for _, el := range p.buffer {
		fmt.Print(getHexStringFromByte(el) + " ")
	}
	fmt.Println("")
}
