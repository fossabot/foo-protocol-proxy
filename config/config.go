package config

type (
	// Configuration type wraps configuration data.
	Configuration struct {
		Listening    string
		Forwarding   string
		HTTPAddress  string
		RecoveryPath string
	}
)
