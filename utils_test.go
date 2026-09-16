package zeroconf

import "testing"

func TestQualifyHostName(t *testing.T) {
	tests := []struct {
		name     string
		hostName string
		domain   string
		want     string
	}{
		// The regression: os.Hostname() on macOS already ends in ".local",
		// and Register passes Domain as "local." with a trailing dot.
		{"macOS host, dotted domain", "Skaarhojs-Mac-mini.local", "local.", "Skaarhojs-Mac-mini.local."},
		{"macOS host, bare domain", "Skaarhojs-Mac-mini.local", "local", "Skaarhojs-Mac-mini.local."},
		{"macOS host, already fqdn", "Skaarhojs-Mac-mini.local.", "local.", "Skaarhojs-Mac-mini.local."},

		// skaarOS device hostnames end in ".skaarhoj" and must still gain
		// the domain -- this is the behaviour that was already correct.
		{"skaarOS host, dotted domain", "446124.SK_COLORFLYV3.skaarhoj", "local.", "446124.SK_COLORFLYV3.skaarhoj.local."},
		{"skaarOS host, bare domain", "446124.SK_COLORFLYV3.skaarhoj", "local", "446124.SK_COLORFLYV3.skaarhoj.local."},

		// A label that merely ends in the domain string is not in the domain.
		{"suffix is not a label boundary", "printer.notlocal", "local", "printer.notlocal.local."},

		{"host equals domain", "local", "local.", "local."},
		{"non-local domain", "host.example.com", "example.com", "host.example.com."},
		{"non-local domain, needs append", "host", "example.com", "host.example.com."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := qualifyHostName(tt.hostName, tt.domain); got != tt.want {
				t.Errorf("qualifyHostName(%q, %q) = %q, want %q", tt.hostName, tt.domain, got, tt.want)
			}
		})
	}
}

// qualifyHostName must be idempotent: Register may run against an entry whose
// HostName was already qualified by an earlier call.
func TestQualifyHostNameIdempotent(t *testing.T) {
	for _, host := range []string{"Skaarhojs-Mac-mini.local", "446124.SK_COLORFLYV3.skaarhoj", "host"} {
		once := qualifyHostName(host, "local.")
		twice := qualifyHostName(once, "local.")
		if once != twice {
			t.Errorf("not idempotent for %q: %q then %q", host, once, twice)
		}
	}
}
