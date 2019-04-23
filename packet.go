package main

type Packet []byte

func (p Packet) SyncByte() uint8 {
	return p[0]
}

func (p Packet) TransportErrorIndicator() uint8 {
	return (p[1] & 0x80) >> 7
}

func (p Packet) PayloadUnitStartIndicator() uint8 {
	return (p[1] & 0x40) >> 6
}

func (p Packet) TransportPriority() uint8 {
	return (p[1] & 0x20) >> 5
}

func (p Packet) PID() uint16 {
	return uint16(p[1]&0x1f)<<8 | uint16(p[2])
}

func (p Packet) TransportScramblingControl() uint8 {
	return (p[3] & 0xc0) >> 6
}

func (p Packet) ContinuityCounter() uint8 {
	return (p[3] & 0x0f)
}

func (p Packet) hasAdaptationField() uint8 {
	return (p[3] & 0x20) >> 5
}

func (p Packet) hasPayload() uint8 {
	return (p[3] & 0x10) >> 4
}

func (p Packet) AdaptarionFieldLength() uint8 {
	if p.hasAdaptationField() != 1 {
		return 0
	}
	return p[4]
}

func (p Packet) Payload() []byte {
	offset := 4
	if p.hasAdaptationField() == 1 {
		offset += 1 + int(p.AdaptarionFieldLength())
	}
	if p.hasPayload() != 1 {
		return []byte{}
	}
	return p[offset:]
}
