package proxy_test

import (
	"bytes"
	. "github.com/zgagnon/auth-route-service/proxy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Secured Proxy", func() {
	roundtripper := FakeTripper{wasCalled: false}
	logger := LoggingRoundTripper{Transporter: roundtripper}
	Context("when the request has a session header", func() {
		headers := make(map[string][]string)
		headers["session_token"] = []string{"a secured session"}

		request := httptest.NewRequest("GET", "http://test.com", bytes.NewBuffer([]byte{}))
		request.Header = headers

		It("should succeed", func() {
			_ = "breakpoint"
			logger.RoundTrip(request)
			Expect(roundtripper.wasCalled).To(BeTrue())
		})
	})
})

type FakeTripper struct {
	wasCalled bool
}

func (ft FakeTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}
