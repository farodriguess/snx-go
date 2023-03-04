package snx

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"snxgo/crypto"
	"snxgo/util"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type SNXParams struct {
	Host         string
	User         string
	Password     string
	Realm        string
	SkipSecurity bool
	Debug        bool
}

type SNXConnect struct {
	Params SNXParams
}

func (snx *SNXConnect) Connect() {
	client, jar := util.CreateHttpClient(&util.HttpClientOptions{SkipSecurity: snx.Params.SkipSecurity, Debug: snx.Params.Debug})

	host := snx.Params.Host
	if !strings.Contains(host, "http") {
		host = fmt.Sprintf("https://%s", host)
	}
	nextUrl, _ := url.Parse(host)

	res, err := client.Get(nextUrl.String())
	checkHttpError(err)
	resBody := snx.readBody(res.Body)

	rsaLocation := snx.getlocationRSA(resBody)

	loginAction := snx.getLoginAction(resBody)

	snx.log(fmt.Sprintf("RSA Location: %s", rsaLocation))
	snx.log(fmt.Sprintf("Login URL: %s", loginAction))

	res, err = client.Get(nextUrl.JoinPath(rsaLocation).String())
	checkHttpError(err)

	resBody = snx.readBody(res.Body)

	modulus, exponent := snx.parseRSAParams(resBody)

	snx.log(fmt.Sprintf("Modulus: %s , Exponent: %d\n", modulus, exponent))

	nextUrl = nextUrl.JoinPath(loginAction)

	pwEncode := crypto.PwEncode{Modulus: modulus, Exponent: exponent, Testing: false, Debug: snx.Params.Debug}

	encodedPWD := pwEncode.EncodePWD(snx.Params.Password)

	formData := url.Values{
		"selectedRealm": {snx.Params.Realm},
		"loginType":     {"Standard"},
		"userName":      {snx.Params.User},
		"vpid_prefix":   {""},
		"pin":           {""},
		"password":      {encodedPWD},
		"HeightData":    {""},
	}

	cookie := http.Cookie{Name: "selected_realm", Value: snx.Params.Realm, Secure: true}
	jar.SetCookies(nextUrl, []*http.Cookie{&cookie})

	fmt.Printf("Logging in Checkpoint VPN: %s. Please check your device to authorize 2FA...\n", snx.Params.Host)

	res, err = client.Post(nextUrl.String(), "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	checkHttpError(err)

	nextUrl = res.Request.URL

	if strings.Contains(nextUrl.Path, "Login/ActivateLogin") {
		query := nextUrl.Query()
		query.Add("ActivateLogin", "activate")
		query.Add("LangSelect", "en_US")
		query.Add("submit", "Continue")
		query.Add("HeightData", "")
		nextUrl.RawQuery = query.Encode()
		nextUrl.Host = regexp.MustCompile(`^https?://`).ReplaceAllString(snx.Params.Host, "")
		res, err = client.Get(nextUrl.String())
		checkHttpError(err)

		nextUrl = res.Request.URL
	}

	if strings.Contains(nextUrl.Path, "Portal/Main") {
		nextUrl.Path = "sslvpn/SNX/extender"
		res, err = client.Get(nextUrl.String())
		checkHttpError(err)

		extenderParams := snx.parseExtender(snx.readBody(res.Body))

		snxExtender := SNXExtender{Params: extenderParams, Debug: snx.Params.Debug}

		snxExtender.CallSNX()

	} else {
		body := snx.readBody(res.Body)
		errorText := snx.parseErrorMessage(body)
		if errorText != "" {
			fmt.Printf("Authentication Failure: %s\n", errorText)
		} else {
			fmt.Println("An error ocurred on connect. Run in debug mode to details...")
			snx.log(fmt.Sprintf("Response Body: %s\n", body))
		}
	}

}

func (s SNXConnect) readBody(Body io.ReadCloser) string {
	if Body == nil {
		return ""
	}
	resBody, err := ioutil.ReadAll(Body)
	checkError(err)
	return string(resBody)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func checkHttpError(err error) {
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		panic(err)
	}
}

func (s *SNXConnect) log(msg string) {
	if s.Params.Debug {
		fmt.Println(msg)
	}
}

func (snxConnect *SNXConnect) getlocationRSA(text string) (data string) {
	tkn := html.NewTokenizer(strings.NewReader(text))

	for {
		tt := tkn.Next()

		switch {
		case tt == html.ErrorToken:
			return ""
		case tt == html.StartTagToken:
			t := tkn.Token()
			if t.Data == "script" {
				for _, s := range t.Attr {
					snxConnect.log(fmt.Sprintf("src attr name: %s. attr value: %s", s.Key, s.Val))
					if s.Key == "src" && strings.Contains(s.Val, "RSA") {
						return s.Val
					}
				}
			}
		}
	}

}

func (snxConnect *SNXConnect) getLoginAction(text string) (data string) {
	tkn := html.NewTokenizer(strings.NewReader(text))

	for {
		tt := tkn.Next()

		switch {
		case tt == html.ErrorToken:
			return ""
		case tt == html.StartTagToken:
			t := tkn.Token()
			if t.Data == "form" {
				attrMap := make(map[string]string)
				for _, s := range t.Attr {
					attrMap[s.Key] = s.Val
				}
				if attrMap["id"] == "loginForm" && attrMap["method"] == "post" {
					return attrMap["action"]
				}
			}
		}
	}

}

func (snx *SNXConnect) parseRSAParams(body string) (modulus string, exponent int) {
	modulusRegex := regexp.MustCompile(`var modulus = '(?P<modulus>.+)'`)
	exponentRegex := regexp.MustCompile(`var exponent = '(?P<exponent>.+)'`)
	modulosNameIndex := modulusRegex.SubexpIndex("modulus")
	exponentNameIndex := exponentRegex.SubexpIndex("exponent")
	for _, line := range strings.Split(body, "\n") {
		matches := modulusRegex.FindStringSubmatch(line)
		if len(matches) >= modulosNameIndex {
			modulus = matches[modulosNameIndex]
		}
		matches = exponentRegex.FindStringSubmatch(line)
		if len(matches) >= exponentNameIndex {
			v, _ := strconv.ParseInt(matches[exponentNameIndex], 16, 0)
			exponent = int(v)
		}
	}

	return modulus, exponent
}

func (snx *SNXConnect) parseExtender(dat string) (params map[string]string) {

	params = map[string]string{}
	m1 := regexp.MustCompile(`.*\.`)
	m2 := regexp.MustCompile(`^ *"|" *$`)

	for _, line := range strings.Split(dat, "\n") {
		if strings.Contains(line, "/* Extender.user_name") {
			stmts := strings.Split(line, ";")
			for _, stmt := range stmts {
				if strings.Contains(stmt, "=") {
					hs := strings.SplitN(stmt, "=", 2)
					lhs := strings.Trim(m1.ReplaceAllString(hs[0], ""), " ")
					rhs := strings.Trim(m2.ReplaceAllString(hs[1], ""), " ")
					params[lhs] = rhs
				}
			}
		}
	}

	for k, v := range params {
		snx.log("Parsed Extender Params:")
		snx.log(fmt.Sprintf("key: \"%s\", value: \"%s\"\n", k, v))
	}

	return params
}

func (snx *SNXConnect) parseErrorMessage(dat string) string {
	tkn := html.NewTokenizer(strings.NewReader(dat))

	found := false

	for {
		tt := tkn.Next()

		switch {
		case tt == html.ErrorToken:
			return ""
		case tt == html.StartTagToken:
			t := tkn.Token()
			if t.Data == "span" {
				for _, s := range t.Attr {
					if s.Key == "class" && s.Val == "errorMessage" {
						found = true
					}
				}
			}
		case tt == html.TextToken:
			if found {
				return strings.TrimSpace(html.UnescapeString(string(tkn.Text())))
			}
		}
	}
}
