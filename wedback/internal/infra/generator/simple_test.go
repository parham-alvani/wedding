package generator_test

import (
	"testing"

	generatorD "github.com/parham-alvani/wedding/wedback/internal/domain/generator"
	generatorI "github.com/parham-alvani/wedding/wedback/internal/infra/generator"
	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	t.Parallel()

	s := new(generatorI.Simple)

	require.Implements(t, new(generatorD.Generator), s)
	require.Len(t, s.ID(), 10)
}
