package mailfull

import "errors"

// Errors for parameter.
var (
	ErrNotEnoughAliasUserTargets = errors.New("AliasUser: targets not enough")
)

// AliasUser represents a AliasUser.
type AliasUser struct {
	name    string
	targets []string
}

// NewAliasUser creates a new AliasUser instance.
func NewAliasUser(name string, targets []string) (*AliasUser, error) {
	if !validAliasUserName(name) {
		return nil, ErrInvalidAliasUserName
	}

	if len(targets) < 1 {
		return nil, ErrNotEnoughAliasUserTargets
	}

	for _, target := range targets {
		if !validAliasUserTarget(target) {
			return nil, ErrInvalidAliasUserTarget
		}
	}

	au := &AliasUser{
		name:    name,
		targets: targets,
	}

	return au, nil
}

// Name returns name.
func (au *AliasUser) Name() string {
	return au.name
}

// Targets returns targets.
func (au *AliasUser) Targets() []string {
	return au.targets
}
