package eucstring

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Mode はDB文字列のエンコーディングモードです。
type Mode int32

const (
	// ModeEUCJP はDB値をEUC-JPとして扱います。
	ModeEUCJP Mode = iota
	// ModeUTF8 はDB値をUTF-8として扱い、変換しません。
	ModeUTF8
)

var currentMode atomic.Int32
var modeMu sync.Mutex
var modeConfigured bool

func init() {
	// 既存利用者との互換性のため、既定値はEUC-JPです。
	currentMode.Store(int32(ModeEUCJP))
}

// SetMode は現在のプロセスのDB文字列モードを設定します。
// アプリケーション起動時に、DB接続を利用する前に一度だけ呼び出してください。
func SetMode(value Mode) error {
	if value != ModeEUCJP && value != ModeUTF8 {
		return fmt.Errorf("eucstring: unsupported mode: %d", value)
	}

	modeMu.Lock()
	defer modeMu.Unlock()
	if modeConfigured && Mode(currentMode.Load()) != value {
		return fmt.Errorf("eucstring: mode is already configured as %s", modeName(Mode(currentMode.Load())))
	}
	currentMode.Store(int32(value))
	modeConfigured = true
	return nil
}

func isUTF8Mode() bool {
	return Mode(currentMode.Load()) == ModeUTF8
}

func modeName(value Mode) string {
	if value == ModeUTF8 {
		return "utf-8"
	}
	return "euc-jp"
}
