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

Memoic's design enforces that functions do not have the ability to override each other, although there are situations 
where pipelines can override one another and that situation is when the function's pipelines also returns an item to heap 
with the same name as an existing one, but in general, all pipelines are not able to write to the heap nor can they write 
to each other's heaps.