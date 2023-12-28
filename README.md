# Humansize
This library parses a human-readable format for writing data measurements in byte counts, or turns a byte count into a data measurement format string.
## Installation
```
go get github.com/MyBlackJay/humansize
```

## Example
[Usage](https://go.dev/play/p/2PHPVZJUGDm)
```go
package main

import (
	"fmt"
	"github.com/MyBlackJay/humansize"
)

func formatAndPrintKB() {
	size := "100KB"

	if parsing, err := humansize.Compile(size); err == nil {
		fmt.Println(parsing.GetInput(), parsing.GetMeasure(), parsing.GetCompiledUInt64())
	}
}

func formatAndPrintMiB() {
	size := "1MiB"

	if parsing, err := humansize.Compile(size); err == nil {
		fmt.Println(parsing.GetInput(), parsing.GetMeasure(), parsing.GetCompiledUInt64())
	}
}

func MustCompileMiWithError() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	size := "1MR"

	parsing := humansize.MustCompile(size)
	fmt.Println(parsing.GetInput(), parsing.GetMeasure(), parsing.GetCompiledUInt64())
}

func validateMeasureAndPrint() {
	measure := "EiR"
	fmt.Println(humansize.ValidateMeasure(measure))
}

func TurnBytesIntoSizeAndPrint() {
	size := 2.596 * float64(1<<60)
	if res, err := humansize.BytesToSize(size, 10); err == nil {
		fmt.Println(res)
	}
}

func TurnBytesIntoSizeInMeasureAndPrint() {
	if res, err := humansize.Compile("2048MB"); err == nil {
		fmt.Println(res.GetCompiledInMeasure("gib"))
	}
}

func main() {
	formatAndPrintKB()                   // 100KB 1024 102400
	formatAndPrintMiB()                  // 1MiB 1048576 1048576
	MustCompileMiWithError()             // unsupported data size format
	validateMeasureAndPrint()            // false
	TurnBytesIntoSizeAndPrint()          // 2.5960000000EB
	TurnBytesIntoSizeInMeasureAndPrint() // 2 <nil>

}
```

