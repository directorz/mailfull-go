package mailfull

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// User represents a User.
type User struct {
	name           string
	hashedPassword string
	forwards       []string
}

// UserSlice attaches the methods of sort.Interface to []*User.
type UserSlice []*User

func (p UserSlice) Len() int           { return len(p) }
func (p UserSlice) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
func (p UserSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// NewUser creates a new User instance.
func NewUser(name, hashedPassword string, forwards []string) (*User, error) {
	u := &User{}

	if err := u.setName(name); err != nil {
		return nil, err
	}

	u.SetHashedPassword(hashedPassword)
	u.SetForwards(forwards)

	return u, nil
}

// setName sets the name.
func (u *User) setName(name string) error {
	if !validUserName(name) {
		return ErrInvalidUserName
	}

	u.name = name

	return nil
}

// Name returns name.
func (u *User) Name() string {
	return u.name
}

// SetHashedPassword sets the hashed password.
func (u *User) SetHashedPassword(hashedPassword string) {
	u.hashedPassword = hashedPassword
}

// HashedPassword returns hashedPassword.
func (u *User) HashedPassword() string {
	return u.hashedPassword
}

// SetForwards sets forwards.
func (u *User) SetForwards(forwards []string) {
	u.forwards = forwards
}

// Forwards returns forwards.
func (u *User) Forwards() []string {
	return u.forwards
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
		if err.(*os.PathError).Err == syscall.ENOENT {
			return nil, nil
		}

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
	defer file.Close()

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

	file, err := os.Open(filepath.Join(r.DirMailDataPath, domainName, userName, FileNameUserForwards))
	if err != nil {
		if err.(*os.PathError).Err == syscall.ENOENT {
			return nil, nil
		}

		return nil, err
	}
	defer file.Close()

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
