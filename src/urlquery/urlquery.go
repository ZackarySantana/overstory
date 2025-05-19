package urlquery

import (
	"net/url"
)

func AddQueryParam(rawURL, key, v string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set(key, v)
	u.RawQuery = q.Encode()

	return u.String(), nil
}
