package main

type EventInformationSection struct {
	TableID                  uint8  //8bits
	SectionSyntaxIndicator   uint8  //1bit
	SectionLength            uint16 //12bits
	ServiceID                uint16 //16bits
	VersionNumber            uint8  //5bits
	CurrentNextIndicator     uint8  //1bits
	SectionNumber            uint8  //8bits
	LastSectionNumber        uint8  //8bits
	TransportStreamID        uint16 //16bits
	OriginalNetworkID        uint16 //16bits
	SegmentLastSectionNumber uint8  //8bits
	LastTableID              uint8  //8bits
	Events                   []EventInformationSectionEvent
}

type EventInformationSectionEvent struct {
	EventID               uint16 //16bits
	StartTime             []byte //40bits
	Duration              []byte //24bits
	RunningStatus         uint8  //3bits
	FreeCAMode            uint8  //1bit
	DescriptionLoopLength uint16 //12bits
	Descriptors           []Descriptor
}

func (t Table) EventInformationTable() EventInformationSection {
	tableID := t[1]
	sectionSyntaxIndicator := (t[2] & 0x80) >> 7
	sectionLength := uint16(t[2]&0x0f)<<8 | uint16(t[3])
	serviceID := uint16(t[4])<<8 | uint16(t[5])
	versionNumber := (t[6] & 0x3e) >> 1
	currentNextIndicator := t[6] & 0x01
	sectionNumber := t[7]
	lastSectionNumber := t[8]
	transportStreamID := uint16(t[9])<<8 | uint16(t[10])
	originalNetworkID := uint16(t[11])<<8 | uint16(t[12])
	segmentLastSectionNumber := t[13]
	lastTableID := t[14]

	events := []EventInformationSectionEvent{}
	offset := uint16(15)
	for {
		if offset >= sectionLength-15 {
			break
		}
		eventID := uint16(t[offset])<<8 | uint16(t[offset+1])
		startTime := t[offset+2 : offset+7]
		duration := t[offset+7 : offset+10]
		runningStatus := (t[offset+10] & 0xe0) >> 5
		freeCAMode := (t[offset+10] & 0x10) >> 4
		descriptionLoopLength := uint16(t[offset+10]&0x0f)<<8 | uint16(t[offset+11])
		descriptorsData := t[offset+12 : offset+12+descriptionLoopLength]
		descriptors := []Descriptor{}
		descriptorOffset := uint16(0)
		for {
			if descriptorOffset >= descriptionLoopLength {
				break
			}
			descriptorLength := descriptorsData[descriptorOffset+1]
			descriptor := Descriptor(descriptorsData[descriptorOffset : descriptorOffset+uint16(descriptorLength)+2])
			descriptors = append(descriptors, descriptor)
			descriptorOffset += uint16(descriptorLength) + 2
		}
		event := EventInformationSectionEvent{
			EventID:               eventID,
			StartTime:             startTime,
			Duration:              duration,
			RunningStatus:         runningStatus,
			FreeCAMode:            freeCAMode,
			DescriptionLoopLength: descriptionLoopLength,
			Descriptors:           descriptors,
		}
		events = append(events, event)
		offset += 12 + descriptionLoopLength
	}

	return EventInformationSection{
		TableID:                  tableID,
		SectionSyntaxIndicator:   sectionSyntaxIndicator,
		SectionLength:            sectionLength,
		ServiceID:                serviceID,
		VersionNumber:            versionNumber,
		CurrentNextIndicator:     currentNextIndicator,
		SectionNumber:            sectionNumber,
		LastSectionNumber:        lastSectionNumber,
		TransportStreamID:        transportStreamID,
		OriginalNetworkID:        originalNetworkID,
		SegmentLastSectionNumber: segmentLastSectionNumber,
		LastTableID:              lastTableID,
		Events:                   events,
	}
}
