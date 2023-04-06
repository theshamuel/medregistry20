//go:build unit
package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
	message string
}

func TestClientTestSuite(t *testing.T)  {
	suite.Run(t, new(ClientTestSuite))
}
func (s *ClientTestSuite) TestClientSetAgePostfix()  {
	tbl := []struct {
		c   Client
		res string
	}{
		{Client{Age: 1}, "год"},
		{Client{Age: 2}, "года"},
		{Client{Age: 4}, "года"},
		{Client{Age: 10}, "лет"},
		{Client{Age: 20}, "лет"},
		{Client{Age: 30}, "лет"},
		{Client{Age: 21}, "год"},
		{Client{Age: 32}, "года"},
		{Client{Age: 37}, "лет"},
		{Client{Age: 44}, "года"},
	}
	for i, tt := range tbl {
		tt.c.SetAgePostfix()
		assert.Equal(s.T(), tt.res, tt.c.AgePostfix, "test case #%d", i)
	}
}