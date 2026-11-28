// GitHub provider credentials are opaque bearer values, distinct from workflow OIDC identity tokens.
package main

import "strings"

// validProviderToken bounds header resources and excludes injection/whitespace without interpreting
// provider prefixes, JWT claims, signatures or the old installation-token length. GitHub authenticates
// the credential. The caller still enforces reviewed source/repository/permission/attempt fences.
func validProviderToken(value string) bool {
	if len(value) < 32 || len(value) > 4096 {
		return false
	}
	for _, ch := range []byte(value) {
		if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || strings.ContainsRune("._~+/=-", rune(ch)) {
			continue
		}
		return false
	}
	return true
}
