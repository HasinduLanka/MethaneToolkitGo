package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var wsroot string = "workspace/"
var wscache string = "workspace/.cache/"

func main() {

	MakeDir(wsroot)
	MakeDir(wscache)

	Print("-------------------------------------")
	Print("          Methane Toolkit            ")
	Print("     By github.com/HasinduLanka      ")
	Print("-------------------------------------")
	Print("")

	DeleteFiles("output.html")

	SetDefaultHTTPClientHostCookies(map[string]map[string]string{"https://quiz.kits2022.com": {"PHPSESSID": "djv6lcpou930f5m8vtkn89pts9"}})

	SA := NewStringAssembler()

	SA.AddStatic("https://quiz.kits2022.com/account.php?q=1,p=")
	SA.Add(NewReusableLineProvider_IntRange(1, 4, 0))

	URIListParser("", SA.Assemble().GetInstance(), false, nil, "output.html")

}

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

/*

Javascript code to generate cookie map:

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
*/
