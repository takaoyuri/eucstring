package eucstring

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// EUCString はEUC-JPエンコードされたデータベースの文字列をUTF-8に変換する型
type EUCString string

// Scan はdatabase/sqlのScannerインターフェースを実装
func (e *EUCString) Scan(src any) error {
	if src == nil {
		*e = ""
		return nil
	}

	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("EUCString: unsupported type: %T", src)
	}

	utf8Str, err := decode(data)
	if err != nil {
		return fmt.Errorf("EUCString: failed to decode database value: %w", err)
	}

	*e = EUCString(utf8Str)
	return nil
}

// ScanText はpgx/v5のTextScannerインターフェースを実装
func (e *EUCString) ScanText(v pgtype.Text) error {
	if !v.Valid {
		*e = ""
		return nil
	}

	utf8Str, err := decode([]byte(v.String))
	if err != nil {
		return fmt.Errorf("EUCString: failed to decode database value: %w", err)
	}

	*e = EUCString(utf8Str)
	return nil
}

// Value はdatabase/sql/driverのValuerインターフェースを実装
func (e EUCString) Value() (driver.Value, error) {
	if isUTF8Mode() {
		if !utf8.ValidString(string(e)) {
			return nil, fmt.Errorf("EUCString: value is not valid UTF-8")
		}
		return string(e), nil
	}

	// UTF-8からEUC-JPに変換して保存
	eucData, err := utf8ToEUCJP(string(e))
	if err != nil {
		return nil, fmt.Errorf("EUCString: failed to convert UTF-8 to EUC-JP: %w", err)
	}
	return eucData, nil
}

func decode(data []byte) (string, error) {
	if isUTF8Mode() {
		if !utf8.Valid(data) {
			return "", fmt.Errorf("value is not valid UTF-8")
		}
		return string(data), nil
	}

	return eucJPToUTF8(data)
}

// String はstringerインターフェースを実装
func (e EUCString) String() string {
	return string(e)
}

// MarshalJSON はJSONマーシャリングを実装
func (e EUCString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

// UnmarshalJSON はJSONアンマーシャリングを実装
func (e *EUCString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*e = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*e = EUCString(s)
	return nil
}

// eucJPToUTF8 はEUC-JPバイト列をUTF-8文字列に変換する
func eucJPToUTF8(data []byte) (string, error) {
	decoder := japanese.EUCJP.NewDecoder()
	utf8Data, _, err := transform.Bytes(decoder, data)
	if err != nil {
		return "", err
	}
	return string(utf8Data), nil
}

// utf8ToEUCJP はUTF-8文字列をEUC-JPバイト列に変換する
func utf8ToEUCJP(str string) ([]byte, error) {
	encoder := japanese.EUCJP.NewEncoder()
	eucData, _, err := transform.Bytes(encoder, []byte(str))
	if err != nil {
		return nil, err
	}
	return eucData, nil
}
