package mailfull

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// AliasDomain represents a AliasDomain.
type AliasDomain struct {
	name   string
	target string
}

// NewAliasDomain creates a new AliasDomain instance.
func NewAliasDomain(name, target string) (*AliasDomain, error) {
	ad := &AliasDomain{}

	if err := ad.setName(name); err != nil {
		return nil, err
	}

	if err := ad.SetTarget(target); err != nil {
		return nil, err
	}

	return ad, nil
}

// setName sets the name.
func (ad *AliasDomain) setName(name string) error {
	if !validAliasDomainName(name) {
		return ErrInvalidAliasDomainName
	}

	ad.name = name

	return nil
}

// Name returns name.
func (ad *AliasDomain) Name() string {
	return ad.name
}

// SetTarget sets the target.
func (ad *AliasDomain) SetTarget(target string) error {
	if !validAliasDomainTarget(target) {
		return ErrInvalidAliasDomainTarget
	}

	ad.target = target

	return nil
}

// Target returns target.
func (ad *AliasDomain) Target() string {
	return ad.target
}

// AliasDomains returns a AliasDomain slice.
func (r *Repository) AliasDomains() ([]*AliasDomain, error) {
	file, err := os.Open(filepath.Join(r.DirMailDataPath, FileNameAliasDomains))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	aliasDomains := make([]*AliasDomain, 0, 10)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), ":")
		if len(words) != 2 {
			return nil, ErrInvalidFormatAliasDomain
		}

		name := words[0]
		target := words[1]

		aliasDomain, err := NewAliasDomain(name, target)
		if err != nil {
			return nil, err
		}

		aliasDomains = append(aliasDomains, aliasDomain)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return aliasDomains, nil
}

// AliasDomain returns a AliasDomain of the input name.
func (r *Repository) AliasDomain(aliasDomainName string) (*AliasDomain, error) {
	aliasDomains, err := r.AliasDomains()
	if err != nil {
		return nil, err
	}

	for _, aliasDomain := range aliasDomains {
		if aliasDomain.Name() == aliasDomainName {
			return aliasDomain, nil
		}
	}

	return nil, nil
}

// AliasDomainCreate creates the input AliasDomain.
func (r *Repository) AliasDomainCreate(aliasDomain *AliasDomain) error {
	aliasDomains, err := r.AliasDomains()
	if err != nil {
		return err
	}

	for _, ad := range aliasDomains {
		if ad.Name() == aliasDomain.Name() {
			return ErrAliasDomainAlreadyExist
		}
	}
	existDomain, err := r.Domain(aliasDomain.Name())
	if err != nil {
		return err
	}
	if existDomain != nil {
		return ErrDomainAlreadyExist
	}
	existDomain, err = r.Domain(aliasDomain.Target())
	if err != nil {
		return err
	}
	if existDomain == nil {
		return ErrDomainNotExist
	}

	aliasDomains = append(aliasDomains, aliasDomain)

	if err := r.writeAliasDomainsFile(aliasDomains); err != nil {
		return err
	}

	return nil
}

// AliasDomainRemove removes a AliasDomain of the input name.
func (r *Repository) AliasDomainRemove(aliasDomainName string) error {
	aliasDomains, err := r.AliasDomains()
	if err != nil {
		return err
	}

	idx := -1
	for i, aliasDomain := range aliasDomains {
		if aliasDomain.Name() == aliasDomainName {
			idx = i
		}
	}
	if idx < 0 {
		return ErrAliasDomainNotExist
	}

	aliasDomains = append(aliasDomains[:idx], aliasDomains[idx+1:]...)

	if err := r.writeAliasDomainsFile(aliasDomains); err != nil {
		return err
	}

	return nil
}

// writeAliasDomainsFile writes a AliasDomain slice to the file.
func (r *Repository) writeAliasDomainsFile(aliasDomains []*AliasDomain) error {
	file, err := os.OpenFile(filepath.Join(r.DirMailDataPath, FileNameAliasDomains), os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	sort.Slice(aliasDomains, func(i, j int) bool { return aliasDomains[i].Name() < aliasDomains[j].Name() })

	for _, aliasDomain := range aliasDomains {
		if _, err := fmt.Fprintf(file, "%s:%s\n", aliasDomain.Name(), aliasDomain.Target()); err != nil {
			return err
		}
	}

	return nil
}
