package main

import (
	"flag"

	"github.com/ningenme/nina/cmd"
)

func main() {
	var t = flag.String("t", "hoge", "personal access token")
	flag.Parse()
	cmd.Execute(*t)
}
