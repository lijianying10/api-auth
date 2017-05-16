package AuthKey_test

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/lijianying10/api-auth/AuthKey"
)

const rfc2822 = "Mon Jan 02 15:04:05 -0700 2006"

func TestSuccess(t *testing.T) {
	fmt.Println("rfc2822 NOW: ", time.Now().Format(rfc2822))
	authKey := AuthKey.AuthKey{
		AppKey:    "app_key_test",
		AppSecret: "secret_key_test",
	}
	method := "POST"
	date := time.Now().Format(rfc2822)
	resource := "test"
	mac := hmac.New(sha1.New, []byte(authKey.AppSecret))
	mac.Write([]byte(method + "\n\n\n" + date + "\n" + resource))
	res, err := authKey.CheckSignature(base64.StdEncoding.EncodeToString(mac.Sum(nil)), method, date, resource)
	if res != true || err != nil {
		t.FailNow()
	}
}
