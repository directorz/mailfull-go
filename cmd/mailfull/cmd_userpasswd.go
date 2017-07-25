package main

import (
	"fmt"
	"strings"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
	"github.com/jsimonetti/pwscheme/ssha"
)

// CmdUserPasswd represents a CmdUserPasswd.
type CmdUserPasswd struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdUserPasswd) Synopsis() string {
	return "Update user's password."
}

// Help returns long-form help text.
func (c *CmdUserPasswd) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [-n] address [password]

Description:
    %s

Required Args:
    address
        The email address that you want to update the password.

Optional Args:
    -n
        Don't update databases.
    password
        Specify the password instead of your typing.
        This option is NOT recommended because the password will be visible in your shell history.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdUserPasswd) Run(args []string) int {
	noCommit, err := noCommitFlag(&args)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	if len(args) != 1 && len(args) != 2 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	address := args[0]
	words := strings.Split(address, "@")
	if len(words) != 2 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	userName := words[0]
	domainName := words[1]

	rawPassword := ""
	if len(args) == 2 {
		rawPassword = args[1]
	}

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	user, err := repo.User(domainName, userName)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}
	if user == nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", mailfull.ErrUserNotExist)
		return 1
	}

	if len(args) != 2 {
		input1, err := c.UI.AskSecret(fmt.Sprintf("Enter new password for %s:", address))
		if err != nil {
			fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
			return 1
		}
		input2, err := c.UI.AskSecret("Retype new password:")
		if err != nil {
			fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
			return 1
		}
		if input1 != input2 {
			fmt.Fprintf(c.UI.ErrorWriter, "[ERR] inputs do not match.\n")
			return 1
		}
		rawPassword = input1
	}

	hashedPassword := mailfull.NeverMatchHashedPassword
	if rawPassword != "" {
		str, err := ssha.Generate(rawPassword, 4)
		if err != nil {
			fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
			return 1
		}
		hashedPassword = str
	}

	user.SetHashedPassword(hashedPassword)

	if err := repo.UserUpdate(domainName, user); err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	if noCommit {
		return 0
	}

	mailData, err := repo.MailData()
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	err = repo.GenerateDatabases(mailData)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	return 0
}
