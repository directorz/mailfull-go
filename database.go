package mailfull

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// GenerateDatabases generates databases from the MailData directory.
func (r *Repository) GenerateDatabases(md *MailData) error {
	sort.Slice(md.Domains, func(i, j int) bool { return md.Domains[i].Name() < md.Domains[j].Name() })
	sort.Slice(md.AliasDomains, func(i, j int) bool { return md.AliasDomains[i].Name() < md.AliasDomains[j].Name() })

	for _, domain := range md.Domains {
		sort.Slice(domain.Users, func(i, j int) bool { return domain.Users[i].Name() < domain.Users[j].Name() })
		sort.Slice(domain.AliasUsers, func(i, j int) bool { return domain.AliasUsers[i].Name() < domain.AliasUsers[j].Name() })
	}

	// Generate files
	if err := r.generateDbDomains(md); err != nil {
		return err
	}
	if err := r.generateDbDestinations(md); err != nil {
		return err
	}
	if err := r.generateDbMaildirs(md); err != nil {
		return err
	}
	if err := r.generateDbLocaltable(md); err != nil {
		return err
	}
	if err := r.generateDbForwards(md); err != nil {
		return err
	}
	if err := r.generateDbPasswords(md); err != nil {
		return err
	}

	// Generate DBs
	if err := exec.Command(r.CmdPostmap, filepath.Join(r.DirDatabasePath, FileNameDbDomains)).Run(); err != nil {
		return err
	}
	if err := exec.Command(r.CmdPostmap, filepath.Join(r.DirDatabasePath, FileNameDbDestinations)).Run(); err != nil {
		return err
	}
	if err := exec.Command(r.CmdPostmap, filepath.Join(r.DirDatabasePath, FileNameDbMaildirs)).Run(); err != nil {
		return err
	}
	if err := exec.Command(r.CmdPostmap, filepath.Join(r.DirDatabasePath, FileNameDbLocaltable)).Run(); err != nil {
		return err
	}
	if err := exec.Command(r.CmdPostalias, filepath.Join(r.DirDatabasePath, FileNameDbForwards)).Run(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) generateDbDomains(md *MailData) error {
	dbDomains, err := os.Create(filepath.Join(r.DirDatabasePath, FileNameDbDomains))
	if err != nil {
		return err
	}
	if err := dbDomains.Chown(r.uid, r.gid); err != nil {
		return err
	}
	defer dbDomains.Close()

	for _, domain := range md.Domains {
		if domain.Disabled() {
			continue
		}

		if _, err := fmt.Fprintf(dbDomains, "%s virtual\n", domain.Name()); err != nil {
			return err
		}
	}

	for _, aliasDomain := range md.AliasDomains {
		if _, err := fmt.Fprintf(dbDomains, "%s virtual\n", aliasDomain.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) generateDbDestinations(md *MailData) error {
	dbDestinations, err := os.Create(filepath.Join(r.DirDatabasePath, FileNameDbDestinations))
	if err != nil {
		return err
	}
	if err := dbDestinations.Chown(r.uid, r.gid); err != nil {
		return err
	}
	defer dbDestinations.Close()

	for _, domain := range md.Domains {
		if domain.Disabled() {
			continue
		}

		// ho-ge.example.com -> ho_ge.example.com
		underscoredDomainName := domain.Name()
		underscoredDomainName = strings.Replace(underscoredDomainName, `-`, `_`, -1)

		for _, user := range domain.Users {
			userName := user.Name()
			if cu := domain.CatchAllUser; cu != nil && cu.Name() == user.Name() {
				userName = ""
			}

			if len(user.Forwards()) > 0 {
				if _, err := fmt.Fprintf(dbDestinations, "%s@%s %s|%s\n", userName, domain.Name(), underscoredDomainName, user.Name()); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(dbDestinations, "%s@%s %s@%s\n", userName, domain.Name(), user.Name(), domain.Name()); err != nil {
					return err
				}
			}

			for _, aliasDomain := range md.AliasDomains {
				if aliasDomain.Target() == domain.Name() {
					if _, err := fmt.Fprintf(dbDestinations, "%s@%s %s@%s\n", userName, aliasDomain.Name(), user.Name(), domain.Name()); err != nil {
						return err
					}
				}
			}
		}

		for _, aliasUser := range domain.AliasUsers {
			if _, err := fmt.Fprintf(dbDestinations, "%s@%s %s\n", aliasUser.Name(), domain.Name(), strings.Join(aliasUser.Targets(), ",")); err != nil {
				return err
			}

			for _, aliasDomain := range md.AliasDomains {
				if aliasDomain.Target() == domain.Name() {
					if _, err := fmt.Fprintf(dbDestinations, "%s@%s %s@%s\n", aliasUser.Name(), aliasDomain.Name(), aliasUser.Name(), domain.Name()); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (r *Repository) generateDbMaildirs(md *MailData) error {
	dbMaildirs, err := os.Create(filepath.Join(r.DirDatabasePath, FileNameDbMaildirs))
	if err != nil {
		return err
	}
	if err := dbMaildirs.Chown(r.uid, r.gid); err != nil {
		return err
	}
	defer dbMaildirs.Close()

	for _, domain := range md.Domains {
		if domain.Disabled() {
			continue
		}

		for _, user := range domain.Users {
			if _, err := fmt.Fprintf(dbMaildirs, "%s@%s %s/%s/Maildir/\n", user.Name(), domain.Name(), domain.Name(), user.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Repository) generateDbLocaltable(md *MailData) error {
	dbLocaltable, err := os.Create(filepath.Join(r.DirDatabasePath, FileNameDbLocaltable))
	if err != nil {
		return err
	}
	if err := dbLocaltable.Chown(r.uid, r.gid); err != nil {
		return err
	}
	defer dbLocaltable.Close()

	for _, domain := range md.Domains {
		if domain.Disabled() {
			continue
		}

		// ho-ge.example.com -> ho_ge\.example\.com
		escapedDomainName := domain.Name()
		escapedDomainName = strings.Replace(escapedDomainName, `-`, `_`, -1)
		escapedDomainName = strings.Replace(escapedDomainName, `.`, `\.`, -1)

		if _, err := fmt.Fprintf(dbLocaltable, "/^%s\\|.*$/ local\n", escapedDomainName); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) generateDbForwards(md *MailData) error {
	dbForwards, err := os.Create(filepath.Join(r.DirDatabasePath, FileNameDbForwards))
	if err != nil {
		return err
	}
	if err := dbForwards.Chown(r.uid, r.gid); err != nil {
		return err
	}
	defer dbForwards.Close()

	for _, domain := range md.Domains {
		if domain.Disabled() {
			continue
		}

		// ho-ge.example.com -> ho_ge.example.com
		underscoredDomainName := domain.Name()
		underscoredDomainName = strings.Replace(underscoredDomainName, `-`, `_`, -1)

		for _, user := range domain.Users {
			if len(user.Forwards()) > 0 {
				if _, err := fmt.Fprintf(dbForwards, "%s|%s:%s\n", underscoredDomainName, user.Name(), strings.Join(user.Forwards(), ",")); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(dbForwards, "%s|%s:/dev/null\n", underscoredDomainName, user.Name()); err != nil {
					return err
				}
			}
		}
	}

	// drop real user
	if _, err := fmt.Fprintf(dbForwards, "%s:/dev/null\n", r.Username); err != nil {
		return err
	}

	return nil
}

func (r *Repository) generateDbPasswords(md *MailData) error {
	dbPasswords, err := os.Create(filepath.Join(r.DirDatabasePath, FileNameDbPasswords))
	if err != nil {
		return err
	}
	if err := dbPasswords.Chown(r.uid, r.gid); err != nil {
		return err
	}
	defer dbPasswords.Close()

	for _, domain := range md.Domains {
		if domain.Disabled() {
			continue
		}

		for _, user := range domain.Users {
			if _, err := fmt.Fprintf(dbPasswords, "%s@%s:%s\n", user.Name(), domain.Name(), user.HashedPassword()); err != nil {
				return err
			}
		}
	}

	return nil
}
