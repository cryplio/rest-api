package stdlogger

import "github.com/Nivl/go-logger"

// Makes sure Creator is a logger.Creator
var _ logger.Creator = (*Creator)(nil)

// NewCreator returns a logger creator that will use the provided token
// to create a new le driver for each single logger
func NewCreator() *Creator {
	return &Creator{}
}

// Creator creates new logger
type Creator struct{}

// New returns a new le client
func (c *Creator) New() (logger.Logger, error) {
	return New(), nil
}
