# wasvm

A simple wasm interpreter/runtime (for leaning purpouses)

### Examples

In the `resources` folder there is some `.wasm` files that you can use to run using this project.

Given a `.wat` file that export a function that sum two integers and return the value

```wasm
(module
    (func (export "sum_i32") (param i32) (param i32) (result i32)
        local.get 0
        local.get 1
        i32.sum
    )
)
```

then you should compile it into `.wasm` (the binary format) using `npx wat2wasm {name}.wat -o {name}.wasm`. Once you have the `.wasm` you can use:

```go
wasm, err := parser.BinaryFormat("path to the .wasm file")
if err != nil { panic(err) }

rt, err := vm.NewRuntime(wasm)
if err != nil { panic(err) }

results, err := rt.Exported["sum_i32"].Call(int32(10), int32(5))
if err != nil { panic(err) }

// as we expect only one result then the value will be in the index 0
fmt.Printf("The result of 10 + 5: %v\n", results[0])
```

### Limitations

- Conditions and Loops
- Float points
- Sections:
  - Memory:
  - Global
  - Start
  - Element
  - Data
  - Data Count
  - Import
  - Table

### Running tests

```
go test ./...
```

#### Useful links

- https://webassembly.github.io/spec/
- https://richardanaya.medium.com/lets-write-a-web-assembly-interpreter-part-1-287298201d75
- https://hacks.mozilla.org/2017/02/a-cartoon-intro-to-webassembly/
- https://blog.scottlogic.com/2019/05/17/webassembly-compiler.html
- https://github.com/twitchyliquid64/golang-asm
- https://github.com/tetratelabs/wazero
