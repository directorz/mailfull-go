package mailfull

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

// Domain represents a Domain.
type Domain struct {
	name         string
	Users        []*User
	AliasUsers   []*AliasUser
	CatchAllUser *CatchAllUser
}

// DomainSlice attaches the methods of sort.Interface to []*Domain.
type DomainSlice []*Domain

func (p DomainSlice) Len() int           { return len(p) }
func (p DomainSlice) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
func (p DomainSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// NewDomain creates a new Domain instance.
func NewDomain(name string) (*Domain, error) {
	if !validDomainName(name) {
		return nil, ErrInvalidDomainName
	}

	d := &Domain{
		name: name,
	}

	return d, nil
}

// Name returns name.
func (d *Domain) Name() string {
	return d.name
}

// Domains returns a Domain slice.
func (r *Repository) Domains() ([]*Domain, error) {
	fileInfos, err := ioutil.ReadDir(r.DirMailDataPath)
	if err != nil {
		return nil, err
	}

	domains := make([]*Domain, 0, len(fileInfos))

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			continue
		}

		name := fileInfo.Name()

		domain, err := NewDomain(name)
		if err != nil {
			continue
		}

		domains = append(domains, domain)
	}

	return domains, nil
}

// Domain returns a Domain of the input name.
func (r *Repository) Domain(domainName string) (*Domain, error) {
	if !validDomainName(domainName) {
		return nil, ErrInvalidDomainName
	}

	fileInfo, err := os.Stat(filepath.Join(r.DirMailDataPath, domainName))
	if err != nil {
		if err.(*os.PathError).Err == syscall.ENOENT {
			return nil, nil
		}

		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, nil
	}

	name := domainName

	domain, err := NewDomain(name)
	if err != nil {
		return nil, err
	}

	return domain, nil
}
