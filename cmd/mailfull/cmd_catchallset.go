package main

import (
	"fmt"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdCatchAllSet represents a CmdCatchAllSet.
type CmdCatchAllSet struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdCatchAllSet) Synopsis() string {
	return "Set a catchall user."
}

// Help returns long-form help text.
func (c *CmdCatchAllSet) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [-n] domain user

Description:
    %s

Required Args:
    domain
        The domain name.
    user
        The user name that you want to set as catchall user.

Optional Args:
    -n
        Don't update databases.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdCatchAllSet) Run(args []string) int {
	noCommit, err := noCommitFlag(&args)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	if len(args) != 2 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	domainName := args[0]
	userName := args[1]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	catchAllUser, err := mailfull.NewCatchAllUser(userName)
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	if err := repo.CatchAllUserSet(domainName, catchAllUser); err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	if noCommit {
		return 0
	}

	mailData, err := repo.MailData()
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	err = repo.GenerateDatabases(mailData)
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	return 0
}
