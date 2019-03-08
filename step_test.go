package stepmachine

import (
	"strings"
	"testing"
)

func TestChain(t *testing.T) {
	s1 := NewStep("first", nil, nil)
	s2 := NewStep("second", nil, nil)
	s3 := NewStep("third", nil, nil)
	Chain(s1, s2, s3)

	if s1.Previous() != nil {
		t.Error("previous step of s1 was not nil")
	}

	if s1.Next() != s2 {
		t.Error("next step of s1 was not s2")
	}

	if s2.Previous() != s1 {
		t.Error("previous step of s2 was not s1")
	}

	if s2.Next() != s3 {
		t.Error("next step of s2 was not s3")
	}

	if s3.Previous() != s2 {
		t.Error("previous step of s3 was not s2")
	}

	if s3.Next() != nil {
		t.Error("next step of s3 was not nil")
	}
}

func TestLogger(t *testing.T) {
	s := NewStep("log-test", func(last Step, current Step) error {
		current.Println("this is a test")
		return nil
	}, nil)

	s.Run()
	if strings.Index(s.Logs(), "this is a test") < 0 {
		t.Error("invalid log stream")
	}
}
