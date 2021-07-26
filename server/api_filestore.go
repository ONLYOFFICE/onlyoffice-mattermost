package main

import (
	"io"

	"github.com/mattermost/mattermost-server/v5/shared/filestore"
)

func getFilestore(p *Plugin) (filestore.FileBackend, error) {
	license := p.API.GetLicense()
	serverConfig := p.API.GetUnsanitizedConfig()
	filestore, err := filestore.NewFileBackend(serverConfig.FileSettings.ToFileBackendSettings(license != nil && *license.Features.Compliance))
	if err != nil {
		return nil, err
	}
	return filestore, nil
}

func (p *Plugin) WriteFile(fr io.Reader, path string) (int64, error) {
	filestore, err := getFilestore(p)
	if err != nil {
		return 0, err
	}

	result, err := filestore.WriteFile(fr, path)
	if err != nil {
		return result, err
	}
	return result, nil
}
