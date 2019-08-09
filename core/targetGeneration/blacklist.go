package targetgeneration

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"

	"github.com/nray-scanner/nray/utils"
	"github.com/zmap/go-iptree/blacklist"
)

// NrayBlacklist allows to add/query ip/net/dns blacklisted items
type NrayBlacklist struct {
	ipBlacklist  *blacklist.Blacklist
	dnsBlacklist *map[string]bool // value type not relevant, taking bool..
}

// NewBlacklist returns a new blacklist
func NewBlacklist() *NrayBlacklist {
	_dnsblacklist := make(map[string]bool)
	return &NrayBlacklist{
		ipBlacklist:  blacklist.New(),
		dnsBlacklist: &_dnsblacklist,
	}
}

// AddToBlacklist can be used if the type of the element
// is unclear
func (blacklist *NrayBlacklist) AddToBlacklist(element string) uint64 {
	if ipv4NetRegexpr.MatchString(element) { // An IPv4 network
		blacklist.AddNetToBlacklist(element)
		_, ipnet, err := net.ParseCIDR(element)
		utils.CheckError(err, true)
		return cidr.AddressCount(ipnet)
	} else if ipv4Regexpr.MatchString(element) { // An IPv4 address
		blacklist.AddNetToBlacklist(fmt.Sprintf("%s/32", element))
		return 1
	} else if mayBeFQDN(element) { // Probably a FQDN
		blacklist.AddDNSNameToBlacklist(element)
		return 1
	} else {
		// Don't care as target generation won't add anything not matching
		// the criteria above
		return 0
	}
}

// AddNetToBlacklist adds a CIDR network range to the blacklist
// <ip>/32 achieves the same for a single IP
func (blacklist *NrayBlacklist) AddNetToBlacklist(net string) {
	blacklist.ipBlacklist.AddEntry(net)
}

// AddDNSNameToBlacklist adds a FQDN to the blacklist
func (blacklist *NrayBlacklist) AddDNSNameToBlacklist(dnsName string) {
	(*blacklist.dnsBlacklist)[dnsName] = true
}

// IsIPBlacklisted returns true if the given IP is contained
// in a network in the blacklist
func (blacklist *NrayBlacklist) IsIPBlacklisted(ip string) bool {
	result, err := blacklist.ipBlacklist.IsBlacklisted(ip)
	utils.CheckError(err, false)
	return result
}

// IsDNSNameBlacklisted returns true if a given DNS name
// is blacklisted
func (blacklist *NrayBlacklist) IsDNSNameBlacklisted(dnsName string) bool {
	_, blacklisted := (*blacklist.dnsBlacklist)[dnsName]
	return blacklisted
}
