package context

type ArgoCredentials struct {
	Username  string
	Password  string
	Host      string
	Token     string
	BasicAuth bool
}

type OutConfig struct {
	CommandsFile        string
	ExportOutUrlCommand string
}

var PluginArgoCredentials = &ArgoCredentials{}
var PluginOutConfig = &OutConfig{}
