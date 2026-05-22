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

## License

MIT
