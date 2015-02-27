package panda

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	HostUS  = "api.pandastream.com"
	HostEU  = "api-eu.pandastream.com"
	HostGCE = "api-gce.pandastream.com"
)

// ClientOptions hold credentials required for authenticating requests to the Panda Cloud
type ClientOptions struct {
	CloudID   string
	AccessKey string
	SecretKey string
}

var queryFixer = strings.NewReplacer("+", "%20", "%5B", "[", "%5D", "]", "%7E", "~")

// Client is capable of sending signed requests to the Panda Cloud
type Client struct {
	Host       string
	Options    *ClientOptions
	HTTPClient *http.Client
}

func (cl *Client) host() string {
	if cl.Host != "" {
		return cl.Host
	}
	return HostUS
}

func (cl *Client) httpclient() *http.Client {
	if cl.HTTPClient != nil {
		return cl.HTTPClient
	}
	return http.DefaultClient
}

func (cl *Client) addAuthParams(v url.Values, t time.Time) {
	v.Set("access_key", cl.Options.AccessKey)
	v.Set("cloud_id", cl.Options.CloudID)
	v.Set("timestamp", t.Format(time.RFC3339Nano))
}

func (cl *Client) fixQuery(s string) string {
	return queryFixer.Replace(s)
}

func (cl *Client) buildSignature(v url.Values, method, u string) (sign string, err error) {
	toSign := fmt.Sprintf("%s\n%s\n%s\n%s", method, cl.host(), u, cl.fixQuery(v.Encode()))
	mac := hmac.New(sha256.New, []byte(cl.Options.SecretKey))
	if _, err = mac.Write([]byte(toSign)); err == nil {
		sign = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	}
	return
}

func (cl *Client) buildURL(v url.Values, urlPath string) *url.URL {
	scheme := "http"
	if strings.HasSuffix(cl.host(), ":443") {
		scheme = "https"
	}
	return &url.URL{
		Scheme:   scheme,
		Host:     cl.host(),
		Path:     path.Join("v2", urlPath),
		RawQuery: v.Encode(),
	}
}

func (cl *Client) do(method, path, cntType string,
	params url.Values, r io.Reader) (b []byte, err error) {
	if params == nil {
		params = url.Values{}
	}
	if err = cl.SignParams(method, path, params); err != nil {
		return
	}
	req, err := http.NewRequest(method, cl.buildURL(params, path).String(), r)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", cntType)
	resp, err := cl.httpclient().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		e := &Error{Code: resp.StatusCode}
		if err = json.Unmarshal(b, e); err != nil {
			return nil, e
		}
		err = e
	}
	return
}

// SignParams signs given parameters by adding required authorization
// fields and values. Params cannot be nil.
func (cl *Client) SignParams(method, path string, params url.Values) error {
	if params == nil {
		panic("params cannot be nil!")
	}
	cl.addAuthParams(params, time.Now().UTC())
	s, err := cl.buildSignature(params, method, path)
	if err != nil {
		return err
	}
	params.Set("signature", s)
	return nil
}

// Get issues a signed GET request to the Panda Cloud
func (cl *Client) Get(url string, params url.Values) ([]byte, error) {
	return cl.do("GET", url, "", params, nil)
}

// Post issues a signed POST request to the Panda Cloud and creates content based on
// the given params
func (cl *Client) Post(url, cntType string, params url.Values, r io.Reader) ([]byte, error) {
	return cl.do("POST", url, cntType, params, r)
}

// Put issues a signed PUT request to the Panda Cloud and updates object according to
// given params
func (cl *Client) Put(url, cntType string, params url.Values, r io.Reader) ([]byte, error) {
	return cl.do("PUT", url, cntType, params, r)
}

// Delete issues a signed DELETE request to the Panda Cloud and deletes content under
// the given url
func (cl *Client) Delete(url string) ([]byte, error) {
	return cl.do("DELETE", url, "", nil, nil)
}
