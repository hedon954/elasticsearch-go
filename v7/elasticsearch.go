package elasticsearch

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/valyala/fasthttp"
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
	req := c.getFasthttpReq()
	resp := fasthttp.AcquireResponse()
	err := fasthttp.Do(req, resp)
	if err != nil {
		ZapLogger.Error(fmt.Sprintf("failed to send es info request: %v", err))
		return false
	}
	defer resp.Reset()

	if resp.StatusCode() < 200 || resp.StatusCode() >= 400 {
		ZapLogger.Error(resp.String())
		return false
	}
	return true
}

// Version returns es version
func (c *Client) Version() string {
	req := c.getFasthttpReq()
	resp := fasthttp.AcquireResponse()
	defer resp.Reset()
	err := fasthttp.Do(req, resp)
	if err != nil {
		ZapLogger.Error(fmt.Sprintf("failed to send es info request: %v", err))
		return ""
	}
	defer resp.Reset()

	if resp.StatusCode() < 200 || resp.StatusCode() >= 400 {
		ZapLogger.Error(resp.String())
		return ""
	}

	ei := esInfo{}
	err = json.Unmarshal(resp.Body(), &ei)
	if err != nil {
		ZapLogger.Error(fmt.Sprintf("unmarshal resp body to esInfo failed: %v", err))
		return ""
	}

	return ei.Version.Number
}

// esInfo is the response body struct of http://127.0.0.1:9200
type esInfo struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
	Version     struct {
		Number                           string `json:"number"`
		BuildFlavor                      string `json:"build_flavor"`
		BuildType                        string `json:"build_type"`
		BuildHash                        string `json:"build_hash"`
		BuildDate                        string `json:"build_date"`
		BuildSnapshot                    bool   `json:"build_snapshot"`
		LuceneVersion                    string `json:"lucene_version"`
		MinimumWireCompatibilityVersion  string `json:"minimum_wire_compatibility_version"`
		MinimumIndexCompatibilityVersion string `json:"minimum_index_compatibility_version"`
	} `json:"version"`
	TagLine string `json:"tag_line"`
}

// getFasthttpReq builds a basic fasthttp.Request
func (c *Client) getFasthttpReq() *fasthttp.Request {
	uri := fasthttp.AcquireURI()
	uri.SetUsername(c.config.Username)
	uri.SetPassword(c.config.Password)
	req := fasthttp.AcquireRequest()

	// TODO: how to choose the request url
	// Using the connection pool,
	// periodically ping all the links in the cluster,
	// maintaining a list of available connections,
	// and selecting a url from the connection pool to request
	req.Header.SetMethod("GET")
	req.Header.SetRequestURI(c.config.Addresses[0])
	parse, _ := url.Parse(c.config.Addresses[0])
	uri.SetScheme(parse.Scheme)
	uri.SetHost(parse.Host)
	req.SetURI(uri)

	return req
}
