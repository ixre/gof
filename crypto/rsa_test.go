package crypto

import (
	"fmt"
	"testing"
)

func TestRsaToken(t *testing.T) {
	publicKey, privateKey, err := GenRsaKeys(2048)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("publicKey=", publicKey)
	fmt.Println("privateKey=", privateKey)
}

func TestParsePrivateKey(t *testing.T) {
	var privk = "-----BEGIN PRIVATE KEY-----\n" +
		"MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDF/hqHZZb7r0S5KuuQ1zE4v6BT+irjybOR0mIBbRqUnlBlIK8eayxs7eTazEn7FIFjepvGMxgH/2tC6R7s45KaoQo5Yq9l/rvziyYI50U4SZor1mV24nlCNLbx5BqDBFcGwxOJqwZGVTelBVjDtOsper10rUjhtwDFcLSe82VoPQUt8k9H4zw8+0lC4DsK0JlNtRJNAi380Fmz5JV19+12D2N8Tn9+pqFXzjyvp2EyJ/hS8uHUXZGy3lh7cbeEkFu5sFcKB2RDSs++8Y5vyeXQ6RLqMlEbJIRcRRAeMaCZ2Vn5OATYQKCvTPmITTzKB7NoOvEOC9FO4V6HMjidZzBTAgMBAAECggEADufwi10EnvI1FFO85GyvEfyrT2c4L2oSENpr8nuKUsIQf2yUgo/DCnhmkGps73A9xYWHkMZr+r4qDyGJ6H/Bm86f/G4HkoA5Gj7RoD35IiG4b7B2dxrZ0jgxxchMjqyW+LVbFTRBBq6Hv+7FHgbS5Y6OEOiy4ftrHXI8xvLAIbbEa9k1EVmH2ZvA5iVTBuZGWsEAQMRrIBNpmyB3Lnmo7iK28vpEPLvxADtlr3/1vpwfIPMb2fUYkuMXsCPuxjGxtkiCNhahUyzzwGG8rvszx/JcP/vWwRC7IQQff+YONdGKrJT5VqchJV1oaKbLg9CbU1/xsuLOn2RZP1A3/ssdsQKBgQDrlYhZ8BYSa2l5euKX7r4NFGETD8UGnyJmCGPy22VstJ77vAvffVLkKSzWrZgOlmW8MdRfFUsLfPaolLx56rCtdgS6mwSh4kqz9nKMuQjQbpECJAJtZL4FuMjVKSL/71Kew3/Bc/MNo6uKGxiK54KjxFu4TXWplKHFAI1MPuhdvQKBgQDXJpkFta6XwWbtrBCrgN5+eROA9qP+xC0WF/Ar8jbNJAntoUYXFLkIMt1HJFKAPND71x54G0ZHHpL7LJCP/NiGhY19/4S1oBP79d67HPku9Kbrm1NXKUzafOv2rPXSK7uGR+XSgnnKbs5GicipcqZP3+OGOajb9xxjer0IpU//TwKBgQCfHy8r4FhoNJjXbsMicCV6XCt9XodsA4yOclhgLwSAujcwPUGfwNx+M7mPf01XfQpWZSnW12EK72sDTwNHLdgMMczb5dzpIxnmGC4jEs/7SNM1KPFixkr7PmaYY+K6EAI0LkRafGDM86Hn9IlNOTYqO3TgNaGl2zixAcBuoYb92QKBgFsS7aerFrMKnWVydsQCkyx6WDU5MoZ/yI4XqAUSTPxdiw5aPG88yG6eCWk6COpb1CMnFrDE6uTkHlfQr4kkAQxAsHprlWPE1XDMzXHre9fSnG4TnB3DT9MVGlWbNZu4A3N+L90CekekzBCz9os0Cw64uXlyIvaqDgxWQnrMb6alAoGBAJ44E3SOo9DD5UOk+6swf/YplhqG2sayJruVib+1D2dlWu/+LxJqQZJGI/jtLVO24q7XGdnlA1YXA85DRI9/VUPPOEaLpUI91KWHUaN0Cgcin/O02UR+UWWvtbNEhI8Huk4BDGOPrBxz1tI2Bw1IvkD6u/mKmiExhzCUX/oAAesT" +

		"\n-----END PRIVATE KEY-----"

	_, err := ParsePrivateKey(privk)
	println("---", err)
}
