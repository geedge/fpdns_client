// dns utils
package utils

import (
	"github.com/miekg/dns"
)

// build the question section of a dns packet
func QuestionMaker(domain string, qclass uint16, qtype uint16) *dns.Question {
	return &dns.Question{Name: dns.Fqdn(domain), Qtype: qtype, Qclass: qclass}
}

// build a specific query message
func QueryMaker(domain string, rd bool, qclass uint16, qtype uint16, edns bool) *dns.Msg {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = rd
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = *QuestionMaker(domain, qclass, qtype)
	if edns {
		msg = msg.SetEdns0(4096, false)
	}
	return msg
}

// query and receive the response
func DNSQuery(addr string, domain string, rd bool, qclass uint16, qtype uint16, edns bool) (*dns.Msg, error) {
	msg := QueryMaker(domain, rd, qclass, qtype, edns)
	res, err := dns.Exchange(msg, addr)
	return res, err
}
