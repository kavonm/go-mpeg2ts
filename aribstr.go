package main

const (
	Kanji = iota
	Alphanumeric
	Hiragana
	Katakana
	MosaicA
	MosaicB
	MosaicC
	MosaicD
	PropAlphanumeric
	PropHiragana
	PropKatakana
	JISX0201Katakana
	JISKanjiPlane1
	JISKanjiPlane2
	AdditionalSymbols
	Unsupported
)

var CodeGSet = map[byte]int{
	0x42: Kanji,
	0x4a: Alphanumeric,
	0x30: Hiragana,
	0x31: Katakana,
	0x32: MosaicA,
	0x33: MosaicB,
	0x34: MosaicC,
	0x35: MosaicD,
	0x36: PropAlphanumeric,
	0x37: PropHiragana,
	0x38: PropKatakana,
	0x39: JISKanjiPlane1,
	0x3a: JISKanjiPlane2,
	0x3b: AdditionalSymbols,
	0x49: JISX0201Katakana,
}

const (
	EscSeqASCII = iota
	EscSeqZenkaku
	EscSeqHankaku
)

var EscSeq = map[int][]byte{
	0: []byte{0x1b, 0x28, 0x42},
	1: []byte{0x1b, 0x24, 0x42},
	2: []byte{0x1b, 0x28, 0x49},
}

type Buffer struct {
	G0 int
	G1 int
	G2 int
	G3 int
}

type Graphic struct {
	GL int
	GR int
}

func AribStr(b []byte) []byte {
	Buffer := Buffer{Kanji, Alphanumeric, Hiragana, Katakana}
	Graphic := Graphic{Buffer.G0, Buffer.G2}
	jis := []byte{}
	var temp bool
	var buf int
	var esc int
	for {
		if len(b) == 0 {
			return jis
		}
		data := b[0]
		b = b[1:]
		if (0x21 <= data && data <= 0x7e) || (0xa1 <= data && data <= 0xfe) {
			//GL or GR table
			firstChar := data
			var secondChar byte = 0x00
			var code int
			if 0x21 <= data && data <= 0x7e {
				code = Graphic.GL
			} else {
				code = Graphic.GR
			}
			if code == Kanji || code == JISKanjiPlane1 || code == JISKanjiPlane2 || code == AdditionalSymbols {
				secondChar = b[0]
				b = b[1:]
			}
			if 0xa1 <= firstChar && firstChar <= 0xfe {
				firstChar = firstChar & 0x7f
				secondChar = secondChar & 0x7f
			}
			if code == Kanji || code == JISKanjiPlane1 || code == JISKanjiPlane2 {
				if esc != EscSeqZenkaku {
					esc = EscSeqZenkaku
					jis = append(jis, EscSeq[esc]...)
				}
				jis = append(jis, firstChar, secondChar)
			} else if code == Alphanumeric || code == PropAlphanumeric {
				if esc != EscSeqASCII {
					esc = EscSeqASCII
					jis = append(jis, EscSeq[esc]...)
				}
				jis = append(jis, firstChar)
			} else if code == Hiragana || code == PropHiragana {
				if esc != EscSeqZenkaku {
					esc = EscSeqZenkaku
					jis = append(jis, EscSeq[esc]...)
				}
				if firstChar >= 0x77 {
					jis = append(jis, 0x21, firstChar)
				} else {
					jis = append(jis, 0x24, firstChar)
				}
			} else if code == Katakana || code == PropKatakana {
				if esc != EscSeqZenkaku {
					esc = EscSeqZenkaku
					jis = append(jis, EscSeq[esc]...)
				}
				if firstChar >= 0x77 {
					jis = append(jis, 0x21, firstChar)
				} else {
					jis = append(jis, 0x24, firstChar)
				}
			} else if code == JISX0201Katakana {
				if esc != EscSeqHankaku {
					esc = EscSeqHankaku
					jis = append(jis, EscSeq[esc]...)
				}
				jis = append(jis, firstChar)
			} else if code == AdditionalSymbols {

			}
			if temp {
				Graphic.GL, temp = buf, false
			}
		} else {
			switch data {
			case 0x1b:
				//ESCb
				data = b[0]
				b = b[1:]
				switch data {
				case 0x6e:
					Graphic.GL = Buffer.G2
				case 0x6f:
					Graphic.GL = Buffer.G3
				case 0x7e:
					Graphic.GR = Buffer.G1
				case 0x7d:
					Graphic.GR = Buffer.G2
				case 0x7c:
					Graphic.GR = Buffer.G3
				case 0x28:
					data = b[0]
					b = b[1:]
					switch data {
					case 0x4a, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x49:
						Buffer.G0 = CodeGSet[data]
					case 0x20:
						data = b[0]
						b = b[1:]
						switch data {
						case 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x70:
							Buffer.G0 = Unsupported
						}
					}
				case 0x29:
					data = b[0]
					b = b[1:]
					switch data {
					case 0x4a, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x49:
						Buffer.G1 = CodeGSet[data]
					case 0x20:
						data = b[0]
						b = b[1:]
						switch data {
						case 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x70:
							Buffer.G1 = Unsupported
						}
					}
				case 0x2a:
					data = b[0]
					b = b[1:]
					switch data {
					case 0x4a, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x49:
						Buffer.G2 = CodeGSet[data]
					case 0x20:
						data = b[0]
						b = b[1:]
						switch data {
						case 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x70:
							Buffer.G2 = Unsupported
						}
					}
				case 0x2b:
					data = b[0]
					b = b[1:]
					switch data {
					case 0x4a, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x49:
						Buffer.G3 = CodeGSet[data]
					case 0x20:
						data = b[0]
						b = b[1:]
						switch data {
						case 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x70:
							Buffer.G3 = Unsupported
						}
					}
				case 0x24:
					data = b[0]
					b = b[1:]
					switch data {
					case 0x42, 0x39, 0x3a, 0x3b:
						Buffer.G0 = CodeGSet[data]
					case 0x28:
						data = b[0]
						b = b[1:]
						switch data {
						case 0x20:
							data = b[0]
							b = b[1:]
							Buffer.G0 = Unsupported
						}
					case 0x29:
						data = b[0]
						b = b[1:]
						switch data {
						case 0x42, 0x39, 0x3a, 0x3b:
							Buffer.G1 = CodeGSet[data]
						case 0x20:
							data = b[0]
							b = b[1:]
							Buffer.G1 = Unsupported
						}
					case 0x2a:
						data = b[0]
						b = b[1:]
						switch data {
						case 0x42, 0x39, 0x3a, 0x3b:
							Buffer.G2 = CodeGSet[data]
						case 0x20:
							data = b[0]
							b = b[1:]
							Buffer.G2 = Unsupported
						}
					case 0x2b:
						data = b[0]
						b = b[1:]
						switch data {
						case 0x42, 0x39, 0x3a, 0x3b:
							Buffer.G3 = CodeGSet[data]
						case 0x20:
							data = b[0]
							b = b[1:]
							Buffer.G3 = Unsupported
						}
					}
				}
			case 0x0f:
				//LS0
				Graphic.GL = Buffer.G0
			case 0x0e:
				//LS1
				Graphic.GL = Buffer.G1
			case 0x19:
				buf, Graphic.GL = Graphic.GL, Buffer.G2
				temp = true
			case 0x1d:
				buf, Graphic.GL = Graphic.GL, Buffer.G3
				temp = true
			}
		}
	}
}
