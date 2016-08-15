package mailfull

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/BurntSushi/toml"
)

// Errors for the Repository.
var (
	ErrInvalidRepository = errors.New("invalid repository")
	ErrNotRepository     = errors.New("not a Mailfull repository (or any of the parent directories)")
	ErrRepositoryExist   = errors.New("a Mailfull repository exists")
)

// Errors for the operation of the Repository.
var (
	ErrDomainNotExist            = errors.New("Domain: not exist")
	ErrDomainAlreadyExist        = errors.New("Domain: already exist")
	ErrDomainIsAliasDomainTarget = errors.New("Domain: is set as alias")

	ErrAliasDomainNotExist     = errors.New("AliasDomain: not exist")
	ErrAliasDomainAlreadyExist = errors.New("AliasDomain: already exist")

	ErrUserNotExist       = errors.New("User: not exist")
	ErrUserAlreadyExist   = errors.New("User: already exist")
	ErrUserIsCatchAllUser = errors.New("User: is set as catchall")

	ErrAliasUserNotExist     = errors.New("AliasUser: not exist")
	ErrAliasUserAlreadyExist = errors.New("AliasUser: already exist")

	ErrInvalidFormatUsersPassword = errors.New("User: password file invalid format")
	ErrInvalidFormatAliasDomain   = errors.New("AliasDomain: file invalid format")
	ErrInvalidFormatAliasUsers    = errors.New("AliasUsers: file invalid format")
)

// RepositoryConfig is used to configure a Repository.
type RepositoryConfig struct {
	DirDatabasePath string `toml:"dir_database"`
	DirMailDataPath string `toml:"dir_maildata"`
	Username        string `toml:"username"`
	CmdPostalias    string `toml:"cmd_postalias"`
	CmdPostmap      string `toml:"cmd_postmap"`
}

// Normalize normalizes paramaters of the RepositoryConfig.
func (c *RepositoryConfig) Normalize(rootPath string) {
	if !filepath.IsAbs(c.DirDatabasePath) {
		c.DirDatabasePath = filepath.Join(rootPath, c.DirDatabasePath)
	}
	if !filepath.IsAbs(c.DirMailDataPath) {
		c.DirMailDataPath = filepath.Join(rootPath, c.DirMailDataPath)
	}

	if filepath.Base(c.CmdPostalias) != c.CmdPostalias {
		if !filepath.IsAbs(c.CmdPostalias) {
			c.CmdPostalias = filepath.Join(rootPath, c.CmdPostalias)
		}
	}

	if filepath.Base(c.CmdPostmap) != c.CmdPostmap {
		if !filepath.IsAbs(c.CmdPostmap) {
			c.CmdPostmap = filepath.Join(rootPath, c.CmdPostmap)
		}
	}
}

// DefaultRepositoryConfig returns a RepositoryConfig with default parameter.
func DefaultRepositoryConfig() *RepositoryConfig {
	c := &RepositoryConfig{
		DirDatabasePath: "./etc",
		DirMailDataPath: "./domains",
		Username:        "",
		CmdPostalias:    "postalias",
		CmdPostmap:      "postmap",
	}

	return c
}

// Repository represents a Repository.
type Repository struct {
	*RepositoryConfig

	uid int
	gid int
}

// NewRepository creates a new Repository instance.
func NewRepository(c *RepositoryConfig) (*Repository, error) {
	u, err := user.Lookup(c.Username)
	if err != nil {
		return nil, err
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return nil, err
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return nil, err
	}

	r := &Repository{
		RepositoryConfig: c,

		uid: uid,
		gid: gid,
	}

	return r, nil
}

// OpenRepository opens a Repository and creates a new Repository instance.
func OpenRepository(basePath string) (*Repository, error) {
	rootPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	for {
		configDirPath := filepath.Join(rootPath, DirNameConfig)

		fi, errStat := os.Stat(configDirPath)
		if errStat != nil {
			if errStat.(*os.PathError).Err != syscall.ENOENT {
				return nil, errStat
			}
		} else {
			if fi.IsDir() {
				break
			} else {
				return nil, ErrInvalidRepository
			}
		}

		parentPath := filepath.Clean(filepath.Join(rootPath, ".."))
		if rootPath == parentPath {
			return nil, ErrNotRepository
		}
		rootPath = parentPath
	}

	configFilePath := filepath.Join(rootPath, DirNameConfig, FileNameConfig)

	fi, err := os.Stat(configFilePath)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, ErrInvalidRepository
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	c := DefaultRepositoryConfig()
	if _, err = toml.DecodeReader(configFile, c); err != nil {
		return nil, err
	}

	c.Normalize(rootPath)

	r, err := NewRepository(c)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// InitRepository initializes the input directory as a Repository.
func InitRepository(rootPath string) error {
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return err
	}

	configDirPath := filepath.Join(rootPath, DirNameConfig)

	fi, err := os.Stat(configDirPath)
	if err != nil {
		if err.(*os.PathError).Err == syscall.ENOENT {
			if err = os.Mkdir(configDirPath, 0777); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if !fi.IsDir() {
			return ErrInvalidRepository
		}
	}

	configFilePath := filepath.Join(configDirPath, FileNameConfig)

	fi, err = os.Stat(configFilePath)
	if err != nil {
		if err.(*os.PathError).Err != syscall.ENOENT {
			return err
		}
	} else {
		if fi.IsDir() {
			return ErrInvalidRepository
		}

		return ErrRepositoryExist
	}

	configFile, err := os.Create(configFilePath)
	if err != nil {
		return nil
	}
	defer configFile.Close()

	u, err := user.Current()
	if err != nil {
		return nil
	}

	c := DefaultRepositoryConfig()
	c.Username = u.Username

	enc := toml.NewEncoder(configFile)
	if err := enc.Encode(c); err != nil {
		return err
	}

	c.Normalize(rootPath)

	fi, err = os.Stat(c.DirDatabasePath)
	if err != nil {
		if err.(*os.PathError).Err == syscall.ENOENT {
			if err = os.Mkdir(c.DirDatabasePath, 0777); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if !fi.IsDir() {
			return ErrInvalidRepository
		}
	}

	fi, err = os.Stat(c.DirMailDataPath)
	if err != nil {
		if err.(*os.PathError).Err == syscall.ENOENT {
			if err = os.Mkdir(c.DirMailDataPath, 0700); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if !fi.IsDir() {
			return ErrInvalidRepository
		}
	}

	aliasDomainsFileName := filepath.Join(c.DirMailDataPath, FileNameAliasDomains)

	fi, err = os.Stat(aliasDomainsFileName)
	if err != nil {
		if err.(*os.PathError).Err != syscall.ENOENT {
			return err
		}
	} else {
		if fi.IsDir() {
			return ErrInvalidRepository
		}
	}

	aliasDomainsFile, err := os.OpenFile(aliasDomainsFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil
	}
	defer aliasDomainsFile.Close()

	return nil
}
