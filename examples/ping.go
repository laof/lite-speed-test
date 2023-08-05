package main

import (
	"fmt"

	"github.com/laof/lite-speed-test/ping"
)

var sss = `
ssr://dDMuZnJlZWdyYWRlbHkueHl6OjIyMjIyOmF1dGhfY2hhaW5fYTpub25lOnBsYWluOlpHOXVaM1JoYVhkaGJtY3VZMjl0Lz9vYmZzcGFyYW09JnJlbWFya3M9ZDNkM0xtUnZibWQwWVdsM1lXNW5MbU52YlNEbXRKdm1uWW5ubjdZ      
ss://YWVzLTI1Ni1nY206ZG9uZ3RhaXdhbmcuY29tQHd3dy5kb25ndGFpd2FuZzQuY29tOjIyMjIy#www.dongtaiwang.com+%e6%b4%9b%e6%9d%89%e7%9f%b6
`

func main() {
	dd, _ := ping.Test(sss)
	fmt.Println(dd)
}
