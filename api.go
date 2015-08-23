package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/blinkat/mapper/image"
	"github.com/blinkat/mapper/resource"
	"net/http"
	"strconv"
)

type ImageFormat int

var mapper image.Mapper

const namespace = "blinkat/mapper:"

const (
	JPEG = ImageFormat(iota)
	PNG
	GIF
)

var (
	warringTop    = image.ImageToBytes(image.DrawTopWarring(), int(PNG))
	warringLeft   = image.ImageToBytes(image.DrawLeftWarring(), int(PNG))
	warringBottom = image.ImageToBytes(image.DrawBottomWarring(), int(PNG))
	warringRight  = image.ImageToBytes(image.DrawRightWarring(), int(PNG))
)

func Handler(response http.ResponseWriter, request *http.Request) {
	var t = request.FormValue("type")
	var err error
	switch t {
	case "jquery":
		err = writehttp(response, resource.JQuery)
	case "get":
		err = handleGet(response, request)
	case "init-data":
		err = handleInitData(response, request)
	case "cent":
		err = handleCent(response, request)
	case "warring-img":
		err = handleWarringImg(response, request)
	case "mouse-wheel":
		err = writehttp(response, resource.MouseWheel)
	case "js":
		err = writehttp(response, resource.MapperJS)
	case "css":
		response.Header().Add("content-type", "text/css")
		err = writehttp(response, resource.MapperCss)
	default:
		response.Write([]byte(fmt.Sprintf(namespace+" unknow type '%s'", t)))
	}

	if err != nil {
		response.Write([]byte(fmt.Sprint(err)))
	}
}

func handleWarringImg(response http.ResponseWriter, request *http.Request) error {
	t := request.FormValue("img-type")
	if t != "" {
		var b []byte
		switch t {
		case "top":
			b = warringTop
		case "left":
			b = warringLeft
		case "right":
			b = warringRight
		case "bottom":
			b = warringBottom
		}
		if b != nil {
			response.Header().Add("content-type", "image/jpg")
			response.Write(b)
			return nil
		}
	}

	return fmt.Errorf("unknow img-type '%s", t)
}

func handleCent(response http.ResponseWriter, request *http.Request) error {
	tmp := request.FormValue("scale")
	tmpN, err := strconv.ParseInt(tmp, 10, 32)
	if err != nil {
		return err
	}
	scale := int(tmpN)
	p := mapper.Center(scale)
	if p != nil {
		var data struct {
			X int `json:"x"`
			Y int `json:"y"`
		}

		data.X = p.X
		data.Y = p.Y

		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		response.Write(b)
		return nil
	}
	return fmt.Errorf("faild scale #%d", scale)
}

func handleInitData(response http.ResponseWriter, request *http.Request) error {
	var data struct {
		Size int `json:"size"`
	}

	data.Size = image.Size
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}

	response.Write(bs)
	return nil
}

func handleGet(response http.ResponseWriter, request *http.Request) error {
	tmp := request.FormValue("scale")
	tmpN, err := strconv.ParseInt(tmp, 10, 32)
	if err != nil {
		return err
	}
	scale := int(tmpN)

	tmp = request.FormValue("x")
	tmpN, err = strconv.ParseInt(tmp, 10, 32)
	if err != nil {
		return err
	}
	x := int(tmpN)

	tmp = request.FormValue("y")
	tmpN, err = strconv.ParseInt(tmp, 10, 32)
	if err != nil {
		return err
	}
	y := int(tmpN)

	bits := mapper.Get(x, y, scale)
	if bits == nil {
		response.WriteHeader(404)
	} else {
		response.Write(bits)
	}
	return nil
}

func writehttp(res http.ResponseWriter, t string) error {
	_, err := res.Write([]byte(t))
	return err
}

func SetMapImage(path string, format ImageFormat) error {
	m, e := image.NewMapper(path, int(format))
	if e != nil {
		return e
	}
	mapper = m
	return nil
}
