package command

import (
	"fmt"

	"github.com/directorz/mailfull-go"
)

// InitCommand represents a InitCommand.
type InitCommand struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *InitCommand) Synopsis() string {
	return "Initializes current directory as a Mailfull repository."
}

// Help returns long-form help text.
func (c *InitCommand) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s

Description:
    Initializes current directory as a Mailfull repository.
`,
		c.CmdName, c.SubCmdName)

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *InitCommand) Run(args []string) int {
	if err := mailfull.InitRepository("."); err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	return 0
}
