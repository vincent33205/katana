package queue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrategyString(t *testing.T) {
	require.Equal(t, "breadth-first", BreadthFirst.String(), "expected breadth first string representation")
	require.Equal(t, "depth-first", DepthFirst.String(), "expected depth first string representation")
	require.Equal(t, "", Strategy(-1).String(), "expected unknown strategy to return empty string")
}
