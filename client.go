package elasticsearch_go

import (
	"fmt"
	"net/http"
)

// Client is a client used to send RESTful request to ElasticSearch
type Client struct {
	config Config

	// the active es addresses
	activePool []string

	// last used es server address
	// if it is active, the following request can be just sent to it
	lastUsedAddress string
}

// Config defines the ElasticSearch config
type Config struct {

	// the addresses of the ElasticSearch cluster
	Addresses []string

	// the basic auth username
	Username string

	// the basic auth password
	Password string
}

// defaultConfig is a default config that uses local ElasticSearch
var defaultConfig = Config{
	Addresses: []string{"http://127.0.0.1:9200"},
}

// DefaultClient creates a default ElasticSearch client
// which is conncted with the local ElasticSearch default
func DefaultClient() *Client {
	defaultClient := &Client{
		config: defaultConfig,
	}
	return defaultClient
}

// NewClient creates a new ElasticSearch client according to Config
func NewClient(config Config) *Client {
	if len(config.Addresses) == 0 {
		config = defaultConfig
	}
	c := &Client{
		config: config,
	}
	return c
}

// Ping pings es server
func (c *Client) Ping() bool {

	// TODO: how to choose the request url
	// Using the connection pool,
	// periodically ping all the links in the cluster,
	// maintaining a list of available connections,
	// and selecting a url from the connection pool to request
	req, err := http.NewRequest("GET", c.config.Addresses[0], nil)
	if err != nil {
		ZapLogger.Error(fmt.Sprintf("failed to create new request: %v", err))
		return false
	}

	if c.config.Username != "" {
		req.SetBasicAuth(c.config.Username, c.config.Password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ZapLogger.Error(fmt.Sprintf("failed to send es info request: %v", err))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		ZapLogger.Error(resp.Status)
		return false
	}

	return true
}
