package main

import (
	"encoding/binary"
	"io"
	"os"
)

type TS struct {
	Section map[uint8]Section
}

type Section struct {
	Event map[uint16]Event
}

type Event struct {
	ShortEventDescriptor    ShortEventDescriptor
	ExtendedEventDescriptor ExtendedEventDescriptor
}

func (p Packet) Parser() ([]byte, []byte) {
	PID := p.PID()
	switch PID {
	case 0x12, 0x26, 0x27:
		if p.hasPayload() == 0 {
			break
		}
		payload := p.Payload()
		if p.PayloadUnitStartIndicator() == 0 {
			return payload, nil
		}
		if payload[0] == 0 {
			return nil, payload
		}
		return payload[1:], nil
	default:
		return nil, nil
	}
	return nil, nil
}

func TSFile(f *os.File, ids []byte) TS {
	ts := TS{map[uint8]Section{}}

	pidBuf := map[uint16][]byte{}

	for {
		p := make([]byte, 188)
		err := binary.Read(f, binary.BigEndian, &p)
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
			table.Parser(ids, &ts)
		}
		if next != nil {
			pidBuf[PID] = next
		}
	}
	return ts
}
