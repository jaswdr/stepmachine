package stepmachine

import (
	"errors"
)

type machine struct {
	name            string
	initialStep     Step
	onStepChangeFn  func(from, to Step)
	onStepErrorFn   func(step Step, err error)
	onStepRestoreFn func(step Step)
}

var (
	ErrStepNotFound = errors.New("step not found")
)

func (m *machine) Run() (Step, error) {
	currentStep := m.initialStep
	for currentStep != nil {
		if err := currentStep.Run(); err != nil {
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

func (m *machine) Resume(stepID string) (Step, error) {
	currentStep := m.initialStep
	for currentStep.ID() != stepID {
		if err := currentStep.Restore(); err != nil {
			if m.onStepErrorFn != nil {
				m.onStepErrorFn(currentStep, err)
			}

			return currentStep, err
		}

		if m.onStepRestoreFn != nil {
			m.onStepRestoreFn(currentStep)
		}

		nextStep := currentStep.Next()
		if nextStep == nil {
			return nil, ErrStepNotFound
		}

		currentStep = currentStep.Next()
	}

	m.initialStep = currentStep
	return m.Run()
}

func (m *machine) Stack() string {
	stack := "\n+++ START OF STACK +++\n"
	stack += "NAME: " + m.name + "\n\n"
	currentStep := m.initialStep
	for currentStep != nil {
		stack += currentStep.Logs()
		currentStep = currentStep.Next()
	}
	stack += "+++  END OF STACK  +++\n"

	return stack
}

func (m *machine) OnStepChange(fn func(from, to Step)) {
	m.onStepChangeFn = fn
}

func (m *machine) OnStepError(fn func(step Step, err error)) {
	m.onStepErrorFn = fn
}

func (m *machine) OnStepRestore(fn func(step Step)) {
	m.onStepRestoreFn = fn
}

func NewMachine(name string, initialStep Step) *machine {
	return &machine{
		name:        name,
		initialStep: initialStep,
	}
}
