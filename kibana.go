package kibana

import "net/http"

// Config represents the client configuration.
type Config struct {
	// base address of kibana
	Address string
}

// Client represents the kibana client
type Client struct {
	Discover

	Config     *Config
	httpClient *http.Client
}

// NewClient returns a new Kibana API client. If a nil httpClient is
// provided, a new http.Client will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(conf *Config, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	client := &Client{
		Config:     conf,
		httpClient: httpClient,
	}

	client.Discover = NewDiscover(conf, httpClient)

	return client
}
