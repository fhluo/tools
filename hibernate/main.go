//go:build windows

package main

import (
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sys/windows"
	"log"
	"os"
)

var ProcSetSuspendState = windows.NewLazyDLL("PowrProf.dll").NewProc("SetSuspendState")

func SetSuspendState(hibernate bool, disableWakeupEvents bool) error {
	var (
		hibernate_           uintptr
		force_               uintptr
		disableWakeupEvents_ uintptr
	)

	if hibernate {
		hibernate_ = 1
	}
	if disableWakeupEvents {
		disableWakeupEvents_ = 1
	}

	r, _, err := ProcSetSuspendState.Call(
		hibernate_,
		force_,
		disableWakeupEvents_,
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
