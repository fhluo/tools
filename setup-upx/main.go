package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
)

const version = "4.2.1"

var (
	source = fmt.Sprintf("https://github.com/upx/upx/releases/download/v%[1]s/upx-%[1]s-win64.zip", version)
	folder = fmt.Sprintf("upx-%s-win64", version)
)

var rootCmd = &cobra.Command{
	Use:   "setup-upx dir",
	Short: "setup upx",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dst := args[0]

		err := getter.GetAny(dst, source, []getter.ClientOption{}...)
		if err != nil {
			return err
		}

		return os.Rename(filepath.Join(dst, folder), filepath.Join(dst, "upx"))
	},
}

func init() {
	log.SetFlags(0)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
