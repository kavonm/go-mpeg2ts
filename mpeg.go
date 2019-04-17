package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Packet struct {
	SyncByte                   uint8
	TransportErrorIndicator    uint8  //1bit
	PayloadUnitStartIndicator  uint8  //1bit
	TransportPriority          uint8  //1bit
	PID                        uint16 // 13bit
	TransportScramblingControl uint8  // 2bit
	AdaptationFieldControl     uint8  // 2bit
	ContinuityCounter          uint8  // 4bit

	//AdaptationFieldControl is 10 or 11
	AdaptationField []byte

	//AdaptationFieldControl is 01 or 11
	DataByte []byte
}

const PaketSize = 188

func parse(raw []byte) Packet {
	SyncByte := uint8(raw[0])
	//fmt.Println(SyncByte)
	TransportErrorIndicator := uint8(raw[1]&0x80) >> 7
	PayloadUnitStartIndicator := uint8(raw[1]&0x40) >> 6
	TransportPriority := uint8(raw[1]&0x20) >> 5
	PID := uint16(raw[1]&0x1F)<<8 | uint16(raw[2])
	TransportScramblingControl := uint8(raw[3]&0xc0) >> 6
	//AdaptationFieldControl := uint8(raw[3]&0x0a) >> 2
	ContinuityCounter := uint8(raw[3] & 0x0f)
	hasAdaptationField := uint8(raw[3]&0x20) >> 5
	hasPayload := uint8(raw[3]&0x10) >> 4

	var AdaptationField []byte
	if hasAdaptationField == 1 {
		start := 4
		AdaptationFieldLength := int(raw[start])
		end := start + AdaptationFieldLength + 1
		AdaptationField = raw[start:end]
	}

	var DataByte []byte
	if hasPayload == 1 {
		start := 4
		if hasAdaptationField == 1 {
			AdaptationFieldLength := int(raw[start])
			start += 1 + AdaptationFieldLength
		}
		if PayloadUnitStartIndicator == 0 {
			//fllowing packet
			DataByte = raw[start:]
		} else {
			//initial packet
			prefix := raw[start : start+3]
			if bytes.Equal(prefix, []byte{0x00, 0x00, 0x01}) {
				//PES
				DataByte = raw[start:]
			} else {
				//PSI
				pointer := int(raw[start])
				start += 1 + pointer
				DataByte = raw[start:]
			}
		}
	}

	return Packet{
		SyncByte:                   SyncByte,
		TransportErrorIndicator:    TransportErrorIndicator,
		PayloadUnitStartIndicator:  PayloadUnitStartIndicator,
		TransportPriority:          TransportPriority,
		PID:                        PID,
		TransportScramblingControl: TransportScramblingControl,
		AdaptationFieldControl:     hasAdaptationField,
		ContinuityCounter:          ContinuityCounter,
		AdaptationField:            AdaptationField,
		DataByte:                   DataByte,
	}
}

type ProgramAssociationSection struct {
	TableID                uint8  //uint8
	SectionSyntaxIndicator uint8  //1bit
	Padding                uint8  //1bit
	Reserved1              uint8  //2bits
	SectionLength          uint16 //12bits
	TransportStreamID      uint16 //16bits
	Reserved2              uint8  //2bits
	VersionNumber          uint8  //5bits
	CurrentNextIndicator   uint8  //1bit
	SectionNumber          uint8  //8bits
	LastSectionNumber      uint8  //8bits
	DataByte               []ProgramPID
	CRC                    uint32 //32bit
}

type ProgramPID struct {
	ProgramNumber uint16
	Reserved      uint8  //3bits
	NetworkPID    uint16 //13bits
	ProgramMapPID uint16 //13bits
}

