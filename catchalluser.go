package mailfull

// CatchAllUser represents a CatchAllUser.
type CatchAllUser struct {
	name string
}

// NewCatchAllUser creates a new CatchAllUser instance.
func NewCatchAllUser(name string) (*CatchAllUser, error) {
	if !validCatchAllUserName(name) {
		return nil, ErrInvalidCatchAllUserName
	}

	cu := &CatchAllUser{
		name: name,
	}

	return cu, nil
}

// Name returns name.
func (cu *CatchAllUser) Name() string {
	return cu.name
}
