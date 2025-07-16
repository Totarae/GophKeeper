package common

import (
	"fmt"
	"go.uber.org/zap"
	"sync"
)

var Logger = zap.NewNop()
var once sync.Once

func Init(name, level string) {
	once.Do(func() {
		atomicLevel, err := zap.ParseAtomicLevel(level)
		if err != nil {
			panic(err)
		}

		config := zap.NewProductionConfig()
		config.Level = atomicLevel
		logger, err := config.Build()
		if err != nil {
			panic(err)
		}

		Logger = logger.Named(name)
	})
}

func PrintBuildInfo(version, date, commit string) {
	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}
