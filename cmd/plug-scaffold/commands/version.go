package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd returns the `plug-scaffold version` command.
func NewVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print plug-scaffold version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("plug-scaffold", version)
		},
	}
}
