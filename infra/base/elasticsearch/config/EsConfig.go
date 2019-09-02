package config

/**
Addresses []string // A list of Elasticsearch nodes to use.
Username  string   // Username for HTTP Basic Authentication.
Password  string   // Password for HTTP Basic Authentication.

Transport http.RoundTripper  // The HTTP transport object.
Logger    estransport.Logger // The logger object.
*/
type EsConfig struct {
	Addrs []string `yaml:"addrs"`
}
