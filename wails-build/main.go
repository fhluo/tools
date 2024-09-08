package main

import (
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wails-build",
	Short: "wails build",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		paths := slices.Collect(func(yield func(string) bool) {
			for pattern := range slices.Values(patterns) {
				r, err := filepath.Glob(pattern)
				if err != nil {
					slog.Warn(err.Error())
					continue
				}

				for p := range slices.Values(r) {
					p, err = filepath.Abs(p)
					if err != nil {
						slog.Warn(err.Error())
						continue
					}

					if !yield(p) {
						return
					}
				}
			}
		})

		err = os.Setenv("path", strings.Join(paths, ";")+";"+os.Getenv("path"))
		if err != nil {
			return err
		}
		slog.Debug("env", "path", os.Getenv("path"))

		args := slices.Collect(func(yield func(string) bool) {
			yield("build")

			if useUPX {
				yield("-upx")
			}

			if useNSIS {
				yield("-nsis")
			}
		})

		// 执行 wails build 命令
		command := exec.Command("wails", args...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		slog.Info(command.String(), "path", os.Getenv("path"))

		return command.Run()
	},
}

var (
	useUPX   bool     // 是否使用 UPX
	useNSIS  bool     // 是否使用 NSIS
	patterns []string // 要添加到环境变量中的路径
)

func init() {
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	))
	log.SetFlags(0)

	rootCmd.Flags().BoolVar(&useUPX, "upx", false, "use UPX")
	rootCmd.Flags().BoolVar(&useNSIS, "nsis", false, "use NSIS")
	rootCmd.Flags().StringSliceVar(&patterns, "path", nil, "paths to add to the environment variable")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
