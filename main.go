package main

import (
	goflags "github.com/jessevdk/go-flags"

	"github.com/Sirupsen/logrus"

	"github.com/modeneis/uminer/src/drivers"
	"github.com/modeneis/uminer/src/model"
	"github.com/modeneis/uminer/src/version"
)

var (
	log *logrus.Entry
)

func init() {
	log = logrus.WithFields(logrus.Fields{
		"app":     "uminer",
		"env":     model.Env(),
		"version": version.String(),
	})

}

func main() {
	var flags = model.Flags{}
	args, err := goflags.Parse(&flags)
	if err != nil {
		if et, ok := err.(*goflags.Error); ok {
			if et.Type == goflags.ErrHelp {
				return
			}
		}
		log.Fatalf("error parsing flags: %v", err)
	}
	if len(args) > 0 {
		log.Fatalf("unexpected arguments: %v", args)
	}

	if flags.Version {

		log.Printf("uminer Version: %s\n", version.String())
		return
	}

	drivers.Run(flags, log)

}
