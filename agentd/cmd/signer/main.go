package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	file := flag.String("file", "", "path to template body file")
	secret := flag.String("secret", "demo-template-hmac", "hmac secret")
	flag.Parse()

	if *file == "" {
		fmt.Println("Usage: signer -file body.txt")
		return
	}

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	mac := hmac.New(sha256.New, []byte(*secret))
	mac.Write(data)
	sum := mac.Sum(nil)
	fmt.Println(hex.EncodeToString(sum))
}
