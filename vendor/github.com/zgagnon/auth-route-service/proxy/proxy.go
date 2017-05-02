package proxy

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	DEFAULT_PORT              = "8080"
	CF_FORWARDED_URL_HEADER   = "X-Cf-Forwarded-Url"
	CF_PROXY_SIGNATURE_HEADER = "X-Cf-Proxy-Signature"
)

func NewProxy(transport http.RoundTripper, skipSslValidation bool) http.Handler {
	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			forwardedURL := req.Header.Get(CF_FORWARDED_URL_HEADER)
			sigHeader := req.Header.Get(CF_PROXY_SIGNATURE_HEADER)

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatalln(err.Error())
			}
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			logRequest(forwardedURL, sigHeader, string(body), req.Header, skipSslValidation)

			err = sleep()
			if err != nil {
				log.Fatalln(err.Error())
			}

			// Note that url.Parse is decoding any url-encoded characters.
			url, err := url.Parse(forwardedURL)
			if err != nil {
				log.Fatalln(err.Error())
			}

			req.URL = url
			req.Host = url.Host
		},
		Transport: transport,
	}
	return reverseProxy
}

func logRequest(forwardedURL, sigHeader, body string, headers http.Header, skipSslValidation bool) {
	log.Printf("Skip ssl validation set to %t", skipSslValidation)
	log.Println("Received request: ")
	log.Printf("%s: %s\n", CF_FORWARDED_URL_HEADER, forwardedURL)
	log.Printf("%s: %s\n", CF_PROXY_SIGNATURE_HEADER, sigHeader)
	log.Println("")
	log.Printf("Headers: %#v\n", headers)
	log.Println("")
	log.Printf("Request Body: %s\n", body)
}

func sleep() error {
	sleepMilliString := os.Getenv("ROUTE_SERVICE_SLEEP_MILLI")
	if sleepMilliString != "" {
		sleepMilli, err := strconv.ParseInt(sleepMilliString, 0, 64)
		if err != nil {
			return err
		}

		log.Printf("Sleeping for %d milliseconds\n", sleepMilli)
		time.Sleep(time.Duration(sleepMilli) * time.Millisecond)

	}
	return nil
}
