package util

import (
	"compress/flate"
	"io/ioutil"
	"log"
	"os"

	"github.com/liuzl/gocc"
)

func Convert(conversion string, in string) (string, error) {
	converter, err := gocc.New(conversion)
	if err != nil {
		log.Fatal(err)
	}
	return converter.Convert(in)
}

func Trad2SimpleConvert(in string) (string, error) {
	return Convert("t2s", in)
}

func DecompressFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fz := flate.NewReader(f)

	defer fz.Close()

	bs, err := ioutil.ReadAll(fz)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
