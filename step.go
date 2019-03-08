package stepmachine

type Step interface {
	ID() string
	SetNext(Step)
	SetPrevious(Step)
	Next() Step
	Previous() Step
	Get(string) interface{}
	Set(string, interface{})
	Restore() error
	Run() error
}

type step struct {
	// ID is a unique identifier of the step
	id string

	// restoreFn is a resume function used to restore the current step
	restoreFn func(previousStep, currentStep Step) error

	// runFn is a function to generate the current step
	runFn func(previousStep, current Step) error

	// nextStep represent the next step after the current
	nextStep Step

	// previousStep represent the previous step before the current
	previousStep Step

	// values is a map with values of the current step
	values map[string]interface{}
}

func (s *step) ID() string {
	return s.id
}

func (s *step) Get(key string) interface{} {
	return s.values[key]
}

func (s *step) Set(key string, value interface{}) {
	s.values[key] = value
}

func (s *step) SetNext(n Step) {
	s.nextStep = n
}

func (s *step) Next() Step {
	return s.nextStep
}

func (s *step) SetPrevious(p Step) {
	s.previousStep = p
}

func (s *step) Previous() Step {
	return s.previousStep
}

func (s *step) Restore() error {
	return s.restoreFn(s.previousStep, s)
}

func (s *step) Run() error {
	return s.runFn(s.previousStep, s)
}

func Chain(steps ...Step) {
	if len(steps) < 2 {
		return
	}

	lastStep := steps[0]
	for i := 1; i < len(steps); i++ {
		lastStep.SetNext(steps[i])
		steps[i].SetPrevious(lastStep)
		lastStep = steps[i]
	}
}

func NewStep(id string, runFn func(previousStep, currentStep Step) error, restoreFn func(previousStep, currentStep Step) error) Step {
	return &step{
		id:        id,
		runFn:     runFn,
		restoreFn: restoreFn,
		values:    make(map[string]interface{}),
	}
}
