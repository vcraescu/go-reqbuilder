package reqbuilder

import (
	"fmt"
	"net/http"
)

var (
	JSONContentHeader       = ContentHeader("application/json")
	URLEncodedContentHeader = ContentHeader("application/x-www-form-urlencoded")
	JSONAcceptHeader        = Header("Accept", "application/json")
)

func ContentHeader(v string) http.Header {
	return Header("Content-Type", v)
}

func Header(k, v string) http.Header {
	h := http.Header{}
	h.Set(k, v)

	return h
}

func AuthHeader(v string) http.Header {
	return Header("Authorization", v)
}

func AuthBearerHeader(token string) http.Header {
	return AuthHeader(fmt.Sprintf("Bearer %s", token))
}
