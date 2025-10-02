package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLinkTag(t *testing.T) {
	header := "<https://api.github.com/user/58276/repos?page=2>; rel=\"next\"," +
		"<https://api.github.com/user/58276/repos?page=10>; rel=\"last\""

	values := ParseLinkTag(header)
	require.ElementsMatch(t, []string{"https://api.github.com/user/58276/repos?page=2", "https://api.github.com/user/58276/repos?page=10"}, values, "could not parse correct links")
}

func TestParseRefreshTag(t *testing.T) {
	header := "999; url=/test/headers/refresh.found"

	values := ParseRefreshTag(header)
	require.Equal(t, "/test/headers/refresh.found", values, "could not parse correct links")
}

func TestTransformIndex(t *testing.T) {
	t.Run("empty array", func(t *testing.T) {
		var arr []int
		require.Equal(t, 0, TransformIndex(arr, 5))
	})

	arr := []int{10, 20, 30}

	t.Run("clamp to first element", func(t *testing.T) {
		require.Equal(t, 0, TransformIndex(arr, -10))
		require.Equal(t, 0, TransformIndex(arr, 0))
		require.Equal(t, 0, TransformIndex(arr, 1))
	})

	t.Run("valid indexes", func(t *testing.T) {
		require.Equal(t, 1, TransformIndex(arr, 2))
		require.Equal(t, 2, TransformIndex(arr, 3))
	})

	t.Run("clamp to last element", func(t *testing.T) {
		require.Equal(t, 2, TransformIndex(arr, 10))
	})
}
