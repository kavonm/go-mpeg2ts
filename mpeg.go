package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("test.ts")
	if err != nil {
		panic(err)
	}
	data := map[uint8]map[uint16]map[uint8][]ShortEventDescriptor{}
	pidBuf := map[uint16][]byte{}
	//tableBuf := map[uint8][]byte{}
	ids := []byte{0x4E, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f}

	for {
		p := make([]byte, 188)
		err = binary.Read(f, binary.BigEndian, &p)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		packet := Packet(p)
		prev, next := packet.Parser()
		if prev == nil && next == nil {
			continue
		}

		PID := packet.PID()
		if prev != nil {
			pidBuf[PID] = append(pidBuf[PID], prev...)
			continue
		}
		if buf, ok := pidBuf[PID]; ok {
			table := Table(buf)
			temp := table.Parser(ids, data)
			if temp != nil {
				data = temp
			}
		}
		if next != nil {
			pidBuf[PID] = next
		}
	}
	fmt.Println(data)
	fmt.Println("=============================")
	for _, t := range data {
		fmt.Println(t)
	}
}
