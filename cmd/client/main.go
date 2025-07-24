package main

import (
	"GophKeeper/internal/client/app"
	"GophKeeper/internal/common/buildlog"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	buildlog.Print(buildVersion, buildDate, buildCommit)

	app, err := app.New()
	if err != nil {
		panic(err)
	}

	app.Run()
}
