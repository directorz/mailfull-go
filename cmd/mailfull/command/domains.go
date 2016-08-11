package command

import (
	"fmt"
	"sort"

	"github.com/directorz/mailfull-go"
)

// DomainsCommand represents a DomainsCommand.
type DomainsCommand struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *DomainsCommand) Synopsis() string {
	return "Show domains."
}

// Help returns long-form help text.
func (c *DomainsCommand) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s

Description:
    %s
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *DomainsCommand) Run(args []string) int {
	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	domains, err := repo.Domains()
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}
	sort.Sort(mailfull.DomainSlice(domains))

	for _, domain := range domains {
		fmt.Fprintf(c.UI.Writer, "%s\n", domain.Name())
	}

	return 0
}
