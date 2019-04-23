package main

type Descriptor []byte

func (d Descriptor) DescriptorTag() uint8 {
	return d[0]
}

func (d Descriptor) DescriptorLength() uint8 {
	return d[1]
}
