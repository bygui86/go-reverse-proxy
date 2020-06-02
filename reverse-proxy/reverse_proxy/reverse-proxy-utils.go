package reverse_proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/bygui86/go-reverse-proxy/reverse-proxy/logging"
)

func setupSingleHostReverseProxy(targetUrlString string) (*httputil.ReverseProxy, error) {
	logging.SugaredLog.Debugf("Setup new single-host reverse proxy to target %s", targetUrlString)

	logging.SugaredLog.Debugf("Parse target URL %s", targetUrlString)
	targetUrl, urlErr := url.Parse(targetUrlString)
	if urlErr != nil {
		logging.SugaredLog.Errorf("Parse downstream url %s failed: %s", targetUrl, urlErr.Error())
		return nil, urlErr
	}

	logging.Log.Debug("Create reverse proxy")
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	logging.Log.Debug("Setup proxy director")
	proxy.Director = createCustomDirector(proxy)

	logging.Log.Debug("Set proxy modify response")
	proxy.ModifyResponse = modifyResponse

	// INFO: if not using gorilla mux router with HTTP server, uncomment this line to directly access to reverse proxy
	// http.HandleFunc(rootEndpoint, r.proxy.ServeHTTP)

	return proxy, nil
}

func createCustomDirector(proxy *httputil.ReverseProxy) func(req *http.Request) {
	logging.Log.Debug("Get default director from proxy")
	defaultDirector := proxy.Director

	logging.Log.Debug("Set custom director on proxy")
	return func(req *http.Request) {
		logging.Log.Debug("Custom proxy director")

		logRequestInfo(req, "BEFORE default redirection")
		// BEFORE DEFAULT REDIRECTION
		// req.Host = localhost:8080
		// req.Header.Get(hostHeaderKey) = -empty-
		// req.URL.Host = -empty-
		// req.URL.Scheme = -empty-

		// WARN: this is required to perform a correct redirection (see httputil.NewSingleHostReverseProxy)
		defaultDirector(req)

		logRequestInfo(req, "AFTER default redirection")
		// AFTER DEFAULT REDIRECTION
		// req.Host = localhost:8080
		// req.Header.Get(hostHeaderKey) = -empty-
		// req.URL.Host = localhost:8081
		// req.URL.Scheme = http

		req.Header.Set(xForwardedHost, req.Host)
	}
}

func logRequestInfo(req *http.Request, prefix string) {
	logging.SugaredLog.Debugf("[%s] Host: req[%s], reqHeader[%s], reqUrl[%s]",
		prefix, req.Host, req.Header.Get(hostHeaderKey), req.URL.Host)
	logging.SugaredLog.Debugf("[%s] Scheme: [%s]", prefix, req.URL.Scheme)
}

func modifyResponse(res *http.Response) error {
	body, err := duplicateResponseBody(res)
	if err != nil {
		return err
	}
	return customBehaviour(body)
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

// INFO: add your custom behaviour here
func customBehaviour(responseBody []byte) error {
	logging.SugaredLog.Infof("Custom behaviour on response body: %s", string(responseBody))
	return nil
}
