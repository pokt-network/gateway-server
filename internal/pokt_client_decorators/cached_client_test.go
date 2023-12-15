package pokt_client_decorators

// Basic imports
import (
	"os-gateway/mocks"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CachedClientTestSuite struct {
	suite.Suite
	mockPocketService *mocks.PocketService
}

func (suite *CachedClientTestSuite) SetupTest() {
	suite.mockPocketService = new(mocks.PocketService)

}

func TestCachedClientTestSuite(t *testing.T) {
	suite.Run(t, new(CachedClientTestSuite))
}
