// Package faux is used exclusive for testing purposes.
package faux

import (
	"github.com/Sirupsen/logrus"

	"fmt"

	"github.com/modeneis/uminer/src/model"
)

var (
	log *logrus.Entry = logrus.WithField("package", "faux")
)

// Provider is used only for testing.
type Provider struct {
	flags *model.Flags
}

// Name is used only for testing.
func (p Provider) Name() string {
	return "faux"
}

// GetType is used to get the TYPE
func (p Provider) GetType() string {
	return "Faux"
}

// ConnectClient is used to connect the client
func (p Provider) ConnectClient(fl *model.Flags) (err error) {
	p.flags = fl

	p.flags = fl

	if p.flags.URL == "" || p.flags.Username == "" || p.flags.Password == "" {
		err = fmt.Errorf("ConnectClient with wrong flags, %v", fl)
		log.WithField("flags", fl).Error("ConnectClient with wrong flags")
	}

	return err
}

// Start is used to start the work
func (p Provider) Start() (err error) {
	return nil
}
