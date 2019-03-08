# Step Machine

> Easiest way to handle errors that happen in the middle of a proccess

### Background

When you have a collections of steps that need to run and you need to handle errors when something bad happen you can use this package to run and resume those workloads.

### Gettings started

Let see a brief example:

```go
package main

import (
	"fmt"

	"github.com/jaswdr/stepmachine"
)

func sum(last stepmachine.Step, current stepmachine.Step) error {
	result := 2 + 2
	current.Set("result", result)
	return nil
}

func restoreSum(last stepmachine.Step, current stepmachine.Step) error {
	current.Set("result", 4)
	return nil
}

func addFive(last stepmachine.Step, current stepmachine.Step) error {
	result := last.Get("result").(int)
	result += 5
	current.Set("result", result)
	return nil
}

func main() {
	// declare all steps
	sumStep := stepmachine.NewStep("sum", sum, restoreSum)
	addFiveStep := stepmachine.NewStep("addFive", addFive, nil)

	// chain all steps together
	stepmachine.Chain(sumStep, addFiveStep)

	// create a new step machine with the sumStep as initial step
	m := stepmachine.NewMachine(sumStep)

	// run the step machine
	lastStep, err := m.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println(lastStep.Get("result").(int)) // 9

    // if something bad happen you only need to resume the machine
    lastStep, err = m.Resume("addFiveStep")
    if err != nil {
        panic(err)
    }

	fmt.Println(lastStep.Get("result").(int)) // 9
}
```

## Author
    Jonathan A. Schweder <jonathanschweder@gmail.com>
