package main

import (
	"fmt"
	"sort"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdUsers represents a CmdUsers.
type CmdUsers struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdUsers) Synopsis() string {
	return "Show users."
}

// Help returns long-form help text.
func (c *CmdUsers) Help() string {
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
func (c *CmdUsers) Run(args []string) int {
	if len(args) != 1 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	targetDomainName := args[0]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	users, err := repo.Users(targetDomainName)
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}
	sort.Slice(users, func(i, j int) bool { return users[i].Name() < users[j].Name() })

	for _, user := range users {
		fmt.Fprintf(c.UI.Writer, "%s\n", user.Name())
	}

	return 0
}
