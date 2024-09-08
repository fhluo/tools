//go:build windows

package main

import (
	"iter"
	"log"
	"log/slog"
	"os"
	"reflect"
	"slices"

	"github.com/spf13/cobra"
	"golang.org/x/sys/windows"
)

var ProcSetSuspendState = windows.NewLazyDLL("PowrProf.dll").NewProc("SetSuspendState")

type SetSuspendState struct {
	Hibernate           bool
	Force               bool
	DisableWakeupEvents bool
}

func (s SetSuspendState) fields() iter.Seq[reflect.Value] {
	return func(yield func(reflect.Value) bool) {
		v := reflect.ValueOf(s)
		for i := 0; i < v.NumField(); i++ {
			if !yield(v.Field(i)) {
				break
			}
		}
	}
}

func (s SetSuspendState) Call() error {
	a := slices.Collect(func(yield func(uintptr) bool) {
		for field := range s.fields() {
			if field.Bool() {
				yield(1)
			} else {
				yield(0)
			}
		}
	})

	r, _, err := ProcSetSuspendState.Call(a...)
	if r == 0 {
		return err
	}

	return nil
}

var rootCmd = cobra.Command{
	Use:  "hibernate",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return setSuspendState.Call()
	},
}

var setSuspendState SetSuspendState

func init() {
	log.SetFlags(0)

	rootCmd.PersistentFlags().BoolP("help", "", false, "help for hibernate")
	rootCmd.Flags().BoolVarP(&setSuspendState.Hibernate, "hibernate", "h", true, "hibernate or suspend")
	rootCmd.Flags().BoolVarP(&setSuspendState.DisableWakeupEvents, "disableWakeupEvents", "d", false, "disables all wake events")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to execute root command", "err", err)
		os.Exit(1)
	}
}
