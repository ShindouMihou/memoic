# Memoic

##

A proof-of-concept, mainly for entertainment, of an in-memory cache service that builds with a functional approach, 
Memoic is based on the concept of memoization, like in React, the service enables developers to create functions using 
a pipeline approach (similar to MongoDB's aggregation), and use them to link keys to functions, allowing automatic 
fetching, refreshing, etc. of data.

## Exploring Pipelines

Memoic isn't meant to be a complete service, but rather a proof-of-concept built for learning. Memoic uses `.json` to 
develop outlines for the function, with built-in native functions backing the infrastructure, this outline is then 
decoded into native Golang structures and then the pipelines are transformed into native anonymous functions.

Pipelines can contain further more pipelines which takes priority over the parent pipeline, this can all be found in the 
[`internal/memoic/functions.go`](internal/memoic/functions.go) code.

So far, we have tested pipelines to work with the [`functions/examples/web/load.json`](functions/examples/web/load.json) under 
a simulated runtime and stack, further work is needed to complete the runtime and stack implementation. Pipelines have their own 
stack which contains their own memory space to store local variables with a runtime that contains a global memory space that is 
accessible to all functions, this prevents functions from overriding one another or accidentally using one another's variables, though 
further more work is needed to complete this section.