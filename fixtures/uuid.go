package fixtures

import (
	"conformity-core/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func UUID(t *testing.T) string {
	id, err := utils.UUID()
	require.Nil(t, err)

	return id
}
