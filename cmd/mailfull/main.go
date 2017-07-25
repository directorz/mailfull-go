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
			return &CmdInit{Meta: meta}, nil
		},
		"genconfig": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdGenConfig{Meta: meta}, nil
		},
		"domains": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdDomains{Meta: meta}, nil
		},
		"domainadd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdDomainAdd{Meta: meta}, nil
		},
		"domaindel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdDomainDel{Meta: meta}, nil
		},
		"domaindisable": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdDomainDisable{Meta: meta}, nil
		},
		"domainenable": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdDomainEnable{Meta: meta}, nil
		},
		"aliasdomains": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdAliasDomains{Meta: meta}, nil
		},
		"aliasdomainadd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdAliasDomainAdd{Meta: meta}, nil
		},
		"aliasdomaindel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdAliasDomainDel{Meta: meta}, nil
		},
		"users": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdUsers{Meta: meta}, nil
		},
		"useradd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdUserAdd{Meta: meta}, nil
		},
		"userdel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdUserDel{Meta: meta}, nil
		},
		"userpasswd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdUserPasswd{Meta: meta}, nil
		},
		"usercheckpw": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdUserCheckPw{Meta: meta}, nil
		},
		"aliasusers": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdAliasUsers{Meta: meta}, nil
		},
		"aliasuseradd": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdAliasUserAdd{Meta: meta}, nil
		},
		"aliasusermod": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdAliasUserMod{Meta: meta}, nil
		},
		"aliasuserdel": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdAliasUserDel{Meta: meta}, nil
		},
		"catchall": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdCatchAll{Meta: meta}, nil
		},
		"catchallset": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdCatchAllSet{Meta: meta}, nil
		},
		"catchallunset": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdCatchAllUnset{Meta: meta}, nil
		},
		"commit": func() (cli.Command, error) {
			meta.SubCmdName = c.Subcommand()
			return &CmdCommit{Meta: meta}, nil
		},
	}

	exitCode, err := c.Run()
	if err != nil {
		fmt.Fprintf(meta.UI.ErrorWriter, "%v\n", err)
	}

	os.Exit(exitCode)
}
