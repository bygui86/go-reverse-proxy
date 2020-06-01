package reverse_proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/bygui86/go-reverse-proxy/reverse-proxy/logging"
)

func createSingleHostReverseProxy(targetUrlString string) (*httputil.ReverseProxy, error) {
	logging.SugaredLog.Debugf("Parse target URL %s", targetUrlString)
	targetUrl, urlErr := url.Parse(targetUrlString)
	if urlErr != nil {
		logging.SugaredLog.Errorf("Parse downstream url %s failed: %s", targetUrl, urlErr.Error())
		return nil, urlErr
	}

	logging.Log.Debug("Create reverse proxy")
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	logging.Log.Debug("Create director")
	director := proxy.Director

	logging.Log.Debug("Set director on reverse proxy")
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Header.Set(additionalHeaderKey, req.Header.Get(hostHeaderKey))
		req.Host = req.URL.Host
	}

	logging.Log.Debug("Set modify response on reverse proxy")
	proxy.ModifyResponse = func(res *http.Response) error {
		body, err := duplicateResponseBody(res)
		if err != nil {
			return err
		}

		return customBehaviour(body)
	}
	return proxy, nil
}

func (r *ReverseProxy) listenAndServe() {
	logging.SugaredLog.Debugf("Listen and server on port %d", r.port)
	r.errChannel <- http.ListenAndServe(":"+strconv.Itoa(r.port), nil)
}

func (r *ReverseProxy) startHttpServerController() {
	logging.Log.Debug("Start HTTP server controller")
	for err := range r.errChannel {
		logging.SugaredLog.Errorf("HTTP server failed and stopped working: %s", err.Error())
		r.running = false
		os.Exit(502)
	}
}

func duplicateResponseBody(res *http.Response) ([]byte, error) {
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	closeErr := res.Body.Close()
	if closeErr != nil {
		logging.SugaredLog.Errorf("Close response body failed: %s", closeErr.Error())
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

// TODO better understand usage
func dumpResponse(res *http.Response) ([]byte, error) {
	return httputil.DumpResponse(res, true)
}

// TODO add your custom behaviour here
func customBehaviour(responseBody []byte) error {
	logging.SugaredLog.Infof("Custom behaviour on response body: %s", string(responseBody))
	return nil
}
