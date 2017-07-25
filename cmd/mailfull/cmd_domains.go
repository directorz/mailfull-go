package main

import (
	"fmt"
	"sort"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdDomains represents a CmdDomains.
type CmdDomains struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdDomains) Synopsis() string {
	return "Show domains."
}

// Help returns long-form help text.
func (c *CmdDomains) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s

Description:
    %s
    Disabled domains are marked "!" the beginning.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdDomains) Run(args []string) int {
	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	domains, err := repo.Domains()
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}
	sort.Sort(mailfull.DomainSlice(domains))

	for _, domain := range domains {
		disableStr := ""
		if domain.Disabled() {
			disableStr = "!"
		}

		fmt.Fprintf(c.UI.Writer, "%s%s\n", disableStr, domain.Name())
	}

	return 0
}
