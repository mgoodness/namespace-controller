package pimapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1.0"
	baseURL        = "http://dev1.pim.nonprod-tmaws.io"
	basePath       = "tkm"
	userAgent      = "pimapi/" + libraryVersion
	mediaType      = "application/json"
)

type Config struct {
	Token      string
	APIVersion string
	BaseURL    string
	Cookies    []*http.Cookie
}

type Client struct {
	// HTTP client used to communicate with the PIM API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Config: Token, APIVersion, BaseURL, Cookies
	Config *Config

	// Services for various endpoints
	Asset             AssetService
	AssetType         AssetTypeService
	AttributeGroup    AttributeGroupService
	AttributeGroupMap AttributeGroupMapService
	AttributeType     AttributeTypeService
	AttributeTypeMap  AttributeTypeMapService
	ExternalSource    ExternalSourceService
	Link              LinkService
	LinkRule          LinkRuleService
	ProductCode       ProductCodeService
	ValidationType    ValidationTypeService
}

type Page struct {
	PageNumber     int64 `json:"pageNumber,omitempty"`
	ResultsPerPage int64 `json:"resultsPerPage,omitempty"`
	TotalCount     int64 `json:"totalCount,omitempty"`
	TotalPages     int64 `json:"totalPages,omitempty"`
}

type IdsResponse struct {
	Ids     []int64 `json:"ids"`
	message string  `json:"message,omitempty"`
}

type IdsRequest struct {
	Ids []int64 `json:"ids"`
}

func NewClient(config *Config) *Client {

	// Default to API Version 1
	if config.APIVersion == "" {
		config.APIVersion = "v1"
	}

	// Default to dev1
	if config.BaseURL == "" {
		config.BaseURL = baseURL
	}

	baseURL, _ := url.Parse(config.BaseURL)
	// add Version to URL base path
	baseURL.Path = basePath + "/" + config.APIVersion + "/"

	cookie := http.Cookie{Name: "token", Value: config.Token}
	config.Cookies = []*http.Cookie{&cookie}

	c := &Client{
		client:    http.DefaultClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
		Config:    config,
	}

	c.Asset = AssetService{client: c}
	c.AssetType = AssetTypeService{client: c}
	c.AttributeGroup = AttributeGroupService{client: c}
	c.AttributeGroupMap = AttributeGroupMapService{client: c}
	c.AttributeType = AttributeTypeService{client: c}
	c.AttributeTypeMap = AttributeTypeMapService{client: c}
	c.ExternalSource = ExternalSourceService{client: c}
	c.Link = LinkService{client: c}
	c.LinkRule = LinkRuleService{client: c}
	c.ProductCode = ProductCodeService{client: c}
	c.ValidationType = ValidationTypeService{client: c}

	return c
}

func (c *Client) SaveCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		c.Config.Cookies = append(c.Config.Cookies, cookie)
	}
}

// Creates an http.Request to the API . A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL. Relative URLS should always be specified without a preceding slash.
func (c *Client) NewRequest(method, urlStr string, body []byte) (req *http.Request, err error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	if body != nil {
		req, err = http.NewRequest(method, u.String(), bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, u.String(), nil)
	}

	if err != nil {
		return nil, err
	}

	/*
	   fmt.Printf("\n================================\n%s %s\n", method, u.String())
	   s := string(body[:])
	   fmt.Printf("BODY: %s\n", s)
	   fmt.Println("================================")
	*/

	for _, cookie := range c.Config.Cookies {
		req.AddCookie(cookie)
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", userAgent)
	return req, nil
}

// Send the request to the API, the response is sent to the caller via the `into` interface
func (c *Client) Do(req *http.Request, into interface{}) error {
	var (
		err  error
		resp *http.Response
		b    []byte
	)

	resp, err = c.client.Do(req)

	defer func() {
		if rerr := resp.Body.Close(); rerr == nil {
			err = rerr
		}
	}()
	if into != nil {
		if w, ok := into.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return err
			}
		} else {
			err := json.NewDecoder(resp.Body).Decode(into)
			if err != nil {
				return err
			}
		}
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s %s\n", resp.Status, b)
		return err
	}
	// Save JsessionID and ELB related cookies
	c.SaveCookies(resp.Cookies())

	return err
}

// GET By Id calls
func (c *Client) ID(path, id string, into interface{}) error {
	path = fmt.Sprintf("%s/%s", path, id)
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	err = c.Do(req, &into)
	if err != nil {
		return err
	}

	return nil
}

// GET calls
func (c *Client) Get(path string, params map[string]string, into interface{}) error {
	if len(params) > 0 {
		path = addOptions(path, params)
	}
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	err = c.Do(req, &into)
	if err != nil {
		return err
	}

	return nil
}

// POST calls
func (c *Client) Create(path string, createRequest interface{}) (IdsResponse, error) {
	var (
		id   IdsResponse
		err  error
		body []byte
	)

	if createRequest == nil {
		err = fmt.Errorf("Create request cannot be `nil`")
		return id, err
	}

	body, err = json.Marshal(createRequest)
	if err != nil {
		return id, err
	}

	req, err := c.NewRequest("POST", path, body)
	if err != nil {
		return id, err
	}

	err = c.Do(req, &id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// PUT calls
func (c *Client) Update(path string, updateRequest interface{}) (IdsResponse, error) {
	var (
		id   IdsResponse
		err  error
		body []byte
	)

	if updateRequest == nil {
		err = fmt.Errorf("Update request cannot be `nil`")
		return id, err
	}

	body, err = json.Marshal(updateRequest)
	if err != nil {
		return id, err
	}

	req, err := c.NewRequest("PUT", path, body)
	if err != nil {
		return id, err
	}

	err = c.Do(req, &id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// DELETE calls
func (c *Client) Delete(path string, deleteRequest interface{}) (IdsResponse, error) {
	var (
		id   IdsResponse
		body []byte
		err  error
	)

	if deleteRequest == nil {
		err = fmt.Errorf("Delete request cannot be `nil`")
		return id, err
	}

	body, err = json.Marshal(deleteRequest)
	if err != nil {
		return id, err
	}

	req, err := c.NewRequest("DELETE", path, body)
	if err != nil {
		return id, err
	}

	err = c.Do(req, &id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// helper function to construct url parameters
func addOptions(basePath string, p map[string]string) string {
	// Specify URL Parameters
	params := url.Values{}
	for k, v := range p {
		params.Add(k, v)
	}

	path := basePath + "?" + params.Encode()
	return path
}
