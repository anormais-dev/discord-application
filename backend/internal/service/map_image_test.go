package service

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand/v2"
	"testing"
)

func TestShrinkTurnsBigPNGIntoHalfSizeJPEG(t *testing.T) {
	// Ruído, como numa foto: PNG não comprime bem e o JPEG fica menor.
	rng := rand.New(rand.NewPCG(1, 2))
	src := image.NewRGBA(image.Rect(0, 0, 1600, 900))
	for y := range 900 {
		for x := range 1600 {
			src.Set(x, y, color.RGBA{uint8(rng.IntN(256)), uint8(y), uint8(x), 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, src)

	got := shrink(buf.Bytes(), "image/png")
	if got.ContentType != "image/jpeg" || len(got.Data) >= buf.Len() {
		t.Fatalf("esperava JPEG menor que o PNG, veio %s com %d bytes (PNG %d)", got.ContentType, len(got.Data), buf.Len())
	}
	cfg, err := jpeg.DecodeConfig(bytes.NewReader(got.Data))
	if err != nil || cfg.Width != 800 || cfg.Height != 450 {
		t.Fatalf("tamanho = %dx%d (%v), esperado 800x450", cfg.Width, cfg.Height, err)
	}
}

func TestShrinkKeepsWhatItCannotDecode(t *testing.T) {
	got := shrink([]byte("não é imagem"), "image/png")
	if string(got.Data) != "não é imagem" || got.ContentType != "image/png" {
		t.Fatalf("deveria devolver o original, veio %+v", got)
	}
}
