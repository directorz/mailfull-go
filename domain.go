package mailfull

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// Domain represents a Domain.
type Domain struct {
	name         string
	disabled     bool
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
	d := &Domain{}

	if err := d.setName(name); err != nil {
		return nil, err
	}

	return d, nil
}

// setName sets the name.
func (d *Domain) setName(name string) error {
	if !validDomainName(name) {
		return ErrInvalidDomainName
	}

	d.name = name

	return nil
}

// Name returns name.
func (d *Domain) Name() string {
	return d.name
}

// SetDisabled disables the Domain if the input is true.
func (d *Domain) SetDisabled(disabled bool) {
	d.disabled = disabled
}

// Disabled returns true if the Domain is disabled.
func (d *Domain) Disabled() bool {
	return d.disabled
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

		disabled, err := r.domainDisabled(name)
		if err != nil {
			return nil, err
		}
		domain.SetDisabled(disabled)

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

	disabled, err := r.domainDisabled(name)
	if err != nil {
		return nil, err
	}
	domain.SetDisabled(disabled)

	return domain, nil
}

// domainDisabled returns true if the input Domain is disabled.
func (r *Repository) domainDisabled(domainName string) (bool, error) {
	if !validDomainName(domainName) {
		return false, ErrInvalidDomainName
	}

	fi, err := os.Stat(filepath.Join(r.DirMailDataPath, domainName, FileNameDomainDisable))

	if err != nil {
		if err.(*os.PathError).Err == syscall.ENOENT {
			return false, nil
		}

		return false, err
	}

	if fi.IsDir() {
		return false, ErrInvalidFormatDomainDisabled
	}

	return true, nil
}

// DomainCreate creates the input Domain.
func (r *Repository) DomainCreate(domain *Domain) error {
	existDomain, err := r.Domain(domain.Name())
	if err != nil {
		return err
	}
	if existDomain != nil {
		return ErrDomainAlreadyExist
	}
	existAliasDomain, err := r.AliasDomain(domain.Name())
	if err != nil {
		return err
	}
	if existAliasDomain != nil {
		return ErrAliasDomainAlreadyExist
	}

	domainDirPath := filepath.Join(r.DirMailDataPath, domain.Name())

	if err := os.Mkdir(domainDirPath, 0700); err != nil {
		return err
	}
	if err := os.Chown(domainDirPath, r.uid, r.gid); err != nil {
		return err
	}

	usersPasswordFile, err := os.OpenFile(filepath.Join(domainDirPath, FileNameUsersPassword), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if err := usersPasswordFile.Chown(r.uid, r.gid); err != nil {
		return err
	}
	usersPasswordFile.Close()

	aliasUsersFile, err := os.OpenFile(filepath.Join(domainDirPath, FileNameAliasUsers), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if err := aliasUsersFile.Chown(r.uid, r.gid); err != nil {
		return err
	}
	aliasUsersFile.Close()

	catchAllUserFile, err := os.OpenFile(filepath.Join(domainDirPath, FileNameCatchAllUser), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if err := catchAllUserFile.Chown(r.uid, r.gid); err != nil {
		return err
	}
	catchAllUserFile.Close()

	if domain.Disabled() {
		if err := r.writeDomainDisabledFile(domain.Name(), domain.Disabled()); err != nil {
			return err
		}
	}

	return nil
}

// DomainUpdate updates the input Domain.
func (r *Repository) DomainUpdate(domain *Domain) error {
	existDomain, err := r.Domain(domain.Name())
	if err != nil {
		return err
	}
	if existDomain == nil {
		return ErrDomainNotExist
	}

	if err := r.writeDomainDisabledFile(domain.Name(), domain.Disabled()); err != nil {
		return err
	}

	return nil
}

// DomainRemove removes a Domain of the input name.
func (r *Repository) DomainRemove(domainName string) error {
	existDomain, err := r.Domain(domainName)
	if err != nil {
		return err
	}
	if existDomain == nil {
		return ErrDomainNotExist
	}

	aliasDomains, err := r.AliasDomains()
	if err != nil {
		return err
	}
	for _, aliasDomain := range aliasDomains {
		if aliasDomain.Target() == domainName {
			return ErrDomainIsAliasDomainTarget
		}
	}

	domainDirPath := filepath.Join(r.DirMailDataPath, domainName)
	domainBackupDirPath := filepath.Join(r.DirMailDataPath, "."+domainName+".deleted."+time.Now().Format("20060102150405"))

	if err := os.Rename(domainDirPath, domainBackupDirPath); err != nil {
		return err
	}

	return nil
}

// writeDomainDisabledFile creates/removes the disabled file.
func (r *Repository) writeDomainDisabledFile(domainName string, disabled bool) error {
	if !validDomainName(domainName) {
		return ErrInvalidDomainName
	}

	nowDisabled, err := r.domainDisabled(domainName)
	if err != nil {
		return err
	}

	domainDisabledFileName := filepath.Join(r.DirMailDataPath, domainName, FileNameDomainDisable)

	if !nowDisabled && disabled {
		file, err := os.OpenFile(domainDisabledFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		if err := file.Chown(r.uid, r.gid); err != nil {
			return err
		}
		file.Close()
	}

	if nowDisabled && !disabled {
		if err := os.Remove(domainDisabledFileName); err != nil {
			return err
		}
	}

	return nil
}
