package context



type ArgoCredentials struct {
	Username string
	Password string
	Host string
}

var PluginArgoCredentials = &ArgoCredentials{}