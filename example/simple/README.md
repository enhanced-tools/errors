# Simple add 2 numbers program

It demonstrates how to use errors module.

### Sample output when no args are passed

```
2022/12/24 16:13:48 --- ERROR --- 109c9ba157d418294384e9c8accb8035 --- 5cbbf66a-86b3-42f8-899d-38ea69af1268 
	CONTENT: error 
	{
	  "argCount": 1,
	  "message": "Wrong number of arguments"
	}
	STACK TRACE: 
	github.com/enhanced-tools/errors.enhancedError.From
	/Users/vashingmachine/projects/enhanced-errors/errors.go:150
	github.com/enhanced-tools/errors.enhancedError.FromEmpty
	/Users/vashingmachine/projects/enhanced-errors/errors.go:156
	main.main
	/Users/vashingmachine/projects/enhanced-errors/example/simple/main.go:26
	runtime.main
	/usr/local/go/src/runtime/proc.go:250
	runtime.goexit
	/usr/local/go/src/runtime/asm_arm64.s:1263
```

Where `109c9ba157d418294384e9c8accb8035` represents hash for stack trace and `5cbbf66a-86b3-42f8-899d-38ea69af1268` represents unique id for error.

It means that the first part is unique when errors happens in the same place. The second part is unique for each error.

### Custom error opts

There are 2 custom errors ops declared in `opts.go` file: `Argument` and `ArgCount`. The are used to describe better what happend in each error case.