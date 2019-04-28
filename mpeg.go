package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	f, err := os.Open("test.ts")
	if err != nil {
		panic(err)
	}

	ids := []byte{0x4E, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f}

	data := TSFile(f, ids)
	fmt.Println(data)
	for _, s := range data.Section {
		for k, e := range s.Event {
			fmt.Println("=========================")
			fmt.Printf("%X", k)
			fmt.Println("\n")
			//fmt.Printf("%X", e.ExtendedEventDescriptor.TextChar)
			scanner := bufio.NewScanner(transform.NewReader(bytes.NewBuffer(AribStr(e.ExtendedEventDescriptor.TextChar)), japanese.ISO2022JP.NewDecoder()))
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			fmt.Println("\n")
			//fmt.Printf("%X", e.ShortEventDescriptor.EventNameChar)
			scanner = bufio.NewScanner(transform.NewReader(bytes.NewBuffer(AribStr(e.ShortEventDescriptor.EventNameChar)), japanese.ISO2022JP.NewDecoder()))
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			fmt.Println("\n")
			//fmt.Printf("%X", e.ShortEventDescriptor.TextChar)
			scanner = bufio.NewScanner(transform.NewReader(bytes.NewBuffer(AribStr(e.ShortEventDescriptor.TextChar)), japanese.ISO2022JP.NewDecoder()))
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			fmt.Println("\n")
		}
	}
}
