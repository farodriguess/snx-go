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
	// reqClone.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0")
	// reqClone.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	// reqClone.Header.Add("Accept-Language", "en-US,pt-BR;q=0.7,en;q=0.3")
	// reqClone.Header.Add("DNT", "1")
	// reqClone.Header.Add("Upgrade-Insecure-Requests", "1")
	// reqClone.Header.Add("Origin", "https://remoteaccess.alelo.com.br")
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
