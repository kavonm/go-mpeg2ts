package main

//the first byte is pointer_field
type Table []byte

func (t Table) TableID() uint8 {
	return t[1]
}

func (t Table) Parser(ids []byte, data map[uint8]map[uint16]map[uint8][]ShortEventDescriptor) map[uint8]map[uint16]map[uint8][]ShortEventDescriptor {
	tableID := t.TableID()
	if b := contain(ids, tableID); b == -1 {
		return nil
	}
	switch tableID {
	case 0x4E, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56,
		0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f:
		eit := t.EventInformationTable()
		data[tableID] = map[uint16]map[uint8][]ShortEventDescriptor{}
		for _, e := range eit.Events {
			data[tableID][e.EventID] = map[uint8][]ShortEventDescriptor{}
			for _, d := range e.Descriptors {
				switch d.DescriptorTag() {
				case 0x4d:
					sed := d.ShortEventDescriptor()
					data[tableID][e.EventID][d.DescriptorTag()] = append(data[tableID][e.EventID][d.DescriptorTag()], sed)
				case 0x4e:

				}
			}
		}
		return data
	default:
		return nil
	}
}

func contain(s []uint8, v uint8) int {
	for i, k := range s {
		if v == k {
			return i
		}
	}
	return -1
}
