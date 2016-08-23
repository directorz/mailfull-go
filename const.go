package mailfull

// Filenames that are contained in the Repository.
const (
	DirNameConfig  = ".mailfull"
	FileNameConfig = "config"

	FileNameDomainDisable = ".vdomaindisable"
	FileNameAliasDomains  = ".valiasdomains"
	FileNameUsersPassword = ".vpasswd"
	FileNameUserForwards  = ".forward"
	FileNameAliasUsers    = ".valiases"
	FileNameCatchAllUser  = ".vcatchall"

	FileNameDbDomains      = "domains"
	FileNameDbDestinations = "destinations"
	FileNameDbMaildirs     = "maildirs"
	FileNameDbLocaltable   = "localtable"
	FileNameDbForwards     = "forwards"
	FileNameDbPasswords    = "vpasswd"
)

// NeverMatchHashedPassword is hash string that is never match with any password.
const NeverMatchHashedPassword = "{SSHA}!!"
