//go:build windows

package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/sys/windows"
)

var ProcSetSuspendState = windows.NewLazyDLL("PowrProf.dll").NewProc("SetSuspendState")

type Bool bool

func (b Bool) Uintptr() uintptr {
	if b {
		return 1
	} else {
		return 0
	}
}

func SetSuspendState(hibernate bool, disableWakeupEvents bool) error {
	force := false

	r, _, err := ProcSetSuspendState.Call(
		Bool(hibernate).Uintptr(),
		Bool(force).Uintptr(),
		Bool(disableWakeupEvents).Uintptr(),
	)
	if r == 0 {
		return err
	}

	return nil
}

var rootCmd = cobra.Command{
	Use:  "hibernate",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return SetSuspendState(hibernate, disableWakeupEvents)
	},
}

var (
	hibernate           bool
	disableWakeupEvents bool
)

func init() {
	log.SetFlags(0)

	rootCmd.PersistentFlags().BoolP("help", "", false, "help for hibernate")
	rootCmd.Flags().BoolVarP(&hibernate, "hibernate", "h", true, "hibernate or suspend")
	rootCmd.Flags().BoolVarP(&disableWakeupEvents, "disableWakeupEvents", "d", false, "disables all wake events")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to execute root command", "err", err)
		os.Exit(1)
	}
}
