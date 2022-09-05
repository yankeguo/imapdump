package imapdump

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSanitizeMessageID(t *testing.T) {
	require.Equal(t, "aaaA---2223@qq.com---", SanitizeMessageID("<aaaA+*&2223@qq.com)))>"))
}
