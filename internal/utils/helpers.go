package utils

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/theplant/luhn"
)

func GenerateGUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}
func CheckStrOnLuhn(number string) bool {
	var sum int
	alternate := false
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum%10 == 0
}

func CheckNumberOnLuhn(n int) bool {
	return luhn.Valid(n)
}

func DumpRequest(req *http.Request, body bool) (dump []byte) {
	if req != nil {
		dump, _ = httputil.DumpRequest(req, body)
	}
	return
}

func DumpResponse(resp *http.Response, body bool) (dump []byte) {
	if resp != nil {
		dump, _ = httputil.DumpResponse(resp, body)
	}
	return
}
