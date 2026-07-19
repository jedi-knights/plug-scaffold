package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jedi-knights/plug-scaffold/internal/adapters/emitter"
	"github.com/jedi-knights/plug-scaffold/internal/adapters/renderer"
	"github.com/jedi-knights/plug-scaffold/internal/domain"
)

// NewNewCmd returns the `plug-scaffold new` command.
func NewNewCmd() *cobra.Command {
	var (
		author string
		org    string
		style  string
		dir    string
	)

	cmd := &cobra.Command{
		Use:   "new <plugin-name>",
		Short: "Scaffold a new Neovim plugin project",
		Long: `Emits a fresh Neovim plugin skeleton (LICENSE, README, .gitignore, plus
style-specific Lua files) into <dir>/<plugin-name>. The default style is
"omar"; --style=tj|prime picks a different convention.

The author defaults to the current git user.name; --author overrides it.
--org is required (no reasonable default — the install snippet needs it).`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runNew(args[0], author, org, style, dir)
		},
	}

	cmd.Flags().StringVar(&author, "author", "", "Author name (defaults to `git config user.name`)")
	cmd.Flags().StringVar(&org, "org", "", "GitHub org or user that will own the repo (required)")
	cmd.Flags().StringVar(&style, "style", string(domain.StyleOmar), "Template style: tj|prime|omar")
	cmd.Flags().StringVar(&dir, "dir", "", "Parent directory for the new plugin (default: current directory)")

	return cmd
}

func runNew(pluginName, author, org, styleName, dir string) error {
	if author == "" {
		author = gitConfigUserName()
	}
	if author == "" {
		return fmt.Errorf("--author is required (git config user.name is not set)")
	}

	style, err := domain.ParseStyle(styleName)
	if err != nil {
		return err
	}

	spec, err := domain.NewPluginSpec(pluginName, author, org, style)
	if err != nil {
		return err
	}

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("resolving current directory: %w", err)
		}
	}
	target := filepath.Join(dir, spec.RepoName)

	// Fail loudly if target exists — silently overwriting a directory
	// is the class of bug this rule most guards against.
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("refusing to overwrite existing path: %s", target)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking target %s: %w", target, err)
	}

	out := emitter.New(target)
	if err := renderer.NewBase().Render(spec, out); err != nil {
		return err
	}

	fmt.Printf("Scaffolded %s at %s (style: %s)\n", spec.RepoName, target, spec.Style)
	return nil
}

// gitConfigUserName reads git config user.name via the git binary. An
// empty return means git is unavailable, the config is unset, or the
// output was blank — all indistinguishable to the caller and treated
// the same.
func gitConfigUserName() string {
	out, err := exec.Command("git", "config", "--get", "user.name").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
