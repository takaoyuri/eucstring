package eucstring

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

type mode int32

const (
	modeEUCJP mode = iota
	modeUTF8
)

const modeEnv = "EUCSTRING_MODE"

var currentMode atomic.Int32
var initOnce sync.Once
var initErr error

func init() {
	// 既存利用者との互換性のため、明示的な設定がない場合はEUC-JPです。
	currentMode.Store(int32(modeEUCJP))
}

func setMode(value mode) error {
	if value != modeEUCJP && value != modeUTF8 {
		return fmt.Errorf("eucstring: unsupported mode: %d", value)
	}
	currentMode.Store(int32(value))
	return nil
}

// InitFromEnv はEUCSTRING_MODEの値を読み込み、DB文字列モードを設定します。
// 値は "euc-jp" または "utf-8" です。未設定時は既定値のEUC-JPを維持します。
func InitFromEnv() error {
	initOnce.Do(func() {
		initErr = initFromEnv()
	})
	return initErr
}

func initFromEnv() error {
	value, ok := os.LookupEnv(modeEnv)
	if !ok {
		return nil
	}

	switch value {
	case "euc-jp":
		return setMode(modeEUCJP)
	case "utf-8":
		return setMode(modeUTF8)
	default:
		return fmt.Errorf("eucstring: invalid %s value %q (want euc-jp or utf-8)", modeEnv, value)
	}
}

func isUTF8Mode() bool {
	return mode(currentMode.Load()) == modeUTF8
}
