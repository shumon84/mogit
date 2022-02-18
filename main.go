package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/shumon84/mogit/inner/object"
)

func main() {
	b, err := object.NewBlobFromPath("inner/object/blob.go")
	if err != nil {
		log.Fatal(err)
	}
	s, err := b.SHA1()
	if err != nil {
		log.Fatal(err)
	}
	decode, err := b.Decode()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("#Type  :", b.Type())
	fmt.Println("#Decode:", string(decode))
	fmt.Println("#SHA1  :", hex.EncodeToString(s))

	f, err := os.OpenFile("hoge.zlib", os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	encode, err := b.Encode()
	if err != nil {
		log.Fatal(err)
	}
	f.Write(encode)
}
