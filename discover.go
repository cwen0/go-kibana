package kibana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/juju/errors"
	"gopkg.in/sakura-internet/go-rison.v3"
)

const (
	discoverBaseURL = "app/kibana#/discover"
	shortURLAPI     = "api/shorten_url"
)

type Discover interface {
	GenerateURL(search *SearchSource) (string, error)
	ShortURL(url string) (string, error)
}

type discover struct {
	config *Config

	httpClient *http.Client
}

func NewDiscover(conf *Config, httpClient *http.Client) Discover {
	return &discover{
		config:     conf,
		httpClient: httpClient,
	}
}

func (d *discover) GenerateURL(search *SearchSource) (string, error) {
	// set default value
	g := []byte("()")
	var err error
	if search.Time != nil {
		g, err = rison.Marshal(search.Time, rison.Rison)
		if err != nil {
			return "", err
		}
	}

	a, err := rison.Marshal(search.Search, rison.Rison)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s?_g=%s&_a=%s", d.config.Address, discoverBaseURL, g, a)
	return url, nil
}

func (d *discover) ShortURL(url string) (string, error) {
	url = strings.TrimPrefix(url, d.config.Address)
	data, err := json.Marshal(map[string]string{"url": url})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s", d.config.Address, shortURLAPI),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("kbn-xsrf", "true")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.BadRequestf("status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	repData := make(map[string]string)
	if err := json.Unmarshal(body, &repData); err != nil {
		return "", err
	}

	url, ok := repData["urlId"]
	if !ok {
		return "", errors.New("short url is empty")
	}

	return fmt.Sprintf("%s/goto/%s", d.config.Address, url), nil
}

type SearchSource struct {
	Time   *TimeFields
	Search SearchFields
}

type TimeFields struct {
	Refresh RefreshInterval `json:"refreshInterval"`
	Time    TimeRange       `json:"time"`
}

type RefreshInterval struct {
	Pause bool  `json:"pause"`
	value int64 `json:"value"`
}

type TimeRange struct {
	From string `json:"from"`
	Mode string `json:"mode"`
	To   string `json:"to"`
}

type SearchFields struct {
	Columns  []string  `json:"columns"`
	Filters  []Filter  `json:"filters"`
	Sort     []string  `json:"sort"`
	Query    QueryMeta `json:"json"`
	Index    string    `json:"index"`
	Interval string    `json:"interval"`
}

type Filter struct {
	State State           `json:"$state"`
	Meta  FilterMeta      `json:"meta"`
	Query FilterQueryMeta `json:"query"`
}

type State struct {
	Store string `json:"store"`
}

type FilterMeta struct {
	Alias   string       `json:"alias"`
	Disable bool         `json:"disable"`
	Index   string       `json:"index"`
	Key     string       `json:"key"`
	Negate  bool         `json:"negate"`
	Params  FilterParams `json:"params"`
	Type    string       `json:"type"`
	Value   string       `json:"value"`
}

type FilterParams struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

type FilterQueryMeta struct {
	Match map[string]FilterParams `json:"match"`
}

type QueryMeta struct {
	Language string `json:"language"`
	Query    string `json:"query"`
}
