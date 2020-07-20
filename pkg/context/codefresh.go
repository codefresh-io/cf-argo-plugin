package context

type CodefreshCredentials struct {
	Token string
	Host  string
}

var PluginCodefreshCredentials = &CodefreshCredentials{}
