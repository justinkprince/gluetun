package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/params"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

func (c *cli) OpenvpnConfig(os os.OS) error {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel)
	if err != nil {
		return err
	}
	paramsReader := params.NewReader(logger, os)
	allSettings, _, err := settings.GetAllSettings(paramsReader)
	if err != nil {
		return err
	}
	allServers, err := storage.New(logger, os, constants.ServersData).
		SyncServers(constants.GetAllServers())
	if err != nil {
		return err
	}
	providerConf := provider.New(allSettings.OpenVPN.Provider.Name, allServers, time.Now)
	connection, err := providerConf.GetOpenVPNConnection(allSettings.OpenVPN.Provider.ServerSelection)
	if err != nil {
		return err
	}
	lines := providerConf.BuildConf(connection, "nonroortuser", allSettings.OpenVPN)
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
