package util

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type HttpClientOptions struct {
	SkipSecurity bool
	Debug        bool
}

func CreateHttpClient(options *HttpClientOptions) (*http.Client, *cookiejar.Jar) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: options.SkipSecurity},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &http.Client{Transport: LoggerRoundTripper{tr, options.Debug}, Jar: jar}, jar

}

type LoggerRoundTripper struct {
	Proxied http.RoundTripper
	Debug   bool
}

func (lrt LoggerRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {

	lrt.log("Send Request...")
	// Send the request, get the response (or the error)
	reqBody := readBody(req.Body)
	reqClone := req.Clone(req.Context())
	reqClone.Body = io.NopCloser(strings.NewReader(reqBody))
	res, e = lrt.Proxied.RoundTrip(reqClone)

	text := `
---------------- request -----------------
%s %s
%s

%s
---------------- response ----------------
%s %s
%s


`

	// Handle the result.
	if e != nil {
		lrt.log(fmt.Sprintf("Error: %v", e))
	} else {
		lrt.log(fmt.Sprintf(text, req.Method, req.URL, printHeader(reqClone.Header), reqBody, res.Status, res.Request.URL, printHeader(res.Header)))
	}

	return
}

func (lrt LoggerRoundTripper) log(msg string) {
	if lrt.Debug {
		fmt.Println(msg)
	}
}

func printHeader(header http.Header) string {
	text := ""
	for k, v := range header {
		text = strings.Join([]string{text, fmt.Sprintf("%s: %s", k, v)}, "\n")
	}
	return text
}

func readBody(Body io.ReadCloser) string {
	if Body == nil {
		return ""
	}
	resBody, _ := ioutil.ReadAll(Body)
	return string(resBody)
}
