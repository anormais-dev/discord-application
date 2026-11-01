package service

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
)

// Acima disso a arte é reduzida à metade: o modal mostra o splash com no máximo 560px.
const maxImageWidth = 1200

// shrink troca o PNG da API (o splash tem 2,3 MB) por um JPEG bem menor, reduzindo as artes
// grandes à metade. Se algo falhar, devolve a imagem original.
func shrink(data []byte, contentType string) MapImage {
	original := MapImage{Data: data, ContentType: contentType}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return original
	}
	if img.Bounds().Dx() > maxImageWidth {
		img = half(img)
	}
	var out bytes.Buffer
	if err := jpeg.Encode(&out, img, &jpeg.Options{Quality: 82}); err != nil || out.Len() >= len(data) {
		return original
	}
	return MapImage{Data: out.Bytes(), ContentType: "image/jpeg"}
}

// half reduz a imagem à metade tirando a média de cada bloco de 2x2 pixels.
func half(src image.Image) image.Image {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx()/2, b.Dy()/2))
	for y := 0; y < dst.Rect.Dy(); y++ {
		for x := 0; x < dst.Rect.Dx(); x++ {
			var r, g, bl, a uint32
			for _, d := range [4][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}} {
				pr, pg, pb, pa := src.At(b.Min.X+2*x+d[0], b.Min.Y+2*y+d[1]).RGBA()
				r, g, bl, a = r+pr, g+pg, bl+pb, a+pa
			}
			i := dst.PixOffset(x, y)
			dst.Pix[i], dst.Pix[i+1], dst.Pix[i+2], dst.Pix[i+3] = uint8(r/4>>8), uint8(g/4>>8), uint8(bl/4>>8), uint8(a/4>>8)
		}
	}
	return dst
}

// Warm baixa todas as artes logo que o servidor sobe, para o primeiro sorteio não esperar a API.
func (s *MapService) Warm(ctx context.Context) {
	if _, err := s.CompetitiveMaps(ctx); err != nil {
		log.Printf("mapas: não deu para pré-carregar as imagens: %v", err)
		return
	}
	for _, img := range s.Images() {
		for _, kind := range []string{"card", "splash"} {
			if _, err := s.Image(ctx, img.UUID, kind); err != nil {
				log.Printf("mapas: imagem %s/%s: %v", img.UUID, kind, err)
			}
		}
	}
}
