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

	log.Printf("Forwarding to: %s\n", request.URL.String())
	res, err = lrt.Transporter.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("")
	log.Printf("Response Headers: %#v\n", res.Header)
	log.Println("")
	log.Printf("Response Body: %s\n", string(body))
	log.Println("")
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	log.Println("Sending response to GoRouter...")

	return res, err
}
