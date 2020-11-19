package errors

import (
	"github.com/pkg/errors"

	"github.com/renehernandez/appfile/internal/log"
)

func CheckAndFailf(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalln(errors.Wrapf(err, msg, args...))
	}
}

func CheckAndFailln(err error, msg string) {
	if err != nil {
		log.Fatalln(errors.Wrap(err, msg))
	}
}

func CheckAndFail(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
