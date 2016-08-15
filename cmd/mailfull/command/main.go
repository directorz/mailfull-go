package command

import (
	"github.com/mitchellh/cli"
)

// Meta is for `*Command` struct.
type Meta struct {
	UI         *cli.BasicUi
	CmdName    string
	SubCmdName string
	Version    string
}
