package vm

import (
	"errors"
	"fmt"
)

var (
	ErrStackOverflow = errors.New("stackoverflow")
	ErrEmptyStack    = errors.New("empty stack")
)

type StackValue struct {
	value   any
	startAt uint
	endAt   uint
}

type Stack []StackValue

func (s *Stack) push(value StackValue) error {
	defPointer := *s
	if len(defPointer) >= cap(defPointer) {
		return fmt.Errorf("%w: limit %d", ErrStackOverflow, cap(defPointer))
	}

	*s = append(defPointer, value)
	return nil
}

func (s *Stack) pop() (value StackValue, err error) {
	defPointer := *s
	if len(defPointer) == 0 {
		return value, ErrEmptyStack
	}

	removeAt := len(defPointer) - 1
	value = defPointer[removeAt]

	*s = defPointer[:removeAt]
	return value, nil
}

func popEnsureType[T any](s *Stack) (value T, err error) {
	stackValue, err := s.pop()
	if err != nil {
		return value, err
	}

	switch stackValue.value.(type) {
	case T:
		return stackValue.value.(T), nil
	default:
		return value, fmt.Errorf("%w: expected i32. got %T",
			ErrWrongType, stackValue.value)
	}

}
