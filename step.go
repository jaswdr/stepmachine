package stepmachine

type Step interface {
	ID() string
	SetNext(Step)
	Next() Step
	Run(machine Machine) error
}

type StepFunc func(machine Machine) error

type step struct {
	// ID is a unique identifier of the step
	id string

	// runFn is a function to generate the current step
	runFn StepFunc

	// nextStep represent the next step after the current
	nextStep Step
}

func (s *step) ID() string {
	return s.id
}

func (s *step) SetNext(n Step) {
	s.nextStep = n
}

func (s *step) Next() Step {
	return s.nextStep
}

func (s *step) Run(machine Machine) error {
	return s.runFn(machine)
}

func NewStep(id string, runFn StepFunc) Step {
	return &step{
		id:    id,
		runFn: runFn,
	}
}
