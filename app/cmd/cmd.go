package cmd

import (
	"strings"
)

type CommonOptionsCommander interface {
	SetCommon(commonOpts CommonOptions)
	Execute(args []string) error
}

type CommonOptions struct {
	MedregAPIV1URL string
	MongoURL       string
	ReportsPath    string
}

func (c *CommonOptions) SetCommon(commonOpts CommonOptions) {
	c.MedregAPIV1URL = strings.TrimSuffix(commonOpts.MedregAPIV1URL, "/")
	c.ReportsPath = strings.TrimSuffix(commonOpts.ReportsPath, "/")
	c.MongoURL = strings.TrimSuffix(commonOpts.MongoURL, "/")
}
