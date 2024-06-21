# Takeaways

## Lexical Analysis and Syntax Analysis(Parsing)

### Lexical Analysis

Lexical analysis is the first phase of a parser/compiler. It takes the raw input and converts it into tokens. A token is a sequence of characters that represent a unit of meaning. For example, the string `int x = 5;` can be broken down into the following tokens:

`[INT, IDENTIFIER("x"), ASSIGN, Integer(5), SEMICOLON]`

**Source Code** -> **Tokens** -> **Abstract Syntax Tree**

Tokens are the smallest unit of meaning in a programming language. They are used to build an abstract syntax tree (AST) which represents the structure of the program. The AST is then used by the parser to generate the intermediate representation (IR) of the program.

**Note:** Whitespace acts as a delimiter between tokens. It is ignored by the parser.
