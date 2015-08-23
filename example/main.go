package main

import (
	"bufio"
	"fmt"
	"github.com/blinkat/mapper"
	"io/ioutil"
	"net/http"
	"os"
)

func getFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	ret, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}
	return ret
}

func example(res http.ResponseWriter, req *http.Request) {
	ret := getFile("./example.html")
	if ret == nil {
		res.WriteHeader(404)
	} else {
		res.Write(ret)
	}
}

func example_content(res http.ResponseWriter, req *http.Request) {
	var ret []byte
	t := req.FormValue("id")
	switch t {
	case "css":
		ret = getFile("./example.css")
		res.Header().Add("content-type", "text/css")
	case "js":
		ret = getFile("./example.js")
	case "map-js":
		ret = getFile("./test/mapper.js")
	case "map-css":
		ret = getFile("./test/mapper.css")
		res.Header().Add("content-type", "text/css")
	}
	if ret == nil {
		res.WriteHeader(404)
	} else {
		res.Write(ret)
	}
}

func main() {
	err := mapper.SetMapImage("./example-map.jpg", mapper.JPEG)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/handler", mapper.Handler)
	http.HandleFunc("/content", example_content)
	http.HandleFunc("/", example)
	http.HandleFunc("/example", example)

	go http.ListenAndServe(":1230", nil)

	fmt.Println("listen and server port:1230")
	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		cmd := string(data)
		if cmd == "quit" {
			break
		}
	}
}
