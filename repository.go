package mailfull

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// Errors for the operation of the Repository.
var (
	ErrDomainNotExist = errors.New("Domain: not exist")
	ErrUserNotExist   = errors.New("User: not exist")

	ErrInvalidFormatUsersPassword = errors.New("User: password file invalid format")
	ErrInvalidFormatAliasDomain   = errors.New("AliasDomain: file invalid format")
	ErrInvalidFormatAliasUsers    = errors.New("AliasUsers: file invalid format")
)

// Repository represents a Repository.
type Repository struct {
	*RepositoryConfig
}

// NewRepository creates a new Repository instance.
func NewRepository(c *RepositoryConfig) (*Repository, error) {
	r := &Repository{
		RepositoryConfig: c,
	}

	return r, nil
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

// Users returns a User slice.
func (r *Repository) Users(domainName string) ([]*User, error) {
	domain, err := r.Domain(domainName)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, ErrDomainNotExist
	}

	hashedPasswords, err := r.usersHashedPassword(domainName)
	if err != nil {
		return nil, err
	}

	fileInfos, err := ioutil.ReadDir(filepath.Join(r.DirMailDataPath, domainName))
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(fileInfos))

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			continue
		}

		name := fileInfo.Name()

		forwards, err := r.userForwards(domainName, name)
		if err != nil {
			return nil, err
		}

		hashedPassword, ok := hashedPasswords[name]
		if !ok {
			hashedPassword = ""
		}

		user, err := NewUser(name, hashedPassword, forwards)
		if err != nil {
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

// User returns a User of the input name.
func (r *Repository) User(domainName, userName string) (*User, error) {
	domain, err := r.Domain(domainName)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, ErrDomainNotExist
	}

	if !validUserName(userName) {
		return nil, ErrInvalidUserName
	}

	hashedPasswords, err := r.usersHashedPassword(domainName)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(filepath.Join(r.DirMailDataPath, domainName, userName))
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, nil
	}

	name := userName

	forwards, err := r.userForwards(domainName, name)
	if err != nil {
		return nil, err
	}

	hashedPassword, ok := hashedPasswords[name]
	if !ok {
		hashedPassword = ""
	}

	user, err := NewUser(name, hashedPassword, forwards)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// usersHashedPassword returns a string map of usernames to the hashed password.
func (r *Repository) usersHashedPassword(domainName string) (map[string]string, error) {
	domain, err := r.Domain(domainName)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, ErrDomainNotExist
	}

	file, err := os.Open(filepath.Join(r.DirMailDataPath, domainName, FileNameUsersPassword))
	if err != nil {
		return nil, err
	}

	hashedPasswords := map[string]string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), ":")
		if len(words) != 2 {
			return nil, ErrInvalidFormatUsersPassword
		}

		name := words[0]
		hashedPassword := words[1]

		hashedPasswords[name] = hashedPassword
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hashedPasswords, nil
}

// userForwards returns a string slice of forwards that the input name has.
func (r *Repository) userForwards(domainName, userName string) ([]string, error) {
	user, err := r.User(domainName, userName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotExist
	}

	file, err := os.Open(filepath.Join(r.DirMailDataPath, domainName, userName, FileNameUserForwards))
	if err != nil {
		if err.(*os.PathError).Err == syscall.ENOENT {
			return []string{}, nil
		}

		return nil, err
	}

	forwards := make([]string, 0, 5)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		forwards = append(forwards, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return forwards, nil
}

// AliasUsers returns a AliasUser slice.
func (r *Repository) AliasUsers(domainName string) ([]*AliasUser, error) {
	domain, err := r.Domain(domainName)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, ErrDomainNotExist
	}

	file, err := os.Open(filepath.Join(r.DirMailDataPath, domainName, FileNameAliasUsers))
	if err != nil {
		return nil, err
	}

	aliasUsers := make([]*AliasUser, 0, 50)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), ":")
		if len(words) != 2 {
			return nil, ErrInvalidFormatAliasUsers
		}

		name := words[0]
		targets := strings.Split(words[1], ",")

		aliasUser, err := NewAliasUser(name, targets)
		if err != nil {
			return nil, err
		}

		aliasUsers = append(aliasUsers, aliasUser)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return aliasUsers, nil
}

// AliasUser returns a AliasUser of the input name.
func (r *Repository) AliasUser(domainName, aliasUserName string) (*AliasUser, error) {
	aliasUsers, err := r.AliasUsers(domainName)
	if err != nil {
		return nil, err
	}

	for _, aliasUser := range aliasUsers {
		if aliasUser.Name() == aliasUserName {
			return aliasUser, nil
		}
	}

	return nil, nil
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
