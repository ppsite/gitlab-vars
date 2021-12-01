package gitlab

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Gitlab struct {
	BaseUrl      string
	ApiPath      string
	Token        string
	Client       *http.Client
}

func NewGitlab(baseUrl string, apiPath string, token string, skipCertVerify bool) *Gitlab {
	config := &tls.Config{InsecureSkipVerify: skipCertVerify}
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: config,
	}
	client := &http.Client{Transport: tr}

	return &Gitlab{
		BaseUrl: baseUrl,
		ApiPath: apiPath,
		Token:   token,
		Client:  client,
	}
}

func (g *Gitlab) ResourceUrl(path string, params map[string]string) *url.URL {
	for key, val := range params {
		path = strings.Replace(path, key, val, -1)
	}

	u, err := url.Parse(g.BaseUrl + g.ApiPath + path)
	if err != nil {
		panic("Error while building gitlab url, unable to parse generated url")
	}

	return u
}

func (g *Gitlab) Request(method, url string, payload []byte) ([]byte, error) {
  var req *http.Request
	var err error

  if payload == nil {
    req, err = http.NewRequest(method, url, nil)
  }else{
    req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
  }

  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  req.Header.Add("PRIVATE-TOKEN", g.Token)
  if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}

  res, err := g.Client.Do(req)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  fmt.Println(string(body))
  return body, err
}


func (g *Gitlab) execRequest(method, url string, body []byte) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		reader := bytes.NewReader(body)
		req, err = http.NewRequest(method, url, reader)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	req.Header.Add("Private-Token", g.Token)
	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}

	if err != nil {
		panic("Error while building gitlab request")
	}

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Client.Do error: %q", err)
	}

	return resp, err
}

func (g *Gitlab) buildAndExecRequest(method, u string, body []byte) ([]byte, error) {
	resp, err := g.execRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return contents, err
}
