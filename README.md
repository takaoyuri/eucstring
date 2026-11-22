# eucstring

Go package for handling EUC-JP database strings as UTF-8 strings.

## Installation

```sh
go get github.com/takaoyuri/eucstring
```

## Usage

```go
import "github.com/takaoyuri/eucstring"

type Book struct {
	Title eucstring.EUCString
}
```

## Features

- `database/sql.Scanner`
- `database/sql/driver.Valuer`
- JSON marshal/unmarshal
- EUC-JPに変換できない文字の置き換え（`ReplaceUnsupported`）

## Replacing unsupported characters

`EUCString`の`Value`は、EUC-JPに変換できない文字（〜 U+301C、絵文字など）があるとエラーを返します。
自由入力欄など、置き換えて保存してよい値は、DBへ書き込む前に`ReplaceUnsupported`で置き換えてください。

```go
replaced, unsupported := eucstring.ReplaceUnsupported(input)
if len(unsupported) > 0 {
	log.Printf("EUC-JPに変換できない文字を?に置き換えました: %q", unsupported)
}
form.Note = eucstring.EUCString(replaced)
```

`ReplaceUnsupported`は次の順で置き換えます。

1. EUC-JPに変換できない文字のうち、対応する文字があるものを変換します。

   | 変換前 | 変換後 |
   | --- | --- |
   | 〜 U+301C WAVE DASH | ～ U+FF5E FULLWIDTH TILDE |
   | − U+2212 MINUS SIGN | － U+FF0D FULLWIDTH HYPHEN-MINUS |
   | — U+2014 EM DASH | ― U+2015 HORIZONTAL BAR |
   | ‖ U+2016 DOUBLE VERTICAL LINE | ∥ U+2225 PARALLEL TO |
   | ¢ U+00A2 CENT SIGN | ￠ U+FFE0 FULLWIDTH CENT SIGN |
   | £ U+00A3 POUND SIGN | ￡ U+FFE1 FULLWIDTH POUND SIGN |
   | ¬ U+00AC NOT SIGN | ￢ U+FFE2 FULLWIDTH NOT SIGN |
   | ¥ U+00A5 YEN SIGN | ￥ U+FFE5 FULLWIDTH YEN SIGN |
   | ‾ U+203E OVERLINE | ￣ U+FFE3 FULLWIDTH MACRON |

2. それでもEUC-JPに変換できない文字と不正なUTF-8のバイトを`?`に置き換えます。

2つ目の戻り値は`?`に置き換えた文字です（出現順、重複なし）。置き換えがなければ`nil`です。
対応表による変換だけが起きたかどうかは、戻り値の文字列と元の文字列を比較して判定します。
`ReplaceUnsupported`は`SetMode`の設定に関係なく、EUC-JPを基準に置き換えます。

## Database encoding migration

`EUCString`は既存環境との互換性のため、既定ではDB値をEUC-JPとして扱います。
UTF-8のDBへ移行する場合は、アプリケーション起動時に設定を変更して変換をスキップできます。

アプリケーションの起動処理で、DB接続を利用する前に`SetMode`を一度だけ呼び出してください。

```go
if err := eucstring.SetMode(eucstring.ModeUTF8); err != nil {
	log.Fatal(err)
}
```

`euc-jp`モードでは読み取り時にEUC-JPからUTF-8へ、書き込み時にUTF-8からEUC-JPへ変換します。
`utf-8`モードでは読み書きとも変換しません。1プロセスが1種類のDBエンコーディングだけを扱う前提です。
EUC-JPとUTF-8のDBを同一プロセスで併用する場合は、このプロセス全体の設定方式ではなく接続単位の型・設定を使用してください。

モード切り替え前に、接続セッションとドライバがアプリケーションの想定するエンコーディングで動作していることを確認してください。
sqlcのoverrideはこれらの値を取得しないため、DB接続後に次のSQLを実行して確認します。

```sql
SHOW server_encoding;
SHOW client_encoding;
```

あわせて、使用するドライバが接続時に返す文字列の形式と、アプリケーションが接続しているDBが想定どおりであることを確認してください。
アプリケーションが受け取る文字列がUTF-8になるよう、接続セッションの`client_encoding`とドライバの返却形式を確認してから、`ModeUTF8`へ切り替えます。

これは移行期間の暫定機能です。移行完了後はsqlcのoverrideを通常の`string`（または変換しないUTF-8専用型）へ変更し、
旧DB向けの`EUCString`を新規用途で使わないことを推奨します。

## License

MIT
