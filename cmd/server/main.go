package main

import (
	"GophKeeper/internal/common"
	"GophKeeper/internal/server/app"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	common.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	app, err := app.New()
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
