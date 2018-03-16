package mailfull

import (
	"errors"
	"regexp"
)

// Errors for incorrect format.
var (
	ErrInvalidDomainName        = errors.New("Domain: name incorrect format")
	ErrInvalidAliasDomainName   = errors.New("AliasDomain: name incorrect format")
	ErrInvalidAliasDomainTarget = errors.New("AliasDomain: target incorrect format")
	ErrInvalidUserName          = errors.New("User: name incorrect format")
	ErrInvalidAliasUserName     = errors.New("AliasUser: name incorrect format")
	ErrInvalidAliasUserTarget   = errors.New("AliasUser: target incorrect format")
	ErrInvalidCatchAllUserName  = errors.New("CatchAllUser: name incorrect format")
)

// validDomainName returns true if the input is correct format.
func validDomainName(name string) bool {
	return regexp.MustCompile(`^([A-Za-z0-9\-]+\.)*[A-Za-z]+$`).MatchString(name)
}

// validAliasDomainName returns true if the input is correct format.
func validAliasDomainName(name string) bool {
	return regexp.MustCompile(`^([A-Za-z0-9\-]+\.)*[A-Za-z]+$`).MatchString(name)
}

// validAliasDomainTarget returns true if the input is correct format.
func validAliasDomainTarget(target string) bool {
	return regexp.MustCompile(`^([A-Za-z0-9\-]+\.)*[A-Za-z]+$`).MatchString(target)
}

// validUserName returns true if the input is correct format.
func validUserName(name string) bool {
	return regexp.MustCompile(`^[^\.\s@][^\s@]+$`).MatchString(name)
}

// validAliasUserName returns true if the input is correct format.
func validAliasUserName(name string) bool {
	return regexp.MustCompile(`^[^\.\s@][^\s@]+$`).MatchString(name)
}

// validAliasUserTarget returns true if the input is correct format.
func validAliasUserTarget(target string) bool {
	return regexp.MustCompile(`^[^\.\s@][^\s@]+@([A-Za-z0-9\-]+\.)*[A-Za-z]+$`).MatchString(target)
}

// validCatchAllUserName returns true if the input is correct format.
func validCatchAllUserName(name string) bool {
	return regexp.MustCompile(`^[^\.\s@][^\s@]+$`).MatchString(name)
}
