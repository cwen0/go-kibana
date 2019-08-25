package kibana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/juju/errors"
	"gopkg.in/sakura-internet/go-rison.v3"
)

const (
	discoverBaseURL = "/app/kibana#/discover"
	shortURLAPI     = "/api/shorten_url"
)

type Discover interface {
	GenerateURL(time *TimeFields, search *SearchFields) (string, error)
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

func (d *discover) GenerateURL(time *TimeFields, search *SearchFields) (string, error) {
	g, err := rison.Marshal(time, rison.Rison)
	if err != nil {
		return "", err
	}

	a, err := rison.Marshal(search, rison.Rison)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s/?_g=%s&_a=%s", d.config.Address, discoverBaseURL, g, a)
	return url, nil
}

func (d *discover) ShortURL(url string) (string, error) {
	data, err := json.Marshal(map[string]string{"url": url})
	if err != nil {
		return "", err
	}

	rep, err := d.httpClient.Post(fmt.Sprintf("%s/%s", d.config.Address, shortURLAPI),
		"application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer rep.Body.Close()

	if rep.StatusCode != http.StatusOK {
		return "", errors.BadRequestf("status: %s", rep.Status)
	}

	body, err := ioutil.ReadAll(rep.Body)
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
	Columns []string  `json:"columns"`
	Filters []Filter  `json:"filters"`
	Sort    []string  `json:"sort"`
	Query   QueryMeta `json:"json"`
}

type Filter struct {
	State State `json:"$state"`
}

type State struct {
	store string `json:"store"`
}

type FilterMeta struct {
	Alias   string `json:"alias"`
	Disable bool   `json:"disable"`
	Index   string `json:"index"`
	Key     string `json:"key"`
	Negate  bool   `json:"negate"`
}

type FilterQueryMeta struct {
	Match map[string]FilterQueryMatchMeta `json:"match"`
}

type FilterQueryMatchMeta struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

type QueryMeta struct {
	Language string `json:"language"`
	Query    string `json:"query"`
}
