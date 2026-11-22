package eucstring

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestReplaceUnsupported(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		input           string
		want            string
		wantUnsupported []rune
	}{
		{
			name:  "置換なし",
			input: "受付テスト①髙橋 ～－―∥￠￡￢￥￣ abc\t\n",
			want:  "受付テスト①髙橋 ～－―∥￠￡￢￥￣ abc\t\n",
		},
		{
			name:  "空文字列",
			input: "",
			want:  "",
		},
		{
			name:  "WAVE DASH",
			input: "10時〜12時",
			want:  "10時～12時",
		},
		{
			name:  "対応表の全文字",
			input: "〜−—‖¢£¬¥‾",
			want:  "～－―∥￠￡￢￥￣",
		},
		{
			name:            "絵文字",
			input:           "よろしく😀お願いします🙏",
			want:            "よろしく?お願いします?",
			wantUnsupported: []rune{'😀', '🙏'},
		},
		{
			name:            "同じ文字は重複せずに返す",
			input:           "😀a😀",
			want:            "?a?",
			wantUnsupported: []rune{'😀'},
		},
		{
			name:            "対応表の変換と?置換の混在",
			input:           "〜😀− ",
			want:            "～?－?",
			wantUnsupported: []rune{'😀', ' '},
		},
		{
			name:            "不正なUTF-8",
			input:           "a\xffb",
			want:            "a?b",
			wantUnsupported: []rune{utf8.RuneError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, unsupported := ReplaceUnsupported(tt.input)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantUnsupported, unsupported)

			_, err := utf8ToEUCJP(got)
			require.NoError(t, err)
		})
	}
}

func TestReplaceUnsupported_ValueSucceeds(t *testing.T) {
	t.Parallel()

	_, err := EUCString("〜😀").Value()
	require.Error(t, err, "Valueの既存動作は変換できない文字でエラーになる")

	got, _ := ReplaceUnsupported("〜😀")
	_, err = EUCString(got).Value()
	require.NoError(t, err)
}

func TestEUCJPSubstitutes(t *testing.T) {
	t.Parallel()

	for from, to := range eucJPSubstitutes {
		_, err := utf8ToEUCJP(string(from))
		require.Errorf(t, err, "U+%04X はEUC-JPに変換できるため対応表は不要", from)

		_, err = utf8ToEUCJP(string(to))
		require.NoErrorf(t, err, "U+%04X の変換先 U+%04X はEUC-JPに変換できる必要がある", from, to)
	}
}

func FuzzReplaceUnsupported(f *testing.F) {
	for _, seed := range []string{"", "テスト", "〜😀", "a\xffb", " —"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		got, unsupported := ReplaceUnsupported(s)
		_, err := utf8ToEUCJP(got)
		require.NoError(t, err)
		if unsupported == nil {
			require.Equal(t, utf8.RuneCountInString(s), utf8.RuneCountInString(got))
		}
	})
}
