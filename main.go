package main

import (
	"os"

	"github.com/backdround/deploy-configs/internal/realmain"
	"github.com/backdround/deploy-configs/pkg/logger"
)

func main() {
	l := logger.New()
	returnCode := realmain.Main(l)
	os.Exit(returnCode)
}
