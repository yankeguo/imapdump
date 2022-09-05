package imapdump

import (
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDumpMessagePath(t *testing.T) {
	path := DumpMessagePath("testdata", &imap.Message{
		Envelope: &imap.Envelope{
			Date:      time.Date(2022, time.February, 1, 12, 12, 12, 0, time.Local),
			MessageId: "<aAA@#$%^@qq.com>",
		},
	})
	require.Equal(t, filepath.Join("testdata", "2022", "02", "aAA@----@qq.com.eml"), path)
}

func TestDumpMessagePathExisted(t *testing.T) {
	msg := &imap.Message{
		Envelope: &imap.Envelope{
			Date:      time.Date(2022, time.February, 1, 12, 12, 12, 0, time.Local),
			MessageId: "<aAA@#$%^@qq.com>",
		},
	}
	path := DumpMessagePath("testdata", msg)
	_ = os.RemoveAll(path)
	ok, err := DumpMessagePathExisted("testdata", msg)
	require.NoError(t, err)
	require.False(t, ok)
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	_ = os.WriteFile(path, []byte("EXISTED"), 0640)
	ok, err = DumpMessagePathExisted("testdata", msg)
	require.NoError(t, err)
	require.True(t, ok)
}
