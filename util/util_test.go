package util

import (
	"testing"
)

func TestConvert(t *testing.T) {
	out, err := Convert("t2s", `女神異聞錄5 [25]`)
	if err != nil {
		t.Error("Error in conversion")
	}
	if out != `女神异闻录5 [25]` {
		t.Error("Conversion error", out)
	}
}

func TestTrad2SimpleConvert(t *testing.T) {
	out, err := Trad2SimpleConvert(`女神異聞錄5 [25]`)
	if err != nil {
		t.Error("Error in conversion")
	}
	if out != `女神异闻录5 [25]` {
		t.Error("Conversion error", out)
	}
}

func TestDecompress(t *testing.T) {
	_, err := DecompressFile(`./example.raw`)
	if err != nil {
		t.Error("Failed to decompress example.raw", err)
	}
}