func ProgramAssociationSectionParse(raw []byte) ProgramAssociationSection {
	TableID := uint8(raw[0])
	SectionSyntaxIndicator := uint8(raw[1]&0x80) >> 7
	Padding := uint8(raw[1]&0x40) >> 6
	Reserved1 := uint8(raw[1]&0x30) >> 4
	SectionLength := uint16(raw[1]&0x0f)<<8 | uint16(raw[2])
	TransportStreamID := uint16(raw[3])<<8 | uint16(raw[4])
	Reserved2 := uint8(raw[5]&0xa0) >> 6
	VersionNumber := uint8(raw[5]&0x3e) >> 1
	CurrentNextIndicator := uint8(raw[5] & 0x01)
	SectionNumber := uint8(raw[6])
	LastSectionNumber := uint8(raw[7])

	length := int(SectionLength) - 9
	data := raw[8 : 8+length+1]
	DataByte := []ProgramPID{}
	for i := 0; i < length; i += 4 {
		ProgramNumber := uint16(data[i])<<8 | uint16(data[i+1])
		Reserved := uint8(data[i+2]&0xe0) >> 5
		var NetworkPID uint16
		var ProgramMapPID uint16
		if ProgramNumber == 0 {
			NetworkPID = uint16(data[i+2]&0x1f)<<8 | uint16(data[i+3])
		} else {
			ProgramMapPID = uint16(data[i+2]&0x1f)<<8 | uint16(data[i+3])
		}
		DataByte = append(DataByte, ProgramPID{
			ProgramNumber: ProgramNumber,
			Reserved:      Reserved,
			NetworkPID:    NetworkPID,
			ProgramMapPID: ProgramMapPID,
		})
	}

	return ProgramAssociationSection{
		TableID:                TableID,
		SectionSyntaxIndicator: SectionSyntaxIndicator,
		Padding:                Padding,
		Reserved1:              Reserved1,
		SectionLength:          SectionLength,
		TransportStreamID:      TransportStreamID,
		Reserved2:              Reserved2,
		VersionNumber:          VersionNumber,
		CurrentNextIndicator:   CurrentNextIndicator,
		SectionNumber:          SectionNumber,
		LastSectionNumber:      LastSectionNumber,
		DataByte:               DataByte,
	}
}

type ProgramMapTable struct {
	TableID                uint8
	SectionSyntaxIndicator uint8  //1bit
	Padding                uint8  //1bit
	Reserved1              uint8  //2bits
	SectionLength          uint16 //12bits
	ProgramNumber          uint16 //16bits
	Reserved2              uint8  //2bits
	VersionNumber          uint8  //5bits
	CurrentNextIndicator   uint8  //1bit
	SectionNumber          uint8  //8bits
	LastSectionNumber      uint8  //8bits
	Reserved3              uint8  //3bits
	PCRPID                 uint16 //13bits
	Reserved4              uint8  //4bits
	ProgramInfoLength      uint16 //12bits

	Descriptor []byte

	DataByte []byte

	CRC uint32 //32bits
}

func TransportStreamFile(r *os.File) {
	counter := map[uint16]int{}
	PAT := map[uint16]uint16{}
	PID := map[uint16]map[uint8]int{}
	//temp := []byte{}
	for {
		buffer := make([]byte, 188)
		err := binary.Read(r, binary.BigEndian, &buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		packet := parse(buffer)
		if _, ok := counter[packet.PID]; !ok {
			counter[packet.PID] = 1
		} else {
			counter[packet.PID]++
		}

		// PAT
		if packet.PID == 0x00 {
			pat := ProgramAssociationSectionParse(packet.DataByte)
			for _, v := range pat.DataByte {
				if v.ProgramNumber == 0 {
					PAT[v.ProgramNumber] = v.NetworkPID
				} else {
					PAT[v.ProgramNumber] = v.ProgramMapPID
				}
			}
		}
		for _, v := range PAT {
			if packet.PID == v {
				if _, ok := PID[packet.PID]; !ok {
					PID[packet.PID] = map[uint8]int{}

					PID[packet.PID][packet.PayloadUnitStartIndicator] = 1
				} else {
					PID[packet.PID][packet.PayloadUnitStartIndicator]++
				}
			}
		}
	}
	fmt.Println(counter)
	fmt.Println(PAT)
	fmt.Println(PID)
}

func main() {
	server()
	/*
		f, err := os.Open("test.ts")
		if err != nil {
			fmt.Println(err)
		}
		TransportStreamFile(f)
	*/
}
