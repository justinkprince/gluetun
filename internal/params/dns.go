package params

import (
	"fmt"
	"net"
	"strings"
	"time"

	dns "github.com/qdm12/dns/pkg/unbound"
	libparams "github.com/qdm12/golibs/params"
)

// GetDNSOverTLS obtains if the DNS over TLS should be enabled
// from the environment variable DOT.
func (r *reader) GetDNSOverTLS() (DNSOverTLS bool, err error) { //nolint:gocritic
	return r.env.OnOff("DOT", libparams.Default("on"))
}

// GetDNSOverTLSProviders obtains the DNS over TLS providers to use
// from the environment variable DOT_PROVIDERS.
func (r *reader) GetDNSOverTLSProviders() (providers []string, err error) {
	s, err := r.env.Get("DOT_PROVIDERS", libparams.Default("cloudflare"))
	if err != nil {
		return nil, err
	}
	for _, provider := range strings.Split(s, ",") {
		_, ok := dns.GetProviderData(provider)
		if !ok {
			return nil, fmt.Errorf("DNS over TLS provider %q is not valid", provider)
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

// GetDNSOverTLSVerbosity obtains the verbosity level to use for Unbound
// from the environment variable DOT_VERBOSITY.
func (r *reader) GetDNSOverTLSVerbosity() (verbosityLevel uint8, err error) {
	n, err := r.env.IntRange("DOT_VERBOSITY", 0, 5, libparams.Default("1"))
	return uint8(n), err
}

// GetDNSOverTLSVerbosityDetails obtains the log level to use for Unbound
// from the environment variable DOT_VERBOSITY_DETAILS.
func (r *reader) GetDNSOverTLSVerbosityDetails() (verbosityDetailsLevel uint8, err error) {
	n, err := r.env.IntRange("DOT_VERBOSITY_DETAILS", 0, 4, libparams.Default("0"))
	return uint8(n), err
}

// GetDNSOverTLSValidationLogLevel obtains the log level to use for Unbound DOT validation
// from the environment variable DOT_VALIDATION_LOGLEVEL.
func (r *reader) GetDNSOverTLSValidationLogLevel() (validationLogLevel uint8, err error) {
	n, err := r.env.IntRange("DOT_VALIDATION_LOGLEVEL", 0, 2, libparams.Default("0"))
	return uint8(n), err
}

// GetDNSMaliciousBlocking obtains if malicious hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_MALICIOUS.
func (r *reader) GetDNSMaliciousBlocking() (blocking bool, err error) {
	return r.env.OnOff("BLOCK_MALICIOUS", libparams.Default("on"))
}

// GetDNSSurveillanceBlocking obtains if surveillance hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_SURVEILLANCE
// and BLOCK_NSA for retrocompatibility.
func (r *reader) GetDNSSurveillanceBlocking() (blocking bool, err error) {
	// Retro-compatibility
	s, err := r.env.Get("BLOCK_NSA")
	if err != nil {
		return false, err
	} else if len(s) != 0 {
		r.logger.Warn("You are using the old environment variable BLOCK_NSA, please consider changing it to BLOCK_SURVEILLANCE") //nolint:lll
		return r.env.OnOff("BLOCK_NSA", libparams.Compulsory())
	}
	return r.env.OnOff("BLOCK_SURVEILLANCE", libparams.Default("off"))
}

// GetDNSAdsBlocking obtains if ads hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_ADS.
func (r *reader) GetDNSAdsBlocking() (blocking bool, err error) {
	return r.env.OnOff("BLOCK_ADS", libparams.Default("off"))
}

// GetDNSUnblockedHostnames obtains a list of hostnames to unblock from block lists
// from the comma separated list for the environment variable UNBLOCK.
func (r *reader) GetDNSUnblockedHostnames() (hostnames []string, err error) {
	s, err := r.env.Get("UNBLOCK")
	if err != nil {
		return nil, err
	} else if len(s) == 0 {
		return nil, nil
	}
	hostnames = strings.Split(s, ",")
	for _, hostname := range hostnames {
		if !r.regex.MatchHostname(hostname) {
			return nil, fmt.Errorf("hostname %q does not seem valid", hostname)
		}
	}
	return hostnames, nil
}

// GetDNSOverTLSCaching obtains if Unbound caching should be enable or not
// from the environment variable DOT_CACHING.
func (r *reader) GetDNSOverTLSCaching() (caching bool, err error) {
	return r.env.OnOff("DOT_CACHING", libparams.Default("on"))
}

// GetDNSOverTLSPrivateAddresses obtains if Unbound caching should be enable or not
// from the environment variable DOT_PRIVATE_ADDRESS.
func (r *reader) GetDNSOverTLSPrivateAddresses() (privateAddresses []string, err error) {
	s, err := r.env.Get("DOT_PRIVATE_ADDRESS")
	if err != nil {
		return nil, err
	} else if len(s) == 0 {
		return nil, nil
	}
	privateAddresses = strings.Split(s, ",")
	for _, address := range privateAddresses {
		ip := net.ParseIP(address)
		_, _, err := net.ParseCIDR(address)
		if ip == nil && err != nil {
			return nil, fmt.Errorf("private address %q is not a valid IP or CIDR range", address)
		}
	}
	return privateAddresses, nil
}

// GetDNSOverTLSIPv6 obtains if Unbound should resolve ipv6 addresses using
// ipv6 DNS over TLS from the environment variable DOT_IPV6.
func (r *reader) GetDNSOverTLSIPv6() (ipv6 bool, err error) {
	return r.env.OnOff("DOT_IPV6", libparams.Default("off"))
}

// GetDNSUpdatePeriod obtains the period to use to update the block lists and cryptographic files
// and restart Unbound from the environment variable DNS_UPDATE_PERIOD.
func (r *reader) GetDNSUpdatePeriod() (period time.Duration, err error) {
	s, err := r.env.Get("DNS_UPDATE_PERIOD", libparams.Default("24h"))
	if err != nil {
		return period, err
	}
	return time.ParseDuration(s)
}

// GetDNSPlaintext obtains the plaintext DNS address to use if DNS over TLS is disabled
// from the environment variable DNS_PLAINTEXT_ADDRESS.
func (r *reader) GetDNSPlaintext() (ip net.IP, err error) {
	s, err := r.env.Get("DNS_PLAINTEXT_ADDRESS", libparams.Default("1.1.1.1"))
	if err != nil {
		return nil, err
	}
	ip = net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("DNS plaintext address %q is not a valid IP address", s)
	}
	return ip, nil
}

// GetDNSKeepNameserver obtains if the nameserver present in /etc/resolv.conf
// should be kept instead of overridden, from the environment variable DNS_KEEP_NAMESERVER.
func (r *reader) GetDNSKeepNameserver() (on bool, err error) {
	return r.env.OnOff("DNS_KEEP_NAMESERVER", libparams.Default("off"))
}
