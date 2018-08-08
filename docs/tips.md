# Tips

Some tips

## Using init() only when needed

Do not use `init()` unless you really-really need it. Using `init()` is confusing sometimes, because it will be executed when a package is imported. The package that import your package might not aware you are using `init()` for some initialization, for eaxmple `parsing flags`