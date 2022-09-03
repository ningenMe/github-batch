package main

import (
	"flag"

	"github.com/ningenme/nina-batch/cmd"
)

func main() {
	var t = flag.String("t", "hoge", "personal access token")
	var s = flag.String("s", "", "start date")
	var e = flag.String("e", "", "end date")

	loginUserContributionUsecase := cmd.LoginUserContributionUsecase{}
	flag.Parse()
	loginUserContributionUsecase.Execute(*t, *s, *e)
}
