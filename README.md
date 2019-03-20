# Step Machine

> Easiest way to handle errors that happen in the middle of a proccess

[![Build Status](https://travis-ci.org/jaswdr/stepmachine.svg?branch=master)](https://travis-ci.org/jaswdr/stepmachine)

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

func addFive(last stepmachine.Step, current stepmachine.Step) error {
	result := last.Get("result").(int)
	result += 5
	current.Set("result", result)
	return nil
}

func main() {
	// declare all steps
	sumStep := stepmachine.NewStep("sum", sum)
	addFiveStep := stepmachine.NewStep("addFive", addFive)

	// create a new step machine with the sumStep as initial step
	m := stepmachine.NewMachine("my step machine name", sumStep, addFiveStep)

	// run the step machine
	lastStep, err := m.Run("", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(lastStep.Get("result").(int)) // 9

	// if something bad happen you only need to resume the machine
    	savedValues := map[string]interface{}{
        	"result": 5,
    	}
    	m.Resume("addFiveStep", savedValues)
	_, err = m.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println(m.Get("result").(int)) // 9
}

```

### Logging

Is possible to log inside a step and print in the end, use `Step.Println()` to log a message and `Machine.Stack()` to retrieve the log stack with all messages.

## Author
    Jonathan A. Schweder <jonathanschweder@gmail.com>
