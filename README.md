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

## Database encoding migration

`EUCString`は既存環境との互換性のため、既定ではDB値をEUC-JPとして扱います。
UTF-8へ移行する期間は、アプリケーション起動時に環境変数を読み込んで変換をスキップできます。

```sh
EUCSTRING_MODE=euc-jp # 移行前
EUCSTRING_MODE=utf-8  # UTF-8 DBへの切り替え後
```

アプリケーションの起動処理で一度だけ`InitFromEnv`を呼び出してください。
不正な値の場合はエラーになります。

```go
if err := eucstring.InitFromEnv(); err != nil {
	log.Fatal(err)
}
```

`euc-jp`モードでは読み取り時にEUC-JPからUTF-8へ、書き込み時にUTF-8からEUC-JPへ変換します。
`utf-8`モードでは読み書きとも変換しません。1プロセスが1種類のDBエンコーディングだけを扱う前提です。
EUC-JPとUTF-8のDBを同一プロセスで併用する場合は、この環境変数方式ではなく接続単位の型・設定を使用してください。

これは移行期間の暫定機能です。移行完了後はsqlcのoverrideを通常の`string`（または変換しないUTF-8専用型）へ変更し、
旧DB向けの`EUCString`を新規用途で使わないことを推奨します。

## License

MIT
