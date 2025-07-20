package fp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapErrors(t *testing.T) {
	errs := []error{
		errors.New("e0"),
		errors.New("e1"),
		errors.New("e2"),
		errors.New("e3"),
	}
	err := WrapErrors(errs[0], errs[1:]...)
	assert.Equal(t, "e0: e1: e2: e3", err.Error())
	for _, e := range errs {
		assert.True(t, errors.Is(err, e))
	}
	err = WrapErrors(nil, errs...)
	assert.Equal(t, "e0: e1: e2: e3", err.Error())
	err = WrapErrors(errs[0], nil)
	assert.True(t, errors.Is(err, errs[0]))
}
