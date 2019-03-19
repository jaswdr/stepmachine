package stepmachine

import (
	"errors"
)

type Machine interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	SetStep(stepID string)
	SetValues(map[string]interface{})
	Values() map[string]interface{}
}

// OnStepChangeFunc is a function called when the machine changes from one step to another
type OnStepChangeFunc func(from, to Step)

// OnStepErrorFunc is a function called when a step returns an error
type OnStepErrorFunc func(step Step, err error)

type machine struct {
	name           string
	initialStep    Step
	values         map[string]interface{}
	onStepChangeFn OnStepChangeFunc
	onStepErrorFn  OnStepErrorFunc
}

var (
	ErrStepNotFound = errors.New("step not found")
)

func (m *machine) Get(key string) interface{} {
	return m.values[key]
}

func (m *machine) Set(key string, value interface{}) {
	m.values[key] = value
}

func (m *machine) SetValues(values map[string]interface{}) {
	if values != nil {
		m.values = values
	}
}

func (m *machine) Values() map[string]interface{} {
	return m.values
}

func (m *machine) SetStep(stepID string) {
	currentStep := m.initialStep
	for currentStep != nil {
		if currentStep.ID() == stepID {
			m.initialStep = currentStep
			return
		}

		currentStep = currentStep.Next()
	}
}

func (m *machine) Resume(stepID string, values map[string]interface{}) {
	m.SetStep(stepID)
	m.SetValues(values)
}

func (m *machine) Run() (Step, error) {
	currentStep := m.initialStep
	for currentStep != nil {
		if err := currentStep.Run(m); err != nil {
			if m.onStepErrorFn != nil {
				m.onStepErrorFn(currentStep, err)
			}

			return currentStep, err
		}

		nextStep := currentStep.Next()
		if m.onStepChangeFn != nil {
			m.onStepChangeFn(currentStep, nextStep)
		}

		if nextStep == nil {
			return currentStep, nil
		}

		currentStep = currentStep.Next()
	}

	return nil, nil
}

func (m *machine) OnStepChange(fn func(from, to Step)) {
	m.onStepChangeFn = fn
}

func (m *machine) OnStepError(fn func(step Step, err error)) {
	m.onStepErrorFn = fn
}

func NewMachine(name string, steps ...Step) *machine {
	// chain steps
	if len(steps) > 1 {
		lastStep := steps[0]
		for i := 1; i < len(steps); i++ {
			lastStep.SetNext(steps[i])
			lastStep = steps[i]
		}
	}

	return &machine{
		name:        name,
		initialStep: steps[0],
		values:      make(map[string]interface{}),
	}
}
