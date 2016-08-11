package mailfull

import (
	"bufio"
	"fmt"
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

// CatchAllUserSet sets a CatchAllUser to the input Domain.
func (r *Repository) CatchAllUserSet(domainName string, catchAllUser *CatchAllUser) error {
	existUser, err := r.User(domainName, catchAllUser.Name())
	if err != nil {
		return err
	}
	if existUser == nil {
		return ErrUserNotExist
	}

	file, err := os.OpenFile(filepath.Join(r.DirMailDataPath, domainName, FileNameCatchAllUser), os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "%s\n", catchAllUser.Name()); err != nil {
		return err
	}

	return nil
}

// CatchAllUserUnset removes a CatchAllUser from the input Domain.
func (r *Repository) CatchAllUserUnset(domainName string) error {
	existDomain, err := r.Domain(domainName)
	if err != nil {
		return err
	}
	if existDomain == nil {
		return ErrDomainNotExist
	}

	file, err := os.OpenFile(filepath.Join(r.DirMailDataPath, domainName, FileNameCatchAllUser), os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}
