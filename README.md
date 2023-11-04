# ConnectReport generated report converter

The intent is to convert large tabulated reports into a more compact format. Currently working with KDL and much less explicit table structure.

## Build

### WASM

standard
`tinygo build -o kdl-interop.wasm -no-debug main.go`

leaky garbage collector (faster, probably not great for a massive report, is that an issue though?)
`tinygo build -o kdl-interop.wasm -no-debug -gc=leaking main.go`

binaryen doesnt seem to make any significant differences
`wasm-opt -Os input.wasm -o output.wasm`

### WASI

`tinygo build -o kdl-interop.wasm -no-debug -target wasi main.go`
