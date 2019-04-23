package main

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
