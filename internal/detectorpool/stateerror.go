package detectorpool

import "fmt"

type stateError struct {
	err string
}

func newStateError(state DetectorState, want ...DetectorState) *stateError {
	return &stateError{
		err: fmt.Sprintf("错误的测试状态: State-[%s],want-%v",
			state,want),
	}
}

func (s *stateError) Error() string {
	return s.err
}
