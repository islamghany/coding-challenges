# Takeaways

## Rune and String

### `rune`

is an alias for `int32` and is used to represent a Unicode code point.

The Unicode is a charset contains 2^21 code points, which is enough to represent every character in every language in the world.

We translate the character's list in binary with `UTF-8` encoding, which is a variable-length encoding (from 1 to 4 bytes).

e.g. the character `ض` is identified by the code point `U+0636`. The UTF-8 encoding of this character is `0xD8 0xB6` which are two bytes.

**this is very important to know that the `rune` is not a character, it's a code point.**

### `string`

is a sequence of bytes, not characters. The `string` type is a collection of bytes, not a collection of runes.

so it we want to get the length of a string, we should use the `len` function, which returns the number of bytes in the string.

```go
len("ض") // 2
```

the output is `2` because the character `ض` is represented by two bytes in UTF-8 encoding.
so be aware when looping over a string to do some operations on the characters.

```go
s := "السلام عليكم"
for i := 0; i < len(s); i++ {
    fmt.Println(string(s[i]))
}
```

the output will be the bytes of the string, not the characters.
to get the characters, we should convert the bytes to runes.

```go
s := "السلام عليكم"
for _, r := range s {
    fmt.Println(string(r))
}
```
