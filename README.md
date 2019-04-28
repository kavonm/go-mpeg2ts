# go-mpeg2ts
*開発中のためpackage mainで実装中

MPEG2TSのParse
```
ARIB STD-B10
ARIB STD-B24 (text関係)
ISO/IEC 13818-1
```

### 使い方
`mpeg.go`を参照

`test.ts`ファイルからEIT取得
```
func main() {
    f, err := os.Open("test.ts")
    if err != nil {
	panic(err)
    }

    ids := []byte{0x4E, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f}

    data := TSFile(f, ids)
    /*
    data = map[
        table_id: map[
            event_id: map[
                descriptor_tag: []
            ]
        ]
    ]
    */
    fmt.Println(data)
}
```

### 現状
- ARIBの文字コードの複合処理はほぼ実装できてる。(コードのリファクタリングが必要)
- 特殊文字のMAPはまだ実装してない。
- DescriptorはShortEventDescriptorとExtendedEventDescriptorだけ

### 課題
- 完全な番組表の取得(現在EITのShortEventDescriptorのみ)
- ARIB STD-B25も実装したい(暗号化と復号化)
