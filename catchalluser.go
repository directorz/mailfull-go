package mailfull

import (
	"bufio"
	"os"
	"path/filepath"
)

// CatchAllUser represents a CatchAllUser.
type CatchAllUser struct {
	name string
}

// NewCatchAllUser creates a new CatchAllUser instance.
func NewCatchAllUser(name string) (*CatchAllUser, error) {
	if !validCatchAllUserName(name) {
		return nil, ErrInvalidCatchAllUserName
	}

	cu := &CatchAllUser{
		name: name,
	}

	return cu, nil
}

// Name returns name.
func (cu *CatchAllUser) Name() string {
	return cu.name
}

// CatchAllUser returns a CatchAllUser that the input name has.
func (r *Repository) CatchAllUser(domainName string) (*CatchAllUser, error) {
	domain, err := r.Domain(domainName)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, ErrDomainNotExist
	}

	file, err := os.Open(filepath.Join(r.DirMailDataPath, domainName, FileNameCatchAllUser))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	name := scanner.Text()

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if name == "" {
		return nil, nil
	}

	catchAllUser, err := NewCatchAllUser(name)
	if err != nil {
		return nil, err
	}

	return catchAllUser, nil
}
