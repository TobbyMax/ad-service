package tests

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestHTTPSuite(t *testing.T) {
	suite.Run(t, new(HTTPSuite))
}
