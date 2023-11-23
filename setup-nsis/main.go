package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"log/slog"

	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
)

const version = "3.09"

var (
	source = fmt.Sprintf("https://onboardcloud.dl.sourceforge.net/project/nsis/NSIS%%203/%[1]s/nsis-%[1]s.zip", version)
	folder = fmt.Sprintf("nsis-%s", version)
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

		return os.Rename(filepath.Join(dst, folder), filepath.Join(dst, "nsis"))
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
