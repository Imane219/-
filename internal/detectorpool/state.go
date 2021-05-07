package detectorpool

type DetectorState string

func newDetectorState(state string) DetectorState {
	return DetectorState(state)
}

var (
	StateRunning = newDetectorState("running")
	StateInit    = newDetectorState("init")
	StateStopped = newDetectorState("stopped")
	StateNull	= newDetectorState("null")
)

