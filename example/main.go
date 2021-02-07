package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/roma8ok/egrul"
)

func main() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	s := egrul.Scrapper{Client: client}

	rows, err := s.Find("Рога и копыта")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", rows)
}
