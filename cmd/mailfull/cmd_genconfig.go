package main

import (
	"fmt"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdGenConfig represents a CmdGenConfig.
type CmdGenConfig struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdGenConfig) Synopsis() string {
	return "Write a Postfix or Dovecot configuration to stdout."
}

// Help returns long-form help text.
func (c *CmdGenConfig) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s name

Description:
    %s

Required Args:
    name
        The software name that you want to generate a configuration.
        Available names are "postfix" and "dovecot".
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdGenConfig) Run(args []string) int {
	if len(args) != 1 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	softwareName := args[0]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	switch softwareName {
	case "postfix":
		fmt.Fprintf(c.UI.Writer, "%s", repo.GenerateConfigPostfix())

	case "dovecot":
		fmt.Fprintf(c.UI.Writer, "%s", repo.GenerateConfigDovecot())

	default:
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] Specify \"postfix\" or \"dovecot\".\n")
		return 1
	}

	return 0
}
