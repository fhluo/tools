package main

import (
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"log"
	"os"
	"path/filepath"
)

var (
	source = "https://github.com/upx/upx/releases/download/v4.0.2/upx-4.0.2-win64.zip"
	folder = "upx-4.0.2-win64"
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

		err = os.Rename(filepath.Join(dst, folder), filepath.Join(dst, "upx"))

		return nil
	},
}

func init() {
	log.SetFlags(0)

	rootCmd.Flags().StringVar(&source, "src", source, "source")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
