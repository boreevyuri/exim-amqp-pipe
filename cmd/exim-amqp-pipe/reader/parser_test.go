package reader

import (
	"net/mail"
	"testing"
)

type testpair struct {
	headers      mail.Header
	resultString string
}

var testHeaders1 = mail.Header{
	"To":          {"user@mail.com"},
	"Envelope-To": {"user@mail.com"},
	"Cc":          {"<user2@mail.com>, user3@mail.com"},
}

var testHeaders2 = mail.Header{
	"To": {"user@mail.com"},
}

var testHeaders3 = mail.Header{
	"To": {"user@mail.com, user2@mail.com"},
}

var testHeaders4 = mail.Header{
	"To":          {"user@mail.com, user2@mail.com"},
	"Envelope-To": {"user@mail.com"},
}

var testHeaders5 = mail.Header{
	"To":          {"user@mail.com"},
	"Envelope-To": {"user@mail.com"},
}

var headersForTest = []testpair{
	{testHeaders1, "user@mail.com, user2@mail.com, user3@mail.com"},
	{testHeaders2, "user@mail.com"},
	{testHeaders3, "user@mail.com, user2@mail.com"},
	{testHeaders4, "user@mail.com, user2@mail.com"},
	{testHeaders5, "user@mail.com"},
}

func TestGetRecipients(t *testing.T) {
	for _, pair := range headersForTest {
		v := GetRecipients(&pair.headers)
		if v != pair.resultString {
			t.Error(
				"For", pair.headers, "\n",
				"expected", pair.resultString, "\n",
				"got", v,
			)
		}
	}
}
