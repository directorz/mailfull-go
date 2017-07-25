package main

import (
	"fmt"
	"sort"

	mailfull "github.com/directorz/mailfull-go"
)

// AliasDomainsCommand represents a AliasDomainsCommand.
type AliasDomainsCommand struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *AliasDomainsCommand) Synopsis() string {
	return "Show aliasdomains."
}

// Help returns long-form help text.
func (c *AliasDomainsCommand) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [domain]

Description:
    %s

Optional Args:
    domain
        Show aliasdomains that the target is "domain".
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *AliasDomainsCommand) Run(args []string) int {
	if len(args) > 1 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	targetDomainName := ""
	if len(args) == 1 {
		targetDomainName = args[0]
	}

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	aliasDomains, err := repo.AliasDomains()
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}
	sort.Sort(mailfull.AliasDomainSlice(aliasDomains))

	for _, aliasDomain := range aliasDomains {
		if targetDomainName != "" {
			if aliasDomain.Target() == targetDomainName {
				fmt.Fprintf(c.UI.Writer, "%s\n", aliasDomain.Name())
			}
		} else {
			fmt.Fprintf(c.UI.Writer, "%s\n", aliasDomain.Name())
		}
	}

	return 0
}
