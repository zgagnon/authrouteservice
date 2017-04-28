package main

import (
	"github.com/zgagnon/auth-route-service/proxy"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	DEFAULT_PORT = "8080"
)

func main() {
	var (
		port              string
		skipSslValidation bool
		err               error
	)

	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}
	if skipSslValidation, err = strconv.ParseBool(os.Getenv("SKIP_SSL_VALIDATION")); err != nil {
		skipSslValidation = true
	}
	log.SetOutput(os.Stdout)

	roundTripper := proxy.NewLoggingRoundTripper(skipSslValidation)
	proxy := proxy.NewProxy(roundTripper, skipSslValidation)

	log.Fatal(http.ListenAndServe(":"+port, proxy))
}
