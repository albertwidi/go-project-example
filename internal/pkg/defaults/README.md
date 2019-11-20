# Defaults

Defaults is a library to set default value to a struct using a struct `tag`. Struct tag `default:"value"` is used to replace struct field value.
The field only set with default value if the field value is empty or zero,  which means:
-  `0` for number 
-  `""` for string

## Boolean is not supported

Boolean is not supported because when the default is `true`, it will always be set to `true`.
Because `false` is similar to `0` for number.

The workaround for this problem is to make your `boolean` field to always has `false` as default value.

## Example

```go
type A struct {
    Str string `default:"asd"`
    I   int    `default:"123"`
}

func main() {
    a := A{
        I:10,
    }
    if err := defaults.SetDefault(&a); err != nil {
        // do something to error
    }
}
```

The `Str` field in struct `A` will be replaced by value in the `default` tag.

## Default Value Replacer (Experimental)

Replacing a struct default value with another struct value. This is useful for overriding some values in a config.

More to be added here.

## Known Issues

The integer and float value is vurnerable to `incompability` of its type. For example, a very large `int` value is being set for `int8` value.

```go
type A struct {
    I8 int8 `default:"100000"`
}
```

Will make `I8` value to `0`