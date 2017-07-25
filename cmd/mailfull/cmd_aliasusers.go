package main

import (
	"fmt"
	"sort"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdAliasUsers represents a CmdAliasUsers.
type CmdAliasUsers struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdAliasUsers) Synopsis() string {
	return "Show aliasusers."
}

// Help returns long-form help text.
func (c *CmdAliasUsers) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s domain

Description:
    %s

Required Args:
    domain
        The domain name.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdAliasUsers) Run(args []string) int {
	if len(args) != 1 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	targetDomainName := args[0]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	aliasUsers, err := repo.AliasUsers(targetDomainName)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}
	sort.Sort(mailfull.AliasUserSlice(aliasUsers))

	for _, aliasUser := range aliasUsers {
		fmt.Fprintf(c.UI.Writer, "%s\n", aliasUser.Name())
	}

	return 0
}
