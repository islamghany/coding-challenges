# Takeways

## GREP : Global Regular Expression Print

- grep is a command-line utility for searching plain-text data sets for lines that match a regular expression.

## Recursive grep with Go

With Go we can directly work with files via `os.File` type which represents an open file descriptor. We can use `os.Open` to open a file and `os.ReadDir` to read the contents of a directory.

This is great but we need to implement the recursive search ourselves. We can do this by using a recursive function that will call itself for each directory it encounters.

### Filesystems and io/fs

The io/fs package defines an interface fs.FS that represents a tree of files.

fs.FS represents A tree of disk files is the obvious example, but if we design our pro‐
gram to operate on an fs.FS value, it can also process ZIP and tar archives, Go mod‐
ules, arbitrary JSON, YAML, or CUE data, or even Web resources addressed by URLs.

Opening a folder as an fs.FS is straightforward. We can do this by calling os.DirFS:

```go
fsys := os.DirFS(".")
```

Notice that there’s no error result to handle here. That’s because we haven’t actually
done any disk operations yet; we’ve just created the abstraction representing the file
tree rooted at ".". If there doesn’t happen to be such a path, or we’re not
allowed to read it, well, too bad: we’ll find that out when we try to do something that
involves reading it.

An fs.FS by itself doesn’t do much. You might be surprised by how small its method
set is:

```go
type FS interface {
    Open(name string) (File, error)
}
```

The Open method returns a File, which is another interface:

```go
type File interface {
    Stat() (fs.FileInfo, error)
    Read([]byte) (int, error)
    Close() error
}
```

**Things we can do with fs.FS:**

### Walking the file tree

Since a filesystem is recursive in nature, the actual recursion operation is always the
same. In principle, the standard library could do this for us, and all we’d need to supply
is the specific code to execute for each file or folder we find.

The fs.WalkDir function does exactly this. It takes a filesystem and some starting
path within it, and recursively walks the tree, visiting every file and folder (in lexical
order; that is, alphabetically).

For each one it finds, it calls some function that you provide, passing it the pathname.
For example:

```go
err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }
    fmt.Println(path)
    return nil
})
```

## Regular Expressions

A regular expression is a sequence of characters that define a search pattern. Usually such patterns are used by string-searching algorithms for "find" or "find and replace" operations on strings, or for input validation.

### Quantifiers

Regex quantifiers check to see how many times you should search for a character.

- `*` - 0 or more e.g. colou\*r will match color, colour, colouur etc.
- `+` - 1 or more e.g. colou+r will match colour but not color
- `?` - 0 or 1 e.g. colou?r will match color and colour
- a|b - Match either “a” or “b”
- `.` - Match any character except newline
- {N} - Exactly N number of occurrences (N is a non-negative integer) .e.g. colou{2}r will match colouur but not colour
- {N,} - N or more occurrences
- {N,M} - Between N and M occurrences .e.g. colou{1,3}r will match colour, colouur and colouuur
- \*? - 0 or more (non-greedy) stops at first true

### Patterns collection

Pattern collections allow you to search for a collection of characters to match against. For example, using the following regex:

```regex
My favorite vowel is [aeiou]
```

This will match the following:

- My favorite vowel is a
- My favorite vowel is e
- My favorite vowel is i ... and so on

Here’s a list of the most common pattern collections:

- [abc] - Match any character in the set
- [A-Z] - Match any character in the range
- [a-z] - Match any character in the range
- [0-9] - Match any character in the range
- [^abc] - Match any character not in the set
- [0-9A-Z] - Match any character in the range

### General Tokens

Not every character is so easily identifiable. While keys like “a” to “z” make sense to match using regex, what about the newline character? Or the tab character? Or even the space character?

- \n - Newline
- \r - Carriage return
- \t - Tab
- \s - any Whitespace character including \n, \r, \t
- \S - any Non-whitespace character
- \w - any Word character (alphanumeric & underscore)
- \W - any Non-word character (the inverse of \w)
- \b - Word boundary: the boundary between a word and a non-word character
- \B - Non-word boundary
- \d - any Digit
- \D - any Non-digit
- ^ - Start of a line
- $ - End of a line
- \ - Escape the next character

### Flags

Flags are used to change the behavior of the regex engine. Here are some of the most common flags:

- g - Global search
- i - Case-insensitive search
- m - Multi-line search ( force ^ and $ to match the start or end of each line)
- s - Single line search ( . matches \n)
