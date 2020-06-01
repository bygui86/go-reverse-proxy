package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

const (
	hostHeaderKey       = "Host"
	additionalHeaderKey = "X-Forwarded-Host"

	rootEndpoint = "/"
)

type ReverseProxy struct {
	proxy *httputil.ReverseProxy
	port  int
}

func newReverseProxy(targetUrlString string, port int) (*ReverseProxy, error) {
	SugaredLog.Debugf("Create new reverse proxy for target %s", targetUrlString)

	SugaredLog.Debugf("Parse target URL %s", targetUrlString)
	targetUrl, urlErr := url.Parse(targetUrlString)
	if urlErr != nil {
		SugaredLog.Errorf("Parse downstream url %s failed: %s", targetUrl, urlErr.Error())
		return nil, urlErr
	}

	Log.Debug("Create reverse proxy")
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	Log.Debug("Create director")
	director := proxy.Director

	Log.Debug("Set director on reverse proxy")
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Header.Set(additionalHeaderKey, req.Header.Get(hostHeaderKey))
		req.Host = req.URL.Host
	}

	Log.Debug("Set modify response on reverse proxy")
	proxy.ModifyResponse = func(res *http.Response) error {
		responseContent := map[string]interface{}{}
		err := parseResponse(res, &responseContent)
		if err != nil {
			return err
		}

		return captureMetrics(responseContent)
	}

	return &ReverseProxy{
		proxy: proxy,
		port:  port,
	}, nil
}

func (r *ReverseProxy) start() {
	Log.Info("Start reverse proxy")

	http.HandleFunc(rootEndpoint, r.proxy.ServeHTTP)
	SugaredLog.Infof("Listening on port %d", r.port)
	SugaredLog.Errorf(http.ListenAndServe(":"+strconv.Itoa(r.port), nil))
}

func (r *ReverseProxy) shutdown(timeout int) {
	Log.Warn("Shutdown reverse proxy")

	// no-op
}

func parseResponse(res *http.Response, unmarshalStruct interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	closeErr := res.Body.Close()
	if closeErr != nil {
		SugaredLog.Errorf("Close response body failed: %s", closeErr.Error())
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return json.Unmarshal(body, unmarshalStruct)
}

// Add your metrics capture code here
func captureMetrics(m map[string]interface{}) error {
	SugaredLog.Infof("captureMetrics = %+v\n", m)
	return nil
}
