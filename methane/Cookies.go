package methane

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func SetDefaultHTTPClientHostCookies(cookies map[string]map[string]string) {
	for host, ck := range cookies {
		SetDefaultHTTPClientCookies(host, ck)
	}
}

func SetDefaultHTTPClientCookies(hosturl string, cookies map[string]string) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		CheckError(err)
	}

	urlObj, _ := url.Parse(hosturl)
	jar.SetCookies(urlObj, ParseCookies(cookies))

	client := &http.Client{Jar: jar}
	http.DefaultClient = client
}

func ParseCookies(cookies map[string]string) []*http.Cookie {
	o := make([]*http.Cookie, len(cookies))

	i := 0
	for k, v := range cookies {
		o[i] = &http.Cookie{
			Name:   k,
			Value:  v,
			MaxAge: 600000,
		}
		i++
	}

	return o
}

const JavascriptCodeToGenerateCookieMap = `
var str = document.cookie;
str = str.split('; ');
var cookies = {};
for (var i = 0; i < str.length; i++) {
    var cur = str[i].split('=');
    cookies[cur[0]] = cur[1];
}

var result = {};
result[document.location.origin] = cookies;

console.log("map[string]map[string]string" + JSON.stringify(result) + "");
`
