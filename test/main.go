package main

import (
	"fmt"

	"github.com/laof/lite-speed-test/data"
)

func main() {
	obj := data.Nodes{Url: "fjson", Max: 1}
	data, err := obj.Get()

	if err != nil {
		return
	}

	// fmt.Println(data)
	res, _ := obj.Test(data)

	fmt.Println(res.SuccessNodes)
	fmt.Println(res.ErrorServers)

}
