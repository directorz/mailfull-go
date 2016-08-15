package mailfull

// MailData represents a MailData.
type MailData struct {
	Domains      []*Domain
	AliasDomains []*AliasDomain
}

// MailData returns a MailData.
func (r *Repository) MailData() (*MailData, error) {
	domains, err := r.Domains()
	if err != nil {
		return nil, err
	}

	aliasDomains, err := r.AliasDomains()
	if err != nil {
		return nil, err
	}

	for _, domain := range domains {
		users, err := r.Users(domain.Name())
		if err != nil {
			return nil, err
		}
		domain.Users = users

		aliasUsers, err := r.AliasUsers(domain.Name())
		if err != nil {
			return nil, err
		}
		domain.AliasUsers = aliasUsers

		catchAllUser, err := r.CatchAllUser(domain.Name())
		if err != nil {
			return nil, err
		}
		domain.CatchAllUser = catchAllUser
	}

	mailData := &MailData{
		Domains:      domains,
		AliasDomains: aliasDomains,
	}

	return mailData, nil
}
