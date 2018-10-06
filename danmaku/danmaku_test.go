package danmaku

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestParseXML(t *testing.T) {
	f, _ := os.Open(`./example.xml`)
	bytes, _ := ioutil.ReadAll(f)
	bc, err := ParseXML(bytes)
	if err != nil {
		t.Error("Parsing danmaku xml error:", err)
	}
	log.Printf("%+v\n", bc.Comments[0])
}

func TestParseJson(t *testing.T) {
	f, _ := os.Open(`./example.info`)
	bytes, _ := ioutil.ReadAll(f)
	info, err := ParseInfo(bytes)
	if err != nil {
		t.Error("Parsing info error:", err)
	}
	log.Printf("%+v\n", info)
}

func TestConvert2RespComment(t *testing.T) {
	f, _ := os.Open(`./example.xml`)
	defer f.Close()
	bytes, _ := ioutil.ReadAll(f)
	bc, _ := ParseXML(bytes)

	log.Printf("%+v\n", bc.Comments[0])
	f1, _ := os.Open(`./example.info`)
	defer f1.Close()
	bytes, _ = ioutil.ReadAll(f1)
	info, _ := ParseInfo(bytes)

	rc, err := Convert2RespComment(bc, info)
	if err != nil {
		log.Println("error:", err)
		return
	}
	fmt.Println(rc)

	// log.Println(rc)
}
