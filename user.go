package mailfull

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
	if !validUserName(name) {
		return nil, ErrInvalidUserName
	}

	u := &User{
		name:           name,
		hashedPassword: hashedPassword,
		forwards:       forwards,
	}

	return u, nil
}

// Name returns name.
func (u *User) Name() string {
	return u.name
}

// HashedPassword returns hashedPassword.
func (u *User) HashedPassword() string {
	return u.hashedPassword
}

// Forwards returns forwards.
func (u *User) Forwards() []string {
	return u.forwards
}
