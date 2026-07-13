// Package passkeys provides WebAuthn (Passkeys) registration, verification,
// authentication, and credential storage endpoints.
package passkeys

import "regexp"

var passkeyNameRegexp = regexp.MustCompile(`^[A-Za-z0-9_\s-]+$`)
