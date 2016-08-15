package mailfull

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
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
	if !validDomainName(domainName) {
		return nil, ErrInvalidDomainName
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
	if !validDomainName(domainName) {
		return nil, ErrInvalidDomainName
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

// UserCreate creates the input User.
func (r *Repository) UserCreate(domainName string, user *User) error {
	existUser, err := r.User(domainName, user.Name())
	if err != nil {
		return err
	}
	if existUser != nil {
		return ErrUserAlreadyExist
	}
	existAliasUser, err := r.AliasUser(domainName, user.Name())
	if err != nil {
		return err
	}
	if existAliasUser != nil {
		return ErrAliasUserAlreadyExist
	}

	userDirPath := filepath.Join(r.DirMailDataPath, domainName, user.Name())

	dirNames := []string{
		userDirPath,
		filepath.Join(userDirPath, "Maildir"),
		filepath.Join(userDirPath, "Maildir/cur"),
		filepath.Join(userDirPath, "Maildir/new"),
		filepath.Join(userDirPath, "Maildir/tmp"),
	}
	for _, dirName := range dirNames {
		if err := os.Mkdir(dirName, 0700); err != nil {
			return err
		}
		if err := os.Chown(dirName, r.uid, r.gid); err != nil {
			return err
		}
	}

	if err := r.UserUpdate(domainName, user); err != nil {
		return err
	}

	return nil
}

// UserUpdate updates the input User.
func (r *Repository) UserUpdate(domainName string, user *User) error {
	existUser, err := r.User(domainName, user.Name())
	if err != nil {
		return err
	}
	if existUser == nil {
		return ErrUserNotExist
	}

	hashedPasswords, err := r.usersHashedPassword(domainName)
	if err != nil {
		return err
	}
	hashedPasswords[user.Name()] = user.HashedPassword()
	if err := r.writeUsersPasswordFile(domainName, hashedPasswords); err != nil {
		return err
	}

	if err := r.writeUserForwardsFile(domainName, user.Name(), user.Forwards()); err != nil {
		return err
	}

	return nil
}

// UserRemove removes a User of the input name.
func (r *Repository) UserRemove(domainName, userName string) error {
	existUser, err := r.User(domainName, userName)
	if err != nil {
		return err
	}
	if existUser == nil {
		return ErrUserNotExist
	}

	catchAllUser, err := r.CatchAllUser(domainName)
	if err != nil {
		return err
	}
	if catchAllUser != nil && catchAllUser.Name() == userName {
		return ErrUserIsCatchAllUser
	}

	hashedPasswords, err := r.usersHashedPassword(domainName)
	if err != nil {
		return err
	}
	delete(hashedPasswords, userName)
	if err := r.writeUsersPasswordFile(domainName, hashedPasswords); err != nil {
		return err
	}

	userDirPath := filepath.Join(r.DirMailDataPath, domainName, userName)
	userBackupDirPath := filepath.Join(r.DirMailDataPath, domainName, "."+userName+".deleted."+time.Now().Format("20060102150405"))

	if err := os.Rename(userDirPath, userBackupDirPath); err != nil {
		return err
	}

	return nil
}

// writeUsersPasswordFile writes passwords of each users to the file.
func (r *Repository) writeUsersPasswordFile(domainName string, hashedPasswords map[string]string) error {
	if !validDomainName(domainName) {
		return ErrInvalidDomainName
	}

	keys := make([]string, 0, len(hashedPasswords))
	for key := range hashedPasswords {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	file, err := os.OpenFile(filepath.Join(r.DirMailDataPath, domainName, FileNameUsersPassword), os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, key := range keys {
		if _, err := fmt.Fprintf(file, "%s:%s\n", key, hashedPasswords[key]); err != nil {
			return err
		}
	}

	return nil
}

// writeUserForwardsFile writes forwards to user's forward file.
func (r *Repository) writeUserForwardsFile(domainName, userName string, forwards []string) error {
	if !validDomainName(domainName) {
		return ErrInvalidDomainName
	}
	if !validUserName(userName) {
		return ErrInvalidUserName
	}

	file, err := os.OpenFile(filepath.Join(r.DirMailDataPath, domainName, userName, FileNameUserForwards), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if err := file.Chown(r.uid, r.gid); err != nil {
		return err
	}
	defer file.Close()

	for _, forward := range forwards {
		if _, err := fmt.Fprintf(file, "%s\n", forward); err != nil {
			return err
		}
	}

	return nil
}
