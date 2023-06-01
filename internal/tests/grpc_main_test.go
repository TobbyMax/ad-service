package tests

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGRPCSuite(t *testing.T) {
	suite.Run(t, new(GRPCSuite))
}
