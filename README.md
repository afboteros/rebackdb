# rebackdb
Simple library for Golang RethinkDB backup

## Installation

``` bash
$ go get github.com/afboteros/rebackdb
```

RethinkDB dump command should be installed in order to use the library:
```
$ sudo pip install rethinkdb
```

## Usage

``` go
import (
    "github.com/afboteros/rebackdb"
)

func main() {
    file, err := rebackdb.Backup(rebackdb.DumpOptions{
        Connection: "localhost:28015",
        OutputFile: "mydatabase",
    }, rebackdb.FormatISO)

    if err != nil {
        fatal(err)
    }

    err = file.Move("/home/user/mybackups")
    if err != nil {
        fatal(err)
    }
}
```

[RethinkDB Dump]: https://rethinkdb.com/docs/backup/