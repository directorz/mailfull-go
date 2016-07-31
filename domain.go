package mailfull

// Domain represents a Domain.
type Domain struct {
	name         string
	Users        []*User
	AliasUsers   []*AliasUser
	CatchAllUser *CatchAllUser
}

// NewDomain creates a new Domain instance.
func NewDomain(name string) (*Domain, error) {
	if !validDomainName(name) {
		return nil, ErrInvalidDomainName
	}

	d := &Domain{
		name: name,
	}

	return d, nil
}

// Name returns name.
func (d *Domain) Name() string {
	return d.name
}
