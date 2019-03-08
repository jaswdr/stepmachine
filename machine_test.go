package stepmachine

import (
	"errors"
	"testing"
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
	return errors.New("an error")
}

func TestRunEmptyMachine(t *testing.T) {
	m := NewMachine(nil)
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
	m := NewMachine(s)
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

	m := NewMachine(s1)
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

	m := NewMachine(s1)
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

	m := NewMachine(s1)
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
