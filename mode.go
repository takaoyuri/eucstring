package eucstring

import (
	"fmt"
	"os"
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

const modeEnv = "EUCSTRING_MODE"

var currentMode atomic.Int32

func init() {
	// 既存利用者との互換性のため、明示的な設定がない場合はEUC-JPです。
	currentMode.Store(int32(ModeEUCJP))
}

// SetMode は現在のプロセスのDB文字列モードを設定します。
// アプリケーション起動時に一度だけ呼び出してください。
func SetMode(mode Mode) error {
	if mode != ModeEUCJP && mode != ModeUTF8 {
		return fmt.Errorf("eucstring: unsupported mode: %d", mode)
	}
	currentMode.Store(int32(mode))
	return nil
}

// InitFromEnv はEUCSTRING_MODEの値を読み込み、DB文字列モードを設定します。
// 値は "euc-jp" または "utf-8" です。未設定時は既定値のEUC-JPを維持します。
func InitFromEnv() error {
	value, ok := os.LookupEnv(modeEnv)
	if !ok {
		return nil
	}

	switch value {
	case "euc-jp":
		return SetMode(ModeEUCJP)
	case "utf-8":
		return SetMode(ModeUTF8)
	default:
		return fmt.Errorf("eucstring: invalid %s value %q (want euc-jp or utf-8)", modeEnv, value)
	}
}

func isUTF8Mode() bool {
	return Mode(currentMode.Load()) == ModeUTF8
}
