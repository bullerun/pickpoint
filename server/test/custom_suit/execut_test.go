package custom_suit

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCustomSuit(t *testing.T) {
	suite.Run(t, &ItemSuite{})
}
