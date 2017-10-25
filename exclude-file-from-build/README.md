# Conditional Compilation Example

If you run `go build` from this directory, it will fail. This is because the 
source code file **excluded.go** is not included in compilation.

Why? Because the first two lines of **excluded.go** look like this:

```go
// +build excluded

```

That's a comment followed by (important!) a blank line. The "excluded" bit can
be any word.

To build this package successfully, you need to tell the compiler to include the
file.

```
go build -tags=excluded
```

