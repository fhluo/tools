package main

import (
	"golang.org/x/exp/slog"
	"os"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
}

func main() {
	Execute()
}
