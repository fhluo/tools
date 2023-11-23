package main

import (
	"github.com/samber/lo"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wails-build",
	Short: "wails build",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		paths := make([]string, 0, len(patterns))

		// 将路径添加到环境变量中
		for _, pattern := range patterns {
			r, err := filepath.Glob(pattern)
			if err != nil {
				return err
			}

			// 获取绝对路径
			for i := range r {
				r[i], err = filepath.Abs(r[i])
				if err != nil {
					return err
				}
			}

			paths = append(paths, r...)
		}

		err = os.Setenv("path", strings.Join(paths, ";")+";"+os.Getenv("path"))
		if err != nil {
			return err
		}

		// 执行 wails build 命令
		command := exec.Command(
			"wails", "build",
			lo.If(useUPX, "-upx").Else(""),
			lo.If(useNSIS, "-nsis").Else(""),
		)
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
