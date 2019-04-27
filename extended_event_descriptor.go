package main

type ExtendedEventDescriptor struct {
	DescriptorTag        uint8  //8bits
	DescriptorLength     uint8  //8bits
	DescriptorNumber     uint8  //4bits
	LastDescriptorNumber uint8  //4bits
	ISO639LanguageCode   uint32 //24bits
	LengthOfItems        uint8  //8bits
	DescriptorItems      []ExtendedEventDescriptorItem
	TextLength           uint8 //8bits
	TextChar             []byte
}

type ExtendedEventDescriptorItem struct {
	ItemDescriptorLength uint8 //8bits
	ItemDescriptorChar   []byte
	ItemLength           uint8 //8bits
	ItemChar             []byte
}

func (d Descriptor) ExtendedEventDescriptor() ExtendedEventDescriptor {
	descriptorTag := d[0]
	descriptorLength := d[1]
	descriptorNumber := (d[2] & 0xf0) >> 4
	lastDescriptorNumber := d[2] & 0x0f
	ISO639LanguageCode := uint32(d[3])<<16 | uint32(d[4])<<8 | uint32(d[5])
	lengthOfItems := d[6]
	descriptorItems := []ExtendedEventDescriptorItem{}
	offset := uint8(7)
	for {
		if offset-6 >= lengthOfItems {
			break
		}
		itemDescriptorLength := d[offset]
		itemDescriptionChar := d[offset+1 : itemDescriptorLength+offset+1]
		itemLength := d[itemDescriptorLength+offset+1]
		itemChar := d[itemDescriptorLength+offset+2 : itemDescriptorLength+offset+itemLength+2]
		offset += itemDescriptorLength + itemLength + 2
		descriptorItem := ExtendedEventDescriptorItem{
			ItemDescriptorLength: itemDescriptorLength,
			ItemDescriptorChar:   itemDescriptionChar,
			ItemLength:           itemLength,
			ItemChar:             itemChar,
		}
		descriptorItems = append(descriptorItems, descriptorItem)
	}
	textLength := d[offset]
	textChar := d[offset+1 : offset+1+textLength]
	return ExtendedEventDescriptor{
		DescriptorTag:        descriptorTag,
		DescriptorLength:     descriptorLength,
		DescriptorNumber:     descriptorNumber,
		LastDescriptorNumber: lastDescriptorNumber,
		ISO639LanguageCode:   ISO639LanguageCode,
		LengthOfItems:        lengthOfItems,
		DescriptorItems:      descriptorItems,
		TextLength:           textLength,
		TextChar:             textChar,
	}
}
