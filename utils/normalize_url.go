package utils

import (
	"net"
	"net/url"
	"path"
	"sort"
	"strings"
)

func NormalizeURL(raw string) (string, error) {
	// 1. Trim spaces
	raw = strings.TrimSpace(raw)

	// 2. Add scheme if missing
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	// 3. Lowercase scheme and host
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// 4. Remove default ports
	host, port, err := net.SplitHostPort(u.Host)
	if err == nil {
		if (u.Scheme == "http" && port == "80") ||
			(u.Scheme == "https" && port == "443") {
			u.Host = host
		}
	}

	// 5. Clean path
	u.Path = path.Clean(u.Path)

	// path.Clean removes trailing slash except root
	if u.Path != "/" && strings.HasSuffix(raw, "/") {
		u.Path += "/"
	}

	// 6. Normalize query params (sort them)
	q := u.Query()
	for key := range q {
		sort.Strings(q[key])
	}
	u.RawQuery = q.Encode()

	// 7. Remove fragment
	u.Fragment = ""

	// 8. Ensure empty path becomes /
	if u.Path == "" {
		u.Path = "/"
	}

	return u.String(), nil
}
