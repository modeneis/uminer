package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/modeneis/uminer/src/providers/faux"
)

func Test_UseProviders(t *testing.T) {
	a := assert.New(t)

	//create test  provider
	fb := &faux.Provider{}
	UseProviders(fb)

	a.Equal(len(GetProviders()), 1)
	a.Equal(GetProviders()[fb.GetType()], fb)
	ClearProviders()
}

func Test_GetProvider(t *testing.T) {
	a := assert.New(t)

	fa := &faux.Provider{}
	UseProviders(fa)

	p, err := GetProvider(fa.GetType())
	a.NoError(err)
	a.Equal(p, fa)

	_, err = GetProvider("unknown")
	a.Error(err)
	a.Equal(err.Error(), "no provider for unknown exists")
	ClearProviders()
}
