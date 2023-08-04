package main

import (
	"fmt"

	"github.com/laof/lite-speed-test/ping"
)

var sss = `
ssr://OTEuMiZnNwYXJhbT0mcmVtYXJrcz01TC1FNTcyWDVwYXZRUSZncm91cD1URzVqYmk1dmNtYw
ssr://MTk1LjEzMy4xMS44OjIwMDg3Om9yaWdpbjLUU1NzJYNXBhdlFnJmdyb3VwPVRHNWpiaTV2Y21j
`

func main() {
	dd, _ := ping.Test(sss)
	fmt.Println(dd)
}
