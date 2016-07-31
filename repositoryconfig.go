package mailfull

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/BurntSushi/toml"
)

// Errors for the Repository.
var (
	ErrInvalidRepository = errors.New("invalid repository")
	ErrNotRepository     = errors.New("not a Mailfull repository (or any of the parent directories)")
	ErrRepositoryExist   = errors.New("a Mailfull repository exists")
)

// RepositoryConfig is used to configure a Repository.
type RepositoryConfig struct {
	DirDatabasePath string `toml:"dir_database"`
	DirMailDataPath string `toml:"dir_maildata"`
	Username        string `toml:"username"`
	Groupname       string `toml:"groupname"`
}

// DefaultRepositoryConfig returns a RepositoryConfig with default parameter.
func DefaultRepositoryConfig() *RepositoryConfig {
	c := &RepositoryConfig{
		DirDatabasePath: "./etc",
		DirMailDataPath: "./domains",
		Username:        "mailfull",
		Groupname:       "mailfull",
	}

	return c
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

	if !filepath.IsAbs(c.DirDatabasePath) {
		c.DirDatabasePath = filepath.Join(rootPath, c.DirDatabasePath)
	}
	if !filepath.IsAbs(c.DirMailDataPath) {
		c.DirMailDataPath = filepath.Join(rootPath, c.DirMailDataPath)
	}

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

	c := DefaultRepositoryConfig()

	enc := toml.NewEncoder(configFile)
	if err := enc.Encode(c); err != nil {
		return err
	}

	if !filepath.IsAbs(c.DirDatabasePath) {
		c.DirDatabasePath = filepath.Join(rootPath, c.DirDatabasePath)
	}
	if !filepath.IsAbs(c.DirMailDataPath) {
		c.DirMailDataPath = filepath.Join(rootPath, c.DirMailDataPath)
	}

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
			if err = os.Mkdir(c.DirMailDataPath, 0777); err != nil {
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

	aliasDomainFileName := filepath.Join(c.DirMailDataPath, FileNameAliasDomains)

	fi, err = os.Stat(aliasDomainFileName)
	if err != nil {
		if err.(*os.PathError).Err != syscall.ENOENT {
			return err
		}
	} else {
		if fi.IsDir() {
			return ErrInvalidRepository
		}
	}

	aliasDomainFile, err := os.Create(aliasDomainFileName)
	if err != nil {
		return nil
	}
	defer aliasDomainFile.Close()

	return nil
}
