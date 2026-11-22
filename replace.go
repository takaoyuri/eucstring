package eucstring

import (
	"strings"

	"golang.org/x/text/encoding/japanese"
)

// eucJPSubstitutes はEUC-JPに変換できないが、意味の同じEUC-JP対応文字がある文字の対応表です。
// WAVE DASH問題などで、UnicodeからJIS X 0208への対応が実装によって異なる文字を対象にします。
var eucJPSubstitutes = map[rune]rune{
	0x301C: 0xFF5E, // 〜 WAVE DASH → ～ FULLWIDTH TILDE
	0x2212: 0xFF0D, // − MINUS SIGN → － FULLWIDTH HYPHEN-MINUS
	0x2014: 0x2015, // — EM DASH → ― HORIZONTAL BAR
	0x2016: 0x2225, // ‖ DOUBLE VERTICAL LINE → ∥ PARALLEL TO
	0x00A2: 0xFFE0, // ¢ CENT SIGN → ￠ FULLWIDTH CENT SIGN
	0x00A3: 0xFFE1, // £ POUND SIGN → ￡ FULLWIDTH POUND SIGN
	0x00AC: 0xFFE2, // ¬ NOT SIGN → ￢ FULLWIDTH NOT SIGN
	0x00A5: 0xFFE5, // ¥ YEN SIGN → ￥ FULLWIDTH YEN SIGN
	0x203E: 0xFFE3, // ‾ OVERLINE → ￣ FULLWIDTH MACRON
	0x00A0: 0x0020, // NO-BREAK SPACE → SPACE
}

// unsupportedReplacement はEUC-JPに変換できない文字の置換文字です。
const unsupportedReplacement = '?'

// ReplaceUnsupported は文字列をEUC-JPに必ず変換できる形へ置き換えて返します。
//
// まず、EUC-JPに変換できない文字のうち対応する文字があるもの（〜 U+301C → ～ U+FF5E など）を変換します。
// それでもEUC-JPに変換できない文字（絵文字など）と不正なUTF-8のバイトは '?' に置き換えます。
//
// 2つ目の戻り値は '?' に置き換えた文字です。出現順に、重複を除いて返します。
// 不正なUTF-8のバイトは utf8.RuneError として返します。'?' への置き換えがなければnilです。
// 対応表による変換だけが起きたかどうかは、戻り値の文字列と引数を比較して判定してください。
//
// EUCString の Value は従来どおり変換できない文字があるとエラーを返します。
// 置き換えてよい入力に限り、DBへ書き込む前にこの関数を呼び出してください。
// この関数はSetModeで設定したモードに関係なく、EUC-JPを基準に置き換えます。
func ReplaceUnsupported(s string) (string, []rune) {
	// 対応表の変換元はいずれもEUC-JPに変換できないため、全体を変換できれば置き換えは不要です。
	if _, err := utf8ToEUCJP(s); err == nil {
		return s, nil
	}

	encoder := japanese.EUCJP.NewEncoder()
	var b strings.Builder
	b.Grow(len(s))
	var unsupported []rune
	// seen は '?' に置き換えた文字の重複判定に使います。変換できない文字が見つかるまで作りません。
	var seen map[rune]struct{}
	for _, r := range s {
		if substitute, ok := eucJPSubstitutes[r]; ok {
			r = substitute
		}
		if _, err := encoder.String(string(r)); err != nil {
			if _, ok := seen[r]; !ok {
				if seen == nil {
					seen = make(map[rune]struct{})
				}
				seen[r] = struct{}{}
				unsupported = append(unsupported, r)
			}
			r = unsupportedReplacement
		}
		b.WriteRune(r)
	}
	return b.String(), unsupported
}
