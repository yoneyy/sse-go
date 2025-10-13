# sse-go

## install
```sh
go get -u github.com/yoneyy/sse-go
```

```go
package main

import (
    "log"
    "net/http"
    "github.com/yoneyy/sse-go"
)

func main() {
    http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		sse := NewSSE(w)
		sse.Data("Hello")
		sse.Data("Hello")
		sse.Data("Hello")
		sse.Data("Hello")
		sse.Data("Hello")
		sse.Err("internal error")
		sse.Done()
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
```