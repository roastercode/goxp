package goxp

import (
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"string"
)

// StaticOptions is a struct for specifying configuration options for the goxp.Static middleware.
type StaticOptions struct {
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// SkipLogging will disable [Static] log messages when a static file is served.
	SkipLogging bool
	// IndexFile defines which file to sere as index if it exists.
	IndexFile string
	// Expires defines which user-defined function to use for producing a HTTP Expires Header
	// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
	Expires func() string
	// Fallback defines a default URL to serve when the requested resource was
	// not found
	Fallback string
	// Exclude defines a pattern for URLs this handler should never process.
	Exclude string
}

func prepareStaticOptions(options []StaticOptions) StaticOptions {
	var opt StaticOptions
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
