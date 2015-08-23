package image

import (
	"bytes"
	"fmt"
	img "image"
	pj "image/jpeg"
	pn "image/png"
)

type Mapper interface {
	Get(x, y, scale int) []byte
	Center(scale int) *img.Point
}

var (
	MaxScale     = 5            // max scale
	ScaleModulus = float32(0.2) //scale modulus
	Size         = 256          //gird size px
)

func NewMapper(path string, format int) (Mapper, error) {
	switch format {
	case 0: //jpeg
		return newJPEG(path)
	case 1: //png
	case 2: //gif
	}
	return nil, fmt.Errorf("unknow img format #%d", format)
}

func ImageToBytes(i img.Image, format int) []byte {
	return saveToMemory(i, format)
}

func saveToMemory(i img.Image, format int) []byte {
	var buf bytes.Buffer

	switch format {
	case 0:
		pj.Encode(&buf, i, &pj.Options{Quality: 100})
	case 1:
		pn.Encode(&buf, i)
	}
	return buf.Bytes()
}

// bits rect
type bitsRect struct {
	x, y, width, height,
	moveTop, moveLeft, moveBottom, moveRight int
}

func (b bitsRect) toRect() img.Rectangle {
	return img.Rect(b.x, b.y, b.x+b.width, b.y+b.height)
}
