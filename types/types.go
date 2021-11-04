package types

// Config is a struct of the framework configuration
type Config struct {
	Principal      string
	LogLevel       string
	MinVersion     string
	AppName        string
	EnableSyslog   bool
	Hostname       string
	Listen         string
	Domain         string
	Credentials    UserCredentials
	PrefixHostname string
	PrefixTaskName string
}

// UserCredentials - The Username and Password to authenticate against this framework
type UserCredentials struct {
	Username string
	Password string
}
