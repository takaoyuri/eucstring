package eucstring

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitFromEnv(t *testing.T) {
	original := Mode(currentMode.Load())
	t.Cleanup(func() { require.NoError(t, SetMode(original)) })

	t.Setenv("EUCSTRING_MODE", "utf-8")
	require.NoError(t, InitFromEnv())

	value, err := EUCString("テスト😀").Value()
	require.NoError(t, err)
	require.Equal(t, "テスト😀", value)
}

func TestInitFromEnvInvalidValue(t *testing.T) {
	original := Mode(currentMode.Load())
	t.Cleanup(func() { require.NoError(t, SetMode(original)) })

	t.Setenv("EUCSTRING_MODE", "invalid")
	require.Error(t, InitFromEnv())
}

func TestInitFromEnvUnsetKeepsCurrentMode(t *testing.T) {
	original := Mode(currentMode.Load())
	t.Cleanup(func() { require.NoError(t, SetMode(original)) })
	require.NoError(t, SetMode(ModeUTF8))
	t.Setenv("EUCSTRING_MODE", "")
	require.NoError(t, os.Unsetenv("EUCSTRING_MODE"))

	require.NoError(t, InitFromEnv())
	value, err := EUCString("テスト😀").Value()
	require.NoError(t, err)
	require.Equal(t, "テスト😀", value)
}

func TestEUCStringUTF8ModeScanAndValue(t *testing.T) {
	original := Mode(currentMode.Load())
	t.Cleanup(func() { require.NoError(t, SetMode(original)) })
	require.NoError(t, SetMode(ModeUTF8))

	var got EUCString
	require.NoError(t, got.Scan([]byte("テスト😀")))
	require.Equal(t, EUCString("テスト😀"), got)

	value, err := got.Value()
	require.NoError(t, err)
	require.Equal(t, "テスト😀", value)
}
