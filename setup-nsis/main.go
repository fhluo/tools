package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var (
	source = "https://onboardcloud.dl.sourceforge.net/project/nsis/NSIS%203/3.09/nsis-3.09.zip"
	folder = "nsis-3.09"
)

var rootCmd = &cobra.Command{
	Use:   "setup-nsis dir",
	Short: "setup nsis",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dst := args[0]

		err := getter.GetAny(dst, source, []getter.ClientOption{}...)
		if err != nil {
			return err
		}

		err = os.Rename(filepath.Join(dst, folder), filepath.Join(dst, "nsis"))

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
