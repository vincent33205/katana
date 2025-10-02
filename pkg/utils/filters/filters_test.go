package filters

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleFilter(t *testing.T) {
	simple, err := NewSimple()
	require.NoError(t, err, "could not create filter")
	defer simple.Close()

	unique := simple.UniqueURL("https://example.com")
	require.True(t, unique, "could not get unique value")

	unique = simple.UniqueURL("https://example.com")
	require.False(t, unique, "could get unique value")
}

func TestSimpleFilterUniqueContent(t *testing.T) {
	simple, err := NewSimple()
	require.NoError(t, err, "could not create filter")
	defer simple.Close()

	payload := []byte("katana")

	unique := simple.UniqueContent(payload)
	require.True(t, unique, "expected new payload to be unique")

	unique = simple.UniqueContent(payload)
	require.False(t, unique, "expected duplicate payload to be rejected")
}

func TestSimpleFilterIsCycle(t *testing.T) {
	simple := &Simple{}

	t.Run("long url", func(t *testing.T) {
		url := strings.Repeat("a", MaxChromeURLLength+1)
		require.True(t, simple.IsCycle(url), "expected overly long url to be considered a cycle")
	})

	t.Run("repeating sequence", func(t *testing.T) {
		var builder strings.Builder
		base := "abcdefghijkl"
		for i := 0; i < MaxSequenceCount; i++ {
			builder.WriteString(base)
			builder.WriteString(fmt.Sprintf("%02d", i))
		}

		require.True(t, simple.IsCycle(builder.String()), "expected highly repetitive sequence to be considered a cycle")
	})

	t.Run("valid url", func(t *testing.T) {
		require.False(t, simple.IsCycle("https://example.com"), "expected typical url to not be a cycle")
	})
}
