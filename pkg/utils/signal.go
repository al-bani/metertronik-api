package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalChannel() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	return sigChan
}
