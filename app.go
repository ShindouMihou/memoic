package main

import (
	"fmt"
	"io"
	"log"
	"memoic/internal/memoic"
	"memoic/internal/memoic/native"
	"memoic/pkg/memoize"
	"os"
)

func main() {
	f, _ := os.Open("functions/examples/web/load.json")
	defer f.Close()
	bytes, _ := io.ReadAll(f)
	root, err := memoize.Parse(bytes)
	if err != nil {
		log.Panicln(err)
	}
	native.Discover()
	err = memoic.ImprintFn(root)
	if err != nil {
		log.Panicln(err)
	}
	fn := memoic.Get("web.load")
	runtime := memoic.NewRuntime(memoic.Sector{
		"link": "https://jsonplaceholder.typicode.com/todos/1",
	})
	err = runtime.Load(fn)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(*runtime.Result)
}
