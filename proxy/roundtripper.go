package proxy

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

type LoggingRoundTripper struct {
	Transporter http.RoundTripper
}

func NewLoggingRoundTripper(skipSslValidation bool) *LoggingRoundTripper {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSslValidation},
	}
	return &LoggingRoundTripper{
		Transporter: tr,
	}
}

func (lrt *LoggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var err error
	var res *http.Response
	if len(request.Header["session_token"]) != 0 {

		log.Printf("Forwarding to: %s\n", request.URL.String())
		res, err = lrt.Transporter.RoundTrip(request)
		if err != nil {
			return nil, err
		}

		log.Println("")
		log.Printf("Response Headers: %#v\n", res.Header)
		log.Println("")
		res.Body = request.Body

		log.Println("Sending response to GoRouter...")

	} else {
		res = &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			StatusCode: 302,
		}

	}
	return res, err
}
