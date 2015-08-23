package image

import (
	img "image"
	"image/color"
)

var warringColor = color.NRGBA{R: 230, G: 53, B: 53, A: 255}

func DrawBottomWarring() img.Image {
	ret := img.NewNRGBA(img.Rect(0, 0, 1, 10))
	for i := 0; i < 10; i++ {
		c := warringColor
		c.A = uint8(255 * (float32(i+1) * 0.075))
		ret.SetNRGBA(0, i, c)
	}
	return ret
}

func DrawTopWarring() img.Image {
	ret := img.NewNRGBA(img.Rect(0, 0, 1, 10))
	for i := 9; i >= 0; i-- {
		c := warringColor
		c.A = uint8(255 * (float32(i+1) * 0.075))
		ret.SetNRGBA(0, 9-i, c)
	}
	return ret
}

func DrawRightWarring() img.Image {
	ret := img.NewNRGBA(img.Rect(0, 0, 10, 1))
	for i := 9; i >= 0; i-- {
		c := warringColor
		c.A = uint8(255 * (float32(i+1) * 0.075))
		ret.SetNRGBA(i, 0, c)
	}
	return ret
}

func DrawLeftWarring() img.Image {
	ret := img.NewNRGBA(img.Rect(0, 0, 10, 1))
	for i := 0; i < 10; i++ {
		c := warringColor
		c.A = uint8(255 * (float32(i+1) * 0.075))
		ret.SetNRGBA(9-i, 0, c)
	}
	return ret
}
