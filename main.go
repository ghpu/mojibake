package main

import (
	"bufio"
	"fmt"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/ianaindex"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func checknil(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var filename string
	var encoding string
	detector := chardet.NewTextDetector()

	fmt.Println("Mojibake v0.1 - UTF-8 file encoding converter")

	fmt.Printf("Input file : ")
	reader := bufio.NewReader(os.Stdin)
	filename, _ = reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	dat, err := ioutil.ReadFile(filename)
	checknil(err)

	result, err := detector.DetectBest(dat)
	checknil(err)

	fmt.Printf("Input encoding (%s) : ", result.Charset)
	encoding, _ = reader.ReadString('\n')
	encoding = strings.TrimSpace(encoding)
	if len(encoding) == 0 {
		encoding = result.Charset
	}

	enc, err := ianaindex.IANA.Encoding(encoding)
	checknil(err)
	name, err := ianaindex.IANA.Name(enc)
	checknil(err)
	fmt.Printf("Selected encoding : %s\n", name)

	f, err := os.Open(filename)
	checknil(err)

	outputname := "utf8." + filename
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Failure to convert file %s\n", filename)
			os.Remove(outputname)
		}
	}()
	r := enc.NewDecoder().Reader(f)

	out, err := os.Create(outputname)
	checknil(err)
	io.Copy(out, r)
	out.Close()
	f.Close()

	fmt.Printf("File successfully written to : %s\n", outputname)

}
