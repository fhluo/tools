package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"slices"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message/pipeline"
)

var (
	config = &pipeline.Config{
		TranslationsPattern: `messages\.(.*)\.json$`,
		OutPattern:          "",
		Format:              "",
		Ext:                 "",
		DeclareVar:          "",
		SetDefault:          false,
	}
	languages      []string
	sourceLanguage string
)

var rootCmd = &cobra.Command{
	Use: "gotext [package]...",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		config.Supported = slices.Collect(func(yield func(tag language.Tag) bool) {
			for lang := range slices.Values(languages) {
				if lang == "" {
					continue
				}

				tag, err := language.Parse(lang)
				if err != nil {
					slog.Error("failed to parse", err, "lang", lang)
					continue
				}

				if !yield(tag) {
					return
				}
			}
		})

		config.SourceLanguage, err = language.Parse(sourceLanguage)
		if err != nil {
			return fmt.Errorf("failed to parse source language %q: %w", sourceLanguage, err)
		}

		config.Packages = args

		state, err := pipeline.Extract(config)
		if err != nil {
			return fmt.Errorf("failed to extract: %w", err)
		}

		if err = state.Import(); err != nil {
			return fmt.Errorf("failed to import: %w", err)
		}
		if err = state.Merge(); err != nil {
			return fmt.Errorf("failed to merge: %w", err)
		}
		if err = state.Export(); err != nil {
			return fmt.Errorf("failed to export: %w", err)
		}
		if err = state.Generate(); err != nil {
			return fmt.Errorf("failed to generate: %w", err)
		}

		return nil
	},
}

func init() {
	log.SetFlags(0)

	rootCmd.Flags().StringSliceVarP(&languages, "lang", "l", []string{"en-US"}, "supported languages")
	rootCmd.Flags().StringVarP(&sourceLanguage, "src", "s", "en-US", "source language")
	rootCmd.Flags().StringVarP(&config.Dir, "dir", "d", "locales", "directory to store translation files")
	rootCmd.Flags().StringVarP(&config.GenFile, "out", "o", "catalog.go", "output file name")
	rootCmd.Flags().StringVarP(&config.GenPackage, "pkg", "p", "", "output package name")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to execute root command", "err", err)
		os.Exit(1)
	}
}
