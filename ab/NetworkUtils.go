package ab

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

// LocalIPAddress returns the local IP address of the machine
func LocalIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				return ip4.String()
			}
		}
	}

	return "127.0.0.1"
}

// ResponseContent gets the content of the response from a given URL
func ResponseContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ResponseContentURI gets the content of the response from a given URI
func ResponseContentURI(uri string) (string, error) {
	resp, err := http.Post(uri, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Spec constructs a URL spec string with the given parameters
func Spec(host, botid, custid, input string) (string, error) {
	var spec string
	inputEncoded := url.QueryEscape(input)

	if custid == "0" {
		spec = fmt.Sprintf("%s?botid=%s&input=%s", host, botid, inputEncoded)
	} else {
		spec = fmt.Sprintf("%s?botid=%s&custid=%s&input=%s", host, botid, custid, inputEncoded)
	}

	return spec, nil
}
