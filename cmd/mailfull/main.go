package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd/mailfull/command"
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

	meta := command.Meta{
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
			return &command.InitCommand{Meta: meta}, nil
		},
		"genconfig": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.GenConfigCommand{Meta: meta}, nil
		},
		"domains": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.DomainsCommand{Meta: meta}, nil
		},
		"domainadd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.DomainAddCommand{Meta: meta}, nil
		},
		"domaindel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.DomainDelCommand{Meta: meta}, nil
		},
		"aliasdomains": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.AliasDomainsCommand{Meta: meta}, nil
		},
		"aliasdomainadd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.AliasDomainAddCommand{Meta: meta}, nil
		},
		"aliasdomaindel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.AliasDomainDelCommand{Meta: meta}, nil
		},
		"users": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.UsersCommand{Meta: meta}, nil
		},
		"useradd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.UserAddCommand{Meta: meta}, nil
		},
		"userdel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.UserDelCommand{Meta: meta}, nil
		},
		"userpasswd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.UserPasswdCommand{Meta: meta}, nil
		},
		"usercheckpw": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.UserCheckPwCommand{Meta: meta}, nil
		},
		"aliasusers": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.AliasUsersCommand{Meta: meta}, nil
		},
		"aliasuseradd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.AliasUserAddCommand{Meta: meta}, nil
		},
		"aliasusermod": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.AliasUserModCommand{Meta: meta}, nil
		},
		"aliasuserdel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.AliasUserDelCommand{Meta: meta}, nil
		},
		"catchall": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.CatchAllCommand{Meta: meta}, nil
		},
		"catchallset": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.CatchAllSetCommand{Meta: meta}, nil
		},
		"catchallunset": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.CatchAllUnsetCommand{Meta: meta}, nil
		},
		"commit": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &command.CommitCommand{Meta: meta}, nil
		},
	}

	exitCode, err := c.Run()
	if err != nil {
		fmt.Fprintf(meta.UI.ErrorWriter, "%v\n", err)
	}

	os.Exit(exitCode)
}
