package main

import (
	"fmt"
	"strings"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
	"github.com/jsimonetti/pwscheme/ssha"
)

// CmdUserCheckPw represents a CmdUserCheckPw.
type CmdUserCheckPw struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdUserCheckPw) Synopsis() string {
	return "Check user's password."
}

// Help returns long-form help text.
func (c *CmdUserCheckPw) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s address [password]

Description:
    %s

Required Args:
    address
        The email address that you want to check the password.

Optional Args:
    password
        Specify the password instead of your typing.
        This option is NOT recommended because the password will be visible in your shell history.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdUserCheckPw) Run(args []string) int {
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
		input, err := c.UI.AskSecret(fmt.Sprintf("Enter password for %s:", address))
		if err != nil {
			fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
			return 1
		}

		rawPassword = input
	}

	if ok, _ := ssha.Validate(rawPassword, user.HashedPassword()); !ok {
		fmt.Fprintf(c.UI.Writer, "The password you entered is incorrect.\n")
		return 1
	}

	fmt.Fprintf(c.UI.Writer, "The password you entered is correct.\n")
	return 0
}
