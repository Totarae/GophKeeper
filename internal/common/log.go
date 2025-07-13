package common

import "fmt"

func PrintBuildInfo(version, date, commit string) {
	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}
