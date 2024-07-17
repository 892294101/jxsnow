package main

import (
	"fmt"
	"github.com/892294101/jxsnow"
)

func main() {
	g, _ := jxsnow.NewGenerator(1)
	n, _ := g.Generate()
	fmt.Println(n)

}
