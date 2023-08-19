package main

import (
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "wails-build",
	Short: "wails build",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		// 向 wails 传递的参数
		args := []string{"build"}

		if useUPX {
			args = append(args, "-upx")
		}
		if useNSIS {
			args = append(args, "-nsis")
		}

		// 将路径添加到环境变量中
		for i := range path {
			path[i], err = filepath.Abs(path[i])
			if err != nil {
				return err
			}
		}

		err = os.Setenv("path", strings.Join(path, ";")+";"+os.Getenv("path"))
		if err != nil {
			return err
		}

		// 执行 wails build 命令
		command := exec.Command("wails", args...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		slog.Info(command.String(), "path", os.Getenv("path"))

		return command.Run()
	},
}

var (
	useUPX  bool     // 是否使用 UPX
	useNSIS bool     // 是否使用 NSIS
	path    []string // 要添加到环境变量中的路径
)

func init() {
	log.SetFlags(0)

	rootCmd.Flags().BoolVar(&useUPX, "upx", false, "use UPX")
	rootCmd.Flags().BoolVar(&useNSIS, "nsis", false, "use NSIS")
	rootCmd.Flags().StringSliceVar(&path, "path", nil, "paths to add to the environment variable")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
