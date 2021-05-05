package tester_pool

type TesterState string

func newTesterState(state string) TesterState {
	return TesterState(state)
}

var (
	StateRunning = newTesterState("running")
	StateInit    = newTesterState("init")
	StateStopped = newTesterState("stopped")
	StateNull	= newTesterState("null")
)

