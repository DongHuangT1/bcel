# bcel

## Install

```shell
go get github.com/DongHuangT1/bcel
```

## Usage

```go
package main

import (
    "github.com/DongHuangT1/bcel"
)

func main() {
    buf := []byte("...")

    str, err := bcel.Encode(buf, true)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s\n", str)

    ret, err := bcel.Decode(str, true)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s\n", ret)
}
```

## Reference

- https://github.com/apache/commons-bcel
