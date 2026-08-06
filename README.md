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
