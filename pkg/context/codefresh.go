package context

type CodefreshCredentials struct {
	Token       string
	Host        string
	Integration string
}

var PluginCodefreshCredentials = &CodefreshCredentials{}
