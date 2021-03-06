package params

import (
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	libparams "github.com/qdm12/golibs/params"
)

// GetCyberghostGroup obtains the server group for the Cyberghost server from the
// environment variable CYBERGHOST_GROUP.
func (r *reader) GetCyberghostGroup() (group string, err error) {
	s, err := r.env.Inside("CYBERGHOST_GROUP",
		constants.CyberghostGroupChoices(), libparams.Default("Premium UDP Europe"))
	return s, err
}

// GetCyberghostRegions obtains the country names for the Cyberghost servers from the
// environment variable REGION.
func (r *reader) GetCyberghostRegions() (regions []string, err error) {
	return r.env.CSVInside("REGION", constants.CyberghostRegionChoices())
}

// GetCyberghostClientKey obtains the client key to use for openvpn
// from the secret file /run/secrets/openvpn_clientkey or from the file
// /gluetun/client.key.
func (r *reader) GetCyberghostClientKey() (clientKey string, err error) {
	b, err := r.getFromFileOrSecretFile("OPENVPN_CLIENTKEY", string(constants.ClientKey))
	if err != nil {
		return "", err
	}
	return extractClientKey(b)
}

func extractClientKey(b []byte) (key string, err error) {
	pemBlock, _ := pem.Decode(b)
	if pemBlock == nil {
		return "", fmt.Errorf("cannot decode PEM block from client key")
	}
	parsedBytes := pem.EncodeToMemory(pemBlock)
	s := string(parsedBytes)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimPrefix(s, "-----BEGIN PRIVATE KEY-----")
	s = strings.TrimSuffix(s, "-----END PRIVATE KEY-----")
	return s, nil
}

// GetCyberghostClientCertificate obtains the client certificate to use for openvpn
// from the secret file /run/secrets/openvpn_clientcrt or from the file
// /gluetun/client.crt.
func (r *reader) GetCyberghostClientCertificate() (clientCertificate string, err error) {
	b, err := r.getFromFileOrSecretFile("OPENVPN_CLIENTCRT", string(constants.ClientCertificate))
	if err != nil {
		return "", err
	}
	return extractClientCertificate(b)
}

func extractClientCertificate(b []byte) (certificate string, err error) {
	pemBlock, _ := pem.Decode(b)
	if pemBlock == nil {
		return "", fmt.Errorf("cannot decode PEM block from client certificate")
	}
	parsedBytes := pem.EncodeToMemory(pemBlock)
	s := string(parsedBytes)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimPrefix(s, "-----BEGIN CERTIFICATE-----")
	s = strings.TrimSuffix(s, "-----END CERTIFICATE-----")
	return s, nil
}
