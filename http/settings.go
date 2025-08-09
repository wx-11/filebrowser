package http

import (
	"encoding/json"
	"net/http"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/filebrowser/filebrowser/v2/settings"
)

type settingsData struct {
	Signup                bool                  `json:"signup"`
	CreateUserDir         bool                  `json:"createUserDir"`
	MinimumPasswordLength uint                  `json:"minimumPasswordLength"`
	UserHomeBasePath      string                `json:"userHomeBasePath"`
	Defaults              settings.UserDefaults `json:"defaults"`
	Rules                 []rules.Rule          `json:"rules"`
	Branding              settings.Branding     `json:"branding"`
	Tus                   settings.Tus          `json:"tus"`
	Shell                 []string              `json:"shell"`
	Commands              map[string][]string   `json:"commands"`
	MimeTypes             map[string]string     `json:"mimeTypes"`
	FilenameMimeTypes     map[string]string     `json:"filenameMimes"`
}

var settingsGetHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	data := &settingsData{
		Signup:                d.settings.Signup,
		CreateUserDir:         d.settings.CreateUserDir,
		MinimumPasswordLength: d.settings.MinimumPasswordLength,
		UserHomeBasePath:      d.settings.UserHomeBasePath,
		Defaults:              d.settings.Defaults,
		Rules:                 d.settings.Rules,
		Branding:              d.settings.Branding,
		Tus:                   d.settings.Tus,
		Shell:                 d.settings.Shell,
		Commands:              d.settings.Commands,
		MimeTypes:             d.settings.MimeTypes,
		FilenameMimeTypes:     d.settings.FilenameMimeTypes,
	}

	return renderJSON(w, r, data)
})

var settingsPutHandler = withAdmin(func(_ http.ResponseWriter, r *http.Request, d *data) (int, error) {
	req := &settingsData{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return http.StatusBadRequest, err
	}

	d.settings.Signup = req.Signup
	d.settings.CreateUserDir = req.CreateUserDir
	d.settings.MinimumPasswordLength = req.MinimumPasswordLength
	d.settings.UserHomeBasePath = req.UserHomeBasePath
	d.settings.Defaults = req.Defaults
	d.settings.Rules = req.Rules
	d.settings.Branding = req.Branding
	d.settings.Tus = req.Tus
	d.settings.Shell = req.Shell
	d.settings.Commands = req.Commands
	d.settings.MimeTypes = req.MimeTypes
	d.settings.FilenameMimeTypes = req.FilenameMimeTypes

	err = d.store.Settings.Save(d.settings)
	if err != nil {
		return errToStatus(err), err
	}

	// Re-initialize MIME types after saving settings
	files.InitializeMimeTypes(d.settings.MimeTypes)
	
	return http.StatusOK, nil
})
