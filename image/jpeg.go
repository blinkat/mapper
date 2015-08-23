package image

import (
	pck "image/jpeg"
	"os"
)

type jpeg struct {
	image
}

func newJPEG(path string) (Mapper, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := pck.Decode(f)
	if err != nil {
		return nil, err
	}

	ret := &jpeg{}
	ret.bounds = m.Bounds()
	ret.source = m

	err = ret.initScale()
	if err != nil {
		return nil, err
	}
	ret.format = 0
	return ret, nil
}
