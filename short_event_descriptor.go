package main

type ShortEventDescriptor struct {
	DescriptorTag      uint8  //8bits
	DescriptorLength   uint8  //8bits
	ISO639LanguageCode uint32 //24bits
	EventNameLength    uint8  //8bits
	EventNameChar      []byte
	TextLength         uint8 //8bits
	TextChar           []byte
}

func (d Descriptor) ShortEventDescriptor() ShortEventDescriptor {
	descriptorTag := d[0]
	descriptorLength := d[1]
	ISO639LanguageCode := uint32(d[2])<<16 | uint32(d[3])<<8 | uint32(d[4])
	eventNameLength := d[5]
	eventNameChar := d[6 : eventNameLength+6]
	textLength := d[eventNameLength+6]
	textChar := d[eventNameLength+7 : textLength+eventNameLength+7]
	return ShortEventDescriptor{
		DescriptorTag:      descriptorTag,
		DescriptorLength:   descriptorLength,
		ISO639LanguageCode: ISO639LanguageCode,
		EventNameLength:    eventNameLength,
		EventNameChar:      eventNameChar,
		TextLength:         textLength,
		TextChar:           textChar,
	}
}
