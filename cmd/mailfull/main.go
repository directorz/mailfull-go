/*
Command mailfull is a CLI application using the mailfull package.
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/directorz/mailfull-go"
	"github.com/mitchellh/cli"
)

var (
	version = mailfull.Version
	gittag  = ""
)

func init() {
	if gittag != "" {
		version = version + "-" + gittag
	}
}

func main() {
	c := &cli.CLI{
		Name:    filepath.Base(os.Args[0]),
		Version: version,
		Args:    os.Args[1:],
	}

	meta := Meta{
		UI: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		CmdName: c.Name,
		Version: c.Version,
	}

	c.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &InitCommand{Meta: meta}, nil
		},
		"genconfig": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &GenConfigCommand{Meta: meta}, nil
		},
		"domains": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &DomainsCommand{Meta: meta}, nil
		},
		"domainadd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &DomainAddCommand{Meta: meta}, nil
		},
		"domaindel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &DomainDelCommand{Meta: meta}, nil
		},
		"domaindisable": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &DomainDisableCommand{Meta: meta}, nil
		},
		"domainenable": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &DomainEnableCommand{Meta: meta}, nil
		},
		"aliasdomains": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &AliasDomainsCommand{Meta: meta}, nil
		},
		"aliasdomainadd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &AliasDomainAddCommand{Meta: meta}, nil
		},
		"aliasdomaindel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &AliasDomainDelCommand{Meta: meta}, nil
		},
		"users": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &UsersCommand{Meta: meta}, nil
		},
		"useradd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &UserAddCommand{Meta: meta}, nil
		},
		"userdel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &UserDelCommand{Meta: meta}, nil
		},
		"userpasswd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &UserPasswdCommand{Meta: meta}, nil
		},
		"usercheckpw": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &UserCheckPwCommand{Meta: meta}, nil
		},
		"aliasusers": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &AliasUsersCommand{Meta: meta}, nil
		},
		"aliasuseradd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &AliasUserAddCommand{Meta: meta}, nil
		},
		"aliasusermod": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &AliasUserModCommand{Meta: meta}, nil
		},
		"aliasuserdel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &AliasUserDelCommand{Meta: meta}, nil
		},
		"catchall": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CatchAllCommand{Meta: meta}, nil
		},
		"catchallset": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CatchAllSetCommand{Meta: meta}, nil
		},
		"catchallunset": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CatchAllUnsetCommand{Meta: meta}, nil
		},
		"commit": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CommitCommand{Meta: meta}, nil
		},
	}

	exitCode, err := c.Run()
	if err != nil {
		fmt.Fprintf(meta.UI.ErrorWriter, "%v\n", err)
	}

	os.Exit(exitCode)
}
