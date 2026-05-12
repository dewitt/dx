// Command declare is the CLI entry point for the `.dx` toolchain.
//
// It contains no LLM logic (per ARCHITECTURE.md §4): every subcommand is a
// deterministic operation over the `.dx` AST.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dewitt/declare/pkg/export"
	"github.com/dewitt/declare/pkg/lint"
)

// version is overwritten at build time via -ldflags.
var version = "0.1.0-dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		// Cobra has already printed the error; exit non-zero so shells
		// and CI runners pick it up.
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "declare",
		Short:         "Toolchain for the .dx declarative specification language",
		Long:          "declare is a deterministic toolchain for authoring, validating, and exporting .dx files.",
		Version:       version,
		SilenceUsage:  true, // do not dump usage on every command-level error
		SilenceErrors: false,
	}
	root.AddCommand(newLintCmd(), newFmtCmd(), newExportCmd())
	return root
}

func newLintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lint [file...]",
		Short: "Validate one or more .dx files against the SPEC",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var failed bool
			for _, path := range args {
				res, err := lint.LintFile(path)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "%s: %v\n", path, err)
					failed = true
					continue
				}
				for _, issue := range res.Issues {
					fmt.Fprintln(cmd.ErrOrStderr(), issue)
				}
				if !res.OK() {
					failed = true
					continue
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: ok\n", path)
			}
			if failed {
				return fmt.Errorf("lint failed")
			}
			return nil
		},
	}
}

func newFmtCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fmt [file...]",
		Short: "Canonicalize the formatting of .dx files",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintln(cmd.ErrOrStderr(), "fmt: not yet implemented")
			return nil
		},
	}
}

func newExportCmd() *cobra.Command {
	var format string
	c := &cobra.Command{
		Use:   "export [file]",
		Short: "Emit the AST in an agent-optimized format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := lint.LintFile(args[0])
			if err != nil {
				return err
			}
			if !res.OK() {
				for _, issue := range res.Issues {
					fmt.Fprintln(cmd.ErrOrStderr(), issue)
				}
				return fmt.Errorf("export aborted: %s has lint errors", args[0])
			}
			if err := export.Write(cmd.OutOrStdout(), res.Declaration, export.Format(format)); err != nil {
				return err
			}
			return nil
		},
	}
	c.Flags().StringVarP(&format, "format", "f", string(export.FormatJSON), "output format (json)")
	return c
}
