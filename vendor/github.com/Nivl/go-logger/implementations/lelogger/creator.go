package lelogger

import (
	"github.com/Nivl/go-logger"
	"github.com/bsphere/le_go"
)

// Makes sure Creator is a logger.Creator
var _ logger.Creator = (*Creator)(nil)

// NewCreator returns a logger creator that will use the provided token
// to create a new le driver for each single logger
func NewCreator(token string) *Creator {
	return &Creator{
		token: token,
	}
}

// NewSharedCreator returns a logger creator that will use the provided token
// to create a new le driver, and will reuse this driver for all logger created
func NewSharedCreator(token string) (*Creator, error) {
	driver, err := le_go.Connect(token)
	if err != nil {
		return nil, err
	}
	return &Creator{
		token:        token,
		shareDriver:  true,
		sharedDriver: driver,
	}, nil
}

// NewCreatorWithDriver returns a logger creator that will always use the
// provider driver
func NewCreatorWithDriver(driver *le_go.Logger) *Creator {
	return &Creator{
		shareDriver:  true,
		sharedDriver: driver,
	}
}

// Creator creates new logetries clients
type Creator struct {
	token        string
	shareDriver  bool
	sharedDriver *le_go.Logger
}

// New returns a new le client
func (c *Creator) New() (logger.Logger, error) {
	if c.shareDriver {
		return NewWithDriver(c.sharedDriver), nil
	}
	return New(c.token)
}
