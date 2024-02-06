package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/debdutdeb/gopark/pkg/progressbar"
)

func main() {
	response, err := http.Get("https://cloud-images.ubuntu.com/releases/23.10/release/ubuntu-23.10-server-cloudimg-amd64.img")
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("ubuntu-23.04-server.img", os.O_CREATE|os.O_WRONLY, 0750)
	if err != nil {
		panic(err)
	}

	bar, err := progressbar.NewWriteProgressBar("ubuntu 23.04", response.ContentLength, f, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Downloading...")

	defer response.Body.Close()
	io.Copy(bar, response.Body)
}
