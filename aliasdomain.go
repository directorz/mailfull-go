package mailfull

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// AliasDomain represents a AliasDomain.
type AliasDomain struct {
	name   string
	target string
}

// AliasDomainSlice attaches the methods of sort.Interface to []*AliasDomain.
type AliasDomainSlice []*AliasDomain

func (p AliasDomainSlice) Len() int           { return len(p) }
func (p AliasDomainSlice) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
func (p AliasDomainSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// NewAliasDomain creates a new AliasDomain instance.
func NewAliasDomain(name, target string) (*AliasDomain, error) {
	if !validAliasDomainName(name) {
		return nil, ErrInvalidAliasDomainName
	}
	if !validAliasDomainTarget(target) {
		return nil, ErrInvalidAliasDomainTarget
	}

	ad := &AliasDomain{
		name:   name,
		target: target,
	}

	return ad, nil
}

// Name returns name.
func (ad *AliasDomain) Name() string {
	return ad.name
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
