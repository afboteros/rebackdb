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
        Connection:      "localhost:28015",
        OutputFileName:  "mybackup",
        PasswordFile:    "password.txt",
        DateFormat:      rebackdb.FormatShort,
        OperativeSystem: rebackdb.Unix,
    })

    if err != nil {
        log.Printf("There was an error: %s\n Command Output: %s\n Command: %s", err.Error(), file.CommandOutput, file.CommandExecuted)
    }

    file, err = file.Move("./mybackups/")
    if err != nil {
        log.Printf("There was an error: %s\n Command Output: %s\n Command: %s", err.Error(), file.CommandOutput, file.CommandExecuted)
    }
}
```

[RethinkDB Dump]: https://rethinkdb.com/docs/backup/