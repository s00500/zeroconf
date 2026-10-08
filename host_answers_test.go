package zeroconf

import (
	"net"
	"testing"

	"github.com/miekg/dns"
)

func hostAnswerServer(enabled bool) *Server {
	entry := NewServiceEntry("dev", "_skaarhoj._tcp", "local.")
	entry.HostName = "454143.SK_BLUEPILL.skaarhoj.local."
	entry.AddrIPv4 = []net.IP{net.ParseIP("192.168.0.50").To4()}
	entry.AddrIPv6 = []net.IP{net.ParseIP("fd1f::1")}
	s := &Server{service: entry, ttl: 3200}
	s.AnswerHostQueries(enabled)
	return s
}

func ask(s *Server, name string, qtype uint16) *dns.Msg {
	query := new(dns.Msg)
	query.SetQuestion(name, qtype)
	resp := &dns.Msg{}
	s.handleQuestion(query.Question[0], resp, query, 0)
	return resp
}

func TestHostQueriesDisabled(t *testing.T) {
	s := hostAnswerServer(false)
	if resp := ask(s, "50.0.168.192.in-addr.arpa.", dns.TypePTR); len(resp.Answer) != 0 {
		t.Fatalf("answered reverse query while disabled: %v", resp.Answer)
	}
	if resp := ask(s, "454143.SK_BLUEPILL.skaarhoj.local.", dns.TypeA); len(resp.Answer) != 0 {
		t.Fatalf("answered host query while disabled: %v", resp.Answer)
	}
}

func TestReversePTR(t *testing.T) {
	s := hostAnswerServer(true)

	resp := ask(s, "50.0.168.192.in-addr.arpa.", dns.TypePTR)
	if len(resp.Answer) != 1 {
		t.Fatalf("want 1 answer, got %v", resp.Answer)
	}
	if ptr := resp.Answer[0].(*dns.PTR); ptr.Ptr != "454143.SK_BLUEPILL.skaarhoj.local." {
		t.Fatalf("wrong PTR target %q", ptr.Ptr)
	}
	if len(resp.Extra) == 0 {
		t.Fatal("want the host's addresses as additional records")
	}

	v6, _ := dns.ReverseAddr("fd1f::1")
	if resp := ask(s, v6, dns.TypePTR); len(resp.Answer) != 1 {
		t.Fatalf("want ip6.arpa answer, got %v", resp.Answer)
	}

	if resp := ask(s, "51.0.168.192.in-addr.arpa.", dns.TypePTR); len(resp.Answer) != 0 {
		t.Fatalf("answered for a foreign address: %v", resp.Answer)
	}
}

func TestHostAddressQueries(t *testing.T) {
	s := hostAnswerServer(true)

	resp := ask(s, "454143.sk_bluepill.skaarhoj.local.", dns.TypeA) // case-insensitive
	if len(resp.Answer) != 1 || resp.Answer[0].(*dns.A).A.String() != "192.168.0.50" {
		t.Fatalf("want one A record for 192.168.0.50, got %v", resp.Answer)
	}
	if resp := ask(s, "454143.SK_BLUEPILL.skaarhoj.local.", dns.TypeAAAA); len(resp.Answer) != 1 {
		t.Fatalf("want one AAAA record, got %v", resp.Answer)
	}
	if resp := ask(s, "454143.SK_BLUEPILL.skaarhoj.local.", dns.TypeANY); len(resp.Answer) != 2 {
		t.Fatalf("want A and AAAA for ANY, got %v", resp.Answer)
	}
}
