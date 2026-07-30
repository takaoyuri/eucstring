package eucstring

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestInitFromEnv(t *testing.T) {
	original := mode(currentMode.Load())
	t.Cleanup(func() { require.NoError(t, setMode(original)) })

	t.Setenv("EUCSTRING_MODE", "utf-8")
	require.NoError(t, InitFromEnv())

	value, err := EUCString("テスト😀").Value()
	require.NoError(t, err)
	require.Equal(t, "テスト😀", value)
}

func TestInitFromEnvInvalidValue(t *testing.T) {
	original := mode(currentMode.Load())
	t.Cleanup(func() { require.NoError(t, setMode(original)) })

	t.Setenv("EUCSTRING_MODE", "invalid")
	require.Error(t, initFromEnv())
}

func TestInitFromEnvUnsetKeepsCurrentMode(t *testing.T) {
	original := mode(currentMode.Load())
	t.Cleanup(func() { require.NoError(t, setMode(original)) })
	require.NoError(t, setMode(modeUTF8))
	t.Setenv("EUCSTRING_MODE", "")
	require.NoError(t, os.Unsetenv("EUCSTRING_MODE"))

	require.NoError(t, initFromEnv())
	value, err := EUCString("テスト😀").Value()
	require.NoError(t, err)
	require.Equal(t, "テスト😀", value)
}

func TestEUCStringUTF8ModeScanAndValue(t *testing.T) {
	original := mode(currentMode.Load())
	t.Cleanup(func() { require.NoError(t, setMode(original)) })
	require.NoError(t, setMode(modeUTF8))

	var got EUCString
	require.NoError(t, got.Scan([]byte("テスト😀")))
	require.Equal(t, EUCString("テスト😀"), got)

	value, err := got.Value()
	require.NoError(t, err)
	require.Equal(t, "テスト😀", value)
}

func TestEUCStringUTF8ModeScanText(t *testing.T) {
	original := mode(currentMode.Load())
	t.Cleanup(func() { require.NoError(t, setMode(original)) })
	require.NoError(t, setMode(modeUTF8))

	var got EUCString
	require.NoError(t, got.ScanText(pgtype.Text{String: "テスト😀", Valid: true}))
	require.Equal(t, EUCString("テスト😀"), got)

	require.Error(t, got.ScanText(pgtype.Text{String: string([]byte{0xff}), Valid: true}))
	require.NoError(t, got.ScanText(pgtype.Text{Valid: false}))
	require.Equal(t, EUCString(""), got)
}
