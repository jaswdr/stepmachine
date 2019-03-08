package stepmachine

import "errors"

type machine struct {
	initialStep Step
}

var (
	ErrStepNotFound = errors.New("step not found")
)

func (m *machine) Run() (Step, error) {
	currentStep := m.initialStep
	for currentStep != nil {
		if err := currentStep.Run(); err != nil {
			return currentStep, err
		}

		nextStep := currentStep.Next()
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
			return currentStep, err
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

func NewMachine(initialStep Step) *machine {
	return &machine{
		initialStep: initialStep,
	}
}
