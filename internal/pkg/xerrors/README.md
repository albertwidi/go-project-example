# XErrors

XErrors package wrap Go 1.13 errors package

This package has `upspin` style error.

## Creating Error

Creating error in `xerrors` is  as simple as in `errors` package

```go
xerrors.New("this is an error)
```

## Error With Kind

Error can be ambigous and hard to categorized, kind is a constant that means to categorized error.

For example:

- Not Found
- Internal Error
- OK
- Bad Request

To create an error with `kind`:

```go
xerrors.New("this is an error", xerrors.KindOK)
```

## Error With Op/Operation

Sometimes function need to be tagged, especially for tracing. With `op` we can tag our error to trace our error more easily.

To create an error with `op`:

```go
xerrors.New(xerrors.Op("doing_something), "this is an error")
```
