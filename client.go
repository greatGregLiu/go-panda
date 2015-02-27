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
	"strconv"
	"strings"
	"time"
)

var queryFixer = strings.NewReplacer("+", "%20", "%5B", "[", "%5D", "]", "%7E", "~")

// Client TODO(jszwec)
type Client interface {
	// Get TODO
	Get(string, url.Values) ([]byte, error)
	// Post TODO
	Post(string, string, url.Values, io.Reader) ([]byte, error)
	// Put TODO
	Put(string, string, url.Values, io.Reader) ([]byte, error)
	// Delete TODO
	Delete(string) ([]byte, error)
	// SetCloud TODO
	SetCloud(string)
}

type client struct {
	host      string
	port      int
	ver       string
	cloudId   string
	accessKey string
	secretKey []byte
	cl        *http.Client
}

// NewClientCustom TODO(jszwec)
func NewClientCustom(host string, port int, ver, cloudId, accessKey,
	secretKey string) Client {
	return &client{
		host,
		port,
		ver,
		cloudId,
		accessKey,
		[]byte(secretKey),
		&http.Client{},
	}
}

// NewClient TODO(jszwec)
func NewClient(cloudId, accessKey, secretKey string) Client {
	return NewClientCustom("api.pandastream.com", 80, "v2", cloudId,
		accessKey, secretKey)
}

// NewClientEU(jszwec)
func NewClientEU(cloudId, accessKey, secretKey string) Client {
	return NewClientCustom("api-eu.pandastream.com", 80, "v2", cloudId,
		accessKey, secretKey)
}

func (cl *client) addAuthParams(v url.Values, t time.Time) {
	v.Set("access_key", cl.accessKey)
	v.Set("cloud_id", cl.cloudId)
	v.Set("timestamp", t.Format(time.RFC3339Nano))
}

func (cl *client) fixQuery(s string) string {
	return queryFixer.Replace(s)
}

func (cl *client) buildSignature(v url.Values, method, u string,
	t time.Time) (sign string, err error) {
	toSign := fmt.Sprintf("%s\n%s\n%s\n%s", method, cl.host, u, cl.fixQuery(v.Encode()))
	mac := hmac.New(sha256.New, cl.secretKey)
	if _, err = mac.Write([]byte(toSign)); err == nil {
		sign = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	}
	return
}

func (cl *client) buildURL(v url.Values, p string) *url.URL {
	scheme := "http"
	if cl.port == 443 {
		scheme = "https"
	}
	return &url.URL{
		Scheme:   scheme,
		Host:     cl.host + ":" + strconv.Itoa(cl.port),
		Path:     path.Join(cl.ver, p),
		RawQuery: v.Encode(),
	}
}

func (cl *client) do(method, path, cntType string,
	params url.Values, r io.Reader) (b []byte, err error) {
	if params == nil {
		params = url.Values{}
	}
	t := time.Now().UTC()
	cl.addAuthParams(params, t)
	sign, err := cl.buildSignature(params, method, path, t)
	if err != nil {
		return
	}
	params.Set("signature", sign)
	req, err := http.NewRequest(method, cl.buildURL(params, path).String(), r)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", cntType)
	resp, err := cl.cl.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		e := &PandaError{Code: resp.StatusCode}
		if err = json.Unmarshal(b, e); err != nil {
			return nil, e
		}
		err = e
	}
	return
}

func (cl *client) Get(url string, params url.Values) ([]byte, error) {
	return cl.do("GET", url, "", params, nil)
}

func (cl *client) Post(url, cntType string, params url.Values, r io.Reader) ([]byte, error) {
	return cl.do("POST", url, cntType, params, r)
}

func (cl *client) Put(url, cntType string, params url.Values, r io.Reader) ([]byte, error) {
	return cl.do("PUT", url, cntType, params, r)
}

func (cl *client) Delete(url string) ([]byte, error) {
	return cl.do("DELETE", url, "", nil, nil)
}

func (cl *client) SetCloud(cloudId string) {
	cl.cloudId = cloudId
}
