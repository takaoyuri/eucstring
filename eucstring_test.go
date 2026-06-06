package eucstring

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

var _ pgtype.TextScanner = (*EUCString)(nil)

func TestEUCString_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value EUCString
	}{
		{
			name:  "通常文字列",
			value: EUCString("テスト文字列"),
		},
		{
			name:  "ダブルクォートを含む",
			value: EUCString(`テスト"文字列"`),
		},
		{
			name:  "改行を含む",
			value: EUCString("1行目\n2行目"),
		},
		{
			name:  "バックスラッシュを含む",
			value: EUCString(`C:\tmp\test.txt`),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.value.MarshalJSON()
			require.NoError(t, err)

			want, err := json.Marshal(string(tt.value))
			require.NoError(t, err)

			require.JSONEq(t, string(want), string(got))
			require.Equal(t, string(want), string(got))
		})
	}
}

func TestEUCString_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
		want EUCString
	}{
		{
			name: "normal string",
			data: `"テスト文字列"`,
			want: EUCString("テスト文字列"),
		},
		{
			name: "escaped quote",
			data: `"テスト\"文字列"`,
			want: EUCString(`テスト"文字列`),
		},
		{
			name: "escaped newline",
			data: `"1行目\n2行目"`,
			want: EUCString("1行目\n2行目"),
		},
		{
			name: "unicode escape",
			data: `"\u30c6\u30b9\u30c8"`,
			want: EUCString("テスト"),
		},
		{
			name: "null",
			data: `null`,
			want: EUCString(""),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got EUCString
			err := got.UnmarshalJSON([]byte(tt.data))
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestEUCString_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  any
		want EUCString
	}{
		{
			name: "nil",
			src:  nil,
			want: EUCString(""),
		},
		{
			name: "EUC-JP bytes",
			src:  []byte{0xa4, 0xc6, 0xa4, 0xb9, 0xa4, 0xc8},
			want: EUCString("てすと"),
		},
		{
			name: "EUC-JP string",
			src:  string([]byte{0xa4, 0xc6, 0xa4, 0xb9, 0xa4, 0xc8}),
			want: EUCString("てすと"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got EUCString
			err := got.Scan(tt.src)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestEUCString_ScanUnsupportedType(t *testing.T) {
	t.Parallel()

	var got EUCString
	err := got.Scan(123)
	require.Error(t, err)
}

func TestEUCString_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value EUCString
		want  driver.Value
	}{
		{
			name:  "normal string",
			value: EUCString("てすと"),
			want:  []byte{0xa4, 0xc6, 0xa4, 0xb9, 0xa4, 0xc8},
		},
		{
			name:  "empty string",
			value: EUCString(""),
			want:  []byte{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.value.Value()
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestEUCString_ValueUnsupportedRune(t *testing.T) {
	t.Parallel()

	_, err := EUCString("😀").Value()
	require.Error(t, err)
}
