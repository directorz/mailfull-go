package mailfull

// AliasDomain represents a AliasDomain.
type AliasDomain struct {
	name   string
	target string
}

// NewAliasDomain creates a new AliasDomain instance.
func NewAliasDomain(name, target string) (*AliasDomain, error) {
	if !validAliasDomainName(name) {
		return nil, ErrInvalidAliasDomainName
	}
	if !validAliasDomainTarget(target) {
		return nil, ErrInvalidAliasDomainTarget
	}

	ad := &AliasDomain{
		name:   name,
		target: target,
	}

	return ad, nil
}

// Name returns name.
func (ad *AliasDomain) Name() string {
	return ad.name
}

// Target returns target.
func (ad *AliasDomain) Target() string {
	return ad.target
}
