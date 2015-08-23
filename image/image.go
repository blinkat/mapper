package image

import (
	"golang.org/x/image/draw"
	img "image"
	//"math"
)

// type imageBit struct {
// 	bytes []byte
// 	cent  img.Point
// }

type imgBitsMap struct {
	bytes map[int]map[int][]byte
	cent  img.Point
}

type image struct {
	//bits   []img.Image
	bits   []*imgBitsMap
	source img.Image
	bounds img.Rectangle
	format int
}

func (r *image) initScale() error {
	r.bits = make([]*imgBitsMap, MaxScale)
	for i := 0; i < MaxScale; i++ {
		every := (float32(i) + 1.0) * ScaleModulus
		wide := int(float32(r.bounds.Dx()) * every)
		high := int(float32(r.bounds.Dy()) * every)

		dst := img.Rect(0, 0, wide, high)
		p := img.NewRGBA(dst)
		draw.ApproxBiLinear.Scale(p, dst, r.source, r.bounds, draw.Src, nil)
		r.bits[i] = r.cutScaled(p)
	}

	return nil
}

func (i *image) cutScaled(bit img.Image) *imgBitsMap {
	center := i.center(bit.Bounds())
	mx := center.X / Size
	my := center.Y / Size

	if center.X%Size != 0 {
		mx += 1
	}
	if center.Y%Size != 0 {
		my += 1
	}

	ret := make(map[int]map[int][]byte)
	for x := -mx; x <= mx; x++ {
		ret[x] = make(map[int][]byte)
		for y := -my; y <= my; y++ {
			bs := i.cutMap(x, y, bit)
			if bs != nil {
				ret[x][y] = bs
			}
		}
	}

	return &imgBitsMap{
		bytes: ret,
		cent:  center,
	}
}

func (i *image) Format() int {
	return i.format
}

func (i *image) Get(x, y, scale int) []byte {
	if scale < 0 || scale >= len(i.bits) {
		return nil
	}
	b := i.bits[scale].bytes
	if mx, ok := b[x]; ok {
		if my, ok := mx[y]; ok {
			return my
		}
	}

	return nil
}

func (i *image) cutMap(x, y int, bit img.Image) []byte {
	rect := img.Rect(x, y, x+Size, y+Size)
	ret := i.cutImg(bit, rect)

	if ret == nil {
		return nil
	}
	return saveToMemory(ret, i.format)
}

func (i *image) cutImg(bit img.Image, rect img.Rectangle) img.Image {
	r_draw := i.imageToDraw(bit.Bounds(), rect)
	start := i.startRect(r_draw, rect).toRect()
	mx := bit.Bounds().Dx()
	my := bit.Bounds().Dy()

	if r_draw.Dx() != 0 && r_draw.Dy() != 0 {
		if (r_draw.Min.X < 0 && r_draw.Min.X+Size < 0 && r_draw.Min.Y < 0 && r_draw.Min.Y+Size < 0) ||
			(r_draw.Min.X > mx && r_draw.Min.Y > my && r_draw.Max.X > mx && r_draw.Max.Y > my) {
			return nil
		}
		ret := img.NewRGBA(img.Rect(0, 0, Size, Size))
		draw.Draw(ret, start, bit, img.Pt(r_draw.Min.X, r_draw.Min.Y), draw.Src)
		return ret
	}
	return nil
}

func (i *image) startRect(drawRect, rect img.Rectangle) bitsRect {
	start := bitsRect{}
	start.x, start.width, start.moveLeft, start.moveRight = i.convertStartRect(drawRect.Min.X, drawRect.Dx())
	start.y, start.height, start.moveTop, start.moveBottom = i.convertStartRect(drawRect.Min.Y, drawRect.Dy())
	return start
}

func (i *image) convertStartRect(min, wide int) (int, int, int, int) {
	var coord, width, move1, move2 int
	if min == 0 && wide < Size {
		coord = Size - wide
		width = Size
		move1 = Size - coord
	} else if wide < Size {
		width = wide
		coord = 0
		move2 = Size - width
	} else {
		coord = 0
		width = Size
	}
	return coord, width, move1, move2
}

// image coord to draw coord
func (i *image) imageToDraw(bit img.Rectangle, rect img.Rectangle) img.Rectangle {
	center := i.center(bit)
	ret := img.ZR
	ret.Min.X, ret.Max.X = i.converMapCoord(rect.Min.X, center.X)
	ret.Min.Y, ret.Max.Y = i.converMapCoord(rect.Min.Y, center.Y)
	return ret
}

// coord = map coord
func (i *image) converMapCoord(coord, center int) (int, int) {
	retx, rety := 0, 0
	coord *= Size
	max := center * 2
	if coord < 0 {
		retx = center + coord
	} else {
		retx = center + coord
	}

	if retx < 0 {
		rety = retx + Size
		retx = 0
	} else if retx > max {
		rety = max
	} else {
		rety = retx + Size
	}
	return retx, rety
}

func (i *image) center(bit img.Rectangle) img.Point {
	w := bit.Dx() / 2
	h := bit.Dy() / 2
	if bit.Dx()%2 != 0 {
		w -= 1
	}
	if bit.Dy()%2 != 0 {
		h -= 1
	}
	return img.Pt(w, h)
}

func (i *image) Center(scale int) *img.Point {
	if scale >= len(i.bits) || scale < 0 {
		return nil
	}
	return &i.bits[scale].cent
}
