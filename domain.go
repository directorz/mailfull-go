package mailfull

// Domain represents a Domain.
type Domain struct {
	name         string
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
