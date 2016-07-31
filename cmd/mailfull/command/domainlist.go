package command

import (
	"fmt"
	"sort"

	"github.com/directorz/mailfull-go"
)

// DomainListCommand represents a DomainListCommand.
type DomainListCommand struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *DomainListCommand) Synopsis() string {
	return "Show all domains."
}

// Help returns long-form help text.
func (c *DomainListCommand) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s

Description:
    Show all domains.
`,
		c.CmdName, c.SubCmdName)

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *DomainListCommand) Run(args []string) int {
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
