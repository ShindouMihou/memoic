package main

import (
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
	fn := memoic.Get("web." + root.Functions[0].Name)
	// simulate the stack since we don't have the runtime complete yet.
	// TODO: complete the runtime stack.
	stack := memoic.Stack{
		Sector: nil,
		Parameters: memoic.Sector{
			"method":  "GET",
			"link":    "https://speed.cloudflare.com/meta",
			"message": "{$local.result}",
		},
		Runtime: &memoic.Runtime{
			Parameters: memoic.Sector{"key": "none"},
			Heap:       memoic.Sector{},
			Stacks:     []memoic.Stack{},
			Result:     nil,
		},
	}
	for _, fnc := range *fn {
		res, err := fnc(&stack)
		if err != nil {
			log.Panicln(err)
		}
		if res != nil {
			stack.Runtime.Heap["result"] = res
		}
	}
}
