package stepmachine

import (
	"errors"
	"testing"
)

var (
	AnError = errors.New("an error")
)

func runSum(m Machine) error {
	result := 2 + 2
	m.Set("result", result)
	return nil
}

func checkSum(m Machine) error {
	if m.Get("result") != 4 {
		return errors.New("invalid last step result")
	}

	m.Set("success", true)
	return nil
}

func errorResult(m Machine) error {
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
	s := NewStep("sum", runSum)
	m := NewMachine("test", s)
	l, err := m.Run()
	if err != nil {
		t.Error(err)
	}

	if l != s {
		t.Errorf("unexpected last step, got %s", l.ID())
	}
}

func TestSumStepMachineResultIsAccessable(t *testing.T) {
	s1 := NewStep("step1", runSum)
	s2 := NewStep("step2", checkSum)
	m := NewMachine("test", s1, s2)

	l, err := m.Run()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if l != s2 {
		t.Errorf("last step was %s, not step2", l.ID())
	}

	if m.Get("success") == nil {
		t.Error("invalid value")
	}

	success := m.Get("success").(bool)
	if !success {
		t.Error("unexpected result")
	}
}

func TestSumStepMachineReturnErrorStepWhenAnErrorHappen(t *testing.T) {
	s1 := NewStep("step1", runSum)
	s2 := NewStep("step2", errorResult)

	m := NewMachine("test", s1, s2)
	l, err := m.Run()
	if err == nil {
		t.Error("an error was expected")
	}

	if l != s2 {
		t.Error("last step was not s2")
	}
}

func TestSumStepMachineRunToCorrectStep(t *testing.T) {
	s1 := NewStep("step1", runSum)
	s2 := NewStep("step2", checkSum)

	values := map[string]interface{}{
		"result": 4,
	}
	m := NewMachine("test", s1, s2)
	m.SetStep("step2")
	m.SetValues(values)
	l, err := m.Run()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if l != s2 {
		t.Error("last step was not step2")
	}

	v := m.Get("result")
	if v == nil {
		t.Error("result is incorrect")
	} else {
		if v.(int) == 0 {
			t.Error("result as an invalid value")
		}
	}
}

func TestOnStepChangeRuns(t *testing.T) {
	s1 := NewStep("step1", runSum)
	s2 := NewStep("step2", checkSum)
	m := NewMachine("test", s1, s2)

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
	s1 := NewStep("step1", runSum)
	s2 := NewStep("step2", errorResult)
	m := NewMachine("test", s1, s2)

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

func TestMachineReturnsCorrectlyValues(t *testing.T) {
	s1 := NewStep("step1", runSum)
	m := NewMachine("test", s1)

	m.Run()

	values := m.Values()
	if values["result"] == nil {
		t.Error("result was not set in returned values list")
	}
}

func TestMachineCorrectlySetValuesWhenRun(t *testing.T) {
	values := map[string]interface{}{
		"testing": "values",
	}

	s1 := NewStep("step1", runSum)
	m := NewMachine("test", s1)
	m.SetValues(values)
	m.Run()

	if m.Values()["testing"] != "values" {
		t.Errorf("values where not correctly restored")
	}
}

func TestMachineCanResume(t *testing.T) {
	values := map[string]interface{}{
		"result":  4,
		"testing": "values",
	}

	s1 := NewStep("step1", runSum)
	s2 := NewStep("step2", checkSum)
	m := NewMachine("test", s1, s2)
	m.Resume("step2", values)

	l, err := m.Run()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if l != s2 {
		t.Errorf("last step was %s, not step2", l.ID())
	}
}

func TestMachineIsSet(t *testing.T) {
	m := NewMachine("test")
	if m.IsSet("value") {
		t.Error("value key is expected to be false")
	}

	m.Set("value", false)
	if !m.IsSet("value") {
		t.Error("value key is expected to be true")
	}
}
