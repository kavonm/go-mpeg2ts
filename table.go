package main

//the first byte is pointer_field
type Table []byte

func (t Table) TableID() uint8 {
	return t[1]
}

func (t Table) Parser(ids []byte, ts *TS) {
	tableID := t.TableID()
	if b := contain(ids, tableID); b == -1 {
		return
	}
	switch tableID {
	case 0x4E, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56,
		0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f:
		eit := t.EventInformationTable()
		//ts.Section[tableID] = Section{}
		for _, e := range eit.Events {
			ts.Section[tableID] = Section{map[uint16]Event{}}
			for _, d := range e.Descriptors {
				switch d.DescriptorTag() {
				case 0x4d:
					sed := d.ShortEventDescriptor()
					event := ts.Section[tableID].Event[e.EventID]
					event.ShortEventDescriptor = sed
					ts.Section[tableID].Event[e.EventID] = event
				case 0x4e:
					eed := d.ExtendedEventDescriptor()
					event := ts.Section[tableID].Event[e.EventID]
					event.ExtendedEventDescriptor = eed
					ts.Section[tableID].Event[e.EventID] = event
				}
			}
		}
		return
	default:
		return
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
