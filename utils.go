package zeroconf

import "strings"

func parseSubtypes(service string) (string, []string) {
	subtypes := strings.Split(service, ",")
	return subtypes[0], subtypes[1:]
}

// trimDot is used to trim the dots from the start or end of a string
func trimDot(s string) string {
	return strings.Trim(s, ".")
}

// qualifyHostName returns hostName as a fully qualified name inside domain.
//
// Callers supply domain both as "local" and as "local.", so both sides are
// compared with their dots stripped: comparing a trimmed host against an
// untrimmed domain never matches, which is how a macOS host that already ends
// in ".local" used to acquire a second ".local" and produce an SRV target that
// resolves nowhere.
//
// The trailing dot on the result is mandatory -- miekg/dns refuses to pack a
// name that is not fully qualified.
func qualifyHostName(hostName, domain string) string {
	domain = trimDot(domain)
	hostName = trimDot(hostName)
	if hostName != domain && !strings.HasSuffix(hostName, "."+domain) {
		hostName += "." + domain
	}
	return hostName + "."
}
