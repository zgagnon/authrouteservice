package proxy_test

import (
	"bytes"
	. "github.com/zgagnon/authrouteservice/proxy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Secured Proxy", func() {
	var wasCalled = Evidence{wasCalled: false}
	roundtripper := &FakeTripper{evidence: &wasCalled}
	logger := &LoggingRoundTripper{Transporter: roundtripper, Okta: "http://callback.com"}

	Context("when the request has a session header", func() {
		headers := make(map[string][]string)
		headers["session_token"] = []string{"a secured session"}

		request := httptest.NewRequest("GET", "http://test.com", bytes.NewBuffer([]byte{}))
		request.Header = headers

		It("should succeed", func() {
			_ = "breakpoint"
			logger.RoundTrip(request)
			Expect(wasCalled.wasCalled).To(BeTrue())
		})
	})

	Context("when the request does not have a session header", func() {

		request := httptest.NewRequest("GET", "http://test.com", bytes.NewBuffer([]byte{}))

		It("should return a 302 response", func() {
			response, err := logger.RoundTrip(request)
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(302))
		})

		It("should contain the okta address in the Location header", func() {
			response, _ := logger.RoundTrip(request)
			Expect(response.Header["location"][0]).To(Equal("http://callback.com"))
		})
	})
})

type Evidence struct {
	wasCalled bool
}

type FakeTripper struct {
	evidence *Evidence
}

func (ft FakeTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ft.evidence.wasCalled = true
	return &http.Response{}, nil
}
