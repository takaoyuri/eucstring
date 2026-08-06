package eucstring

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestSetMode(t *testing.T) {
	resetModeStateForTest()
	t.Cleanup(resetModeStateForTest)

	require.NoError(t, SetMode(ModeUTF8))
	value, err := EUCString("テスト😀").Value()
	require.NoError(t, err)
	require.Equal(t, "テスト😀", value)

	require.NoError(t, SetMode(ModeUTF8))
}

func TestSetModeInvalidValue(t *testing.T) {
	resetModeStateForTest()
	t.Cleanup(resetModeStateForTest)
	require.Error(t, SetMode(Mode(99)))
}

func TestSetModeCannotChangeAfterConfiguration(t *testing.T) {
	resetModeStateForTest()
	t.Cleanup(resetModeStateForTest)

	require.NoError(t, SetMode(ModeUTF8))
	err := SetMode(ModeEUCJP)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already configured")
}

func TestEUCStringUTF8ModeScanAndValue(t *testing.T) {
	resetModeStateForTest()
	t.Cleanup(resetModeStateForTest)
	require.NoError(t, SetMode(ModeUTF8))

	var got EUCString
	require.NoError(t, got.Scan([]byte("テスト😀")))
	require.Equal(t, EUCString("テスト😀"), got)

	value, err := got.Value()
	require.NoError(t, err)
	require.Equal(t, "テスト😀", value)

	require.Error(t, got.Scan([]byte{0xff}))
	_, err = EUCString(string([]byte{0xff})).Value()
	require.Error(t, err)
}

func TestEUCStringUTF8ModeScanText(t *testing.T) {
	resetModeStateForTest()
	t.Cleanup(resetModeStateForTest)
	require.NoError(t, SetMode(ModeUTF8))

	var got EUCString
	require.NoError(t, got.ScanText(pgtype.Text{String: "テスト😀", Valid: true}))
	require.Equal(t, EUCString("テスト😀"), got)

	require.Error(t, got.ScanText(pgtype.Text{String: string([]byte{0xff}), Valid: true}))
	require.NoError(t, got.ScanText(pgtype.Text{Valid: false}))
	require.Equal(t, EUCString(""), got)
}

func resetModeStateForTest() {
	modeMu.Lock()
	currentMode.Store(int32(ModeEUCJP))
	modeConfigured = false
	modeMu.Unlock()
}
