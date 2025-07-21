package apartment

import (
	"fmt"
	"net/url"
	"testing"
)


func TestXxx(t *testing.T) {
	u := url.URL{
		Scheme:      "http",
		Host:        "127.0.0.1:8080",
		Path:        "/api/v1/apartment/invite/accept",
		RawQuery:    fmt.Sprintf("%s=%s", "token", "xxx-xxx"),
	}
	fmt.Printf("----> %q\n", u.String())
}