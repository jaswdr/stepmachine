package stepmachine

import (
	"errors"
	"testing"
)

var (
	AnError = errors.New("an error")
)

func runSum(last Step, current Step) error {
	result := 2 + 2
	current.Set("result", result)
	current.Println("runSum executed")
	return nil
}

func restoreSum(last Step, current Step) error {
	current.Set("result", 4)
	current.Set("restored", true)
	current.Println("restoreSum executed")
	return nil
}

func checkSum(last Step, current Step) error {
	if last.Get("result") != 4 {
		return errors.New("invalid last step result")
	}

	current.Set("success", true)
	current.Println("checkSum executed")
	return nil
}

func errorResult(last Step, current Step) error {
	return AnError
}

func TestRunEmptyMachine(t *testing.T) {
	m := NewMachine("test", nil)
	l, err := m.Run()
	if err != nil {
		t.Errorf("running a empty machine returned a error: %s", err)
	}

	if l != nil {
		t.Errorf("last step was not nil: %+v", l)
	}
}

func TestSumStepMachine(t *testing.T) {
	s := NewStep("sum", runSum, restoreSum)
	m := NewMachine("test", s)
	l, err := m.Run()
	if err != nil {
		t.Error(err)
	}

	if l == nil {
		t.Error("last successfull step was nil in a not empty machine")
	}
}

func TestSumStepMachineResultIsAccessable(t *testing.T) {
	s1 := NewStep("step1", runSum, restoreSum)
	s2 := NewStep("step2", checkSum, nil)
	Chain(s1, s2)

	m := NewMachine("test", s1)
	l, err := m.Run()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if l == nil {
		t.Error("last successfull step was nil")
	} else {
		if l.Get("success") == nil {
			t.Error("invalid value")
		}
	}
}

func TestSumStepMachineReturnErrorStepWhenAnErrorHappen(t *testing.T) {
	s1 := NewStep("step1", runSum, restoreSum)
	s2 := NewStep("step2", errorResult, nil)
	Chain(s1, s2)

	m := NewMachine("test", s1)
	l, err := m.Run()
	if err == nil {
		t.Error("an error was expected")
	}

	if l != s2 {
		t.Error("last step was not s2")
	}
}

func TestSumStepMachineResumeToCorrectStep(t *testing.T) {
	s1 := NewStep("step1", runSum, restoreSum)
	s2 := NewStep("step2", checkSum, nil)
	Chain(s1, s2)

	m := NewMachine("test", s1)
	l, err := m.Resume("step2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if l != s2 {
		t.Error("last step was not s2")
	}

	v := s1.Get("restored")
	if v == nil {
		t.Error("s1 was not restored")
	} else {
		if v.(bool) != true {
			t.Error("s1 was not correctly restored")
		}
	}
}

func TestOnStepChangeRuns(t *testing.T) {
	s1 := NewStep("step1", runSum, restoreSum)
	s2 := NewStep("step2", checkSum, nil)
	Chain(s1, s2)

	m := NewMachine("test", s1)

	fromReceived := []Step{}
	toReceived := []Step{}
	m.OnStepChange(func(from, to Step) {
		fromReceived = append(fromReceived, from)
		toReceived = append(toReceived, to)
	})

	m.Run()

	if fromReceived[0] != s1 && toReceived[0] != s2 {
		t.Error("invalid sequence received at index 0")
	}

	if fromReceived[1] != s2 && toReceived[1] != nil {
		t.Error("invalid sequence received at index 1")
	}
}

func TestOnStepErrorRuns(t *testing.T) {
	s1 := NewStep("step1", runSum, restoreSum)
	s2 := NewStep("step2", errorResult, nil)
	Chain(s1, s2)

	m := NewMachine("test", s1)

	var stepReceived Step
	var errorReceived error
	m.OnStepError(func(step Step, err error) {
		stepReceived = step
		errorReceived = err
	})

	m.Run()

	if stepReceived != s2 {
		t.Errorf("unexpected step received: %v", stepReceived)
	}

	if errorReceived != AnError {
		t.Errorf("unexpected error received: %v", errorReceived)
	}
}

func TestOnStepRestoreRuns(t *testing.T) {
	s1 := NewStep("step1", nil, restoreSum)
	s2 := NewStep("step2", checkSum, nil)
	Chain(s1, s2)

	m := NewMachine("test", s1)

	var stepReceived Step
	m.OnStepRestore(func(step Step) {
		stepReceived = step
	})

	m.Resume("step2")
	if stepReceived != s1 {
		t.Errorf("unexpected step was restored: %v", stepReceived)
	}
}
