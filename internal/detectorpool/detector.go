package detectorpool

import (
	"context"
	"log"
	"os"
	"os/exec"
	"sync"
)

type detector struct {
	id         string
	state      DetectorState
	cmd        *exec.Cmd
	rwMutex    sync.RWMutex
	cancelFunc context.CancelFunc
}

func newDetector(id string) *detector {
	d := &detector{
		id:    id,
		state: StateInit,
	}
	return d
}

func (dt *detector) delete() error {
	dt.rwMutex.Lock()
	if dt.state == StateRunning {
		dt.cancelFunc()
	}
	dt.rwMutex.Unlock()
	if err := os.RemoveAll(poolSetting.OyenteOutputPath + "/" + dt.id); err != nil {
		return err
	}
	if err := os.RemoveAll(poolSetting.SfuzzOutputPath + "/" + dt.id); err != nil {
		return err
	}
	if err := os.RemoveAll(poolSetting.UploadSavePath + "/" + dt.id); err != nil {
		return err
	}
	return nil
}

func (dt *detector) State() DetectorState {
	dt.rwMutex.RLock()
	defer dt.rwMutex.RUnlock()
	return dt.state
}

func (dt *detector) checkOutputState() error {
	dt.rwMutex.RLock()
	defer dt.rwMutex.RUnlock()
	if dt.state != StateRunning && dt.state != StateStopped {
		return newStateError(dt.state)
	}
	return nil
}

func (dt *detector) checkResetDetectionState() error {
	dt.rwMutex.Lock()
	defer dt.rwMutex.Unlock()
	if dt.state != StateStopped {
		return newStateError(dt.state, StateStopped)
	}
	dt.state = StateNull
	return nil
}

func (dt *detector) checkUploadState() error {
	dt.rwMutex.RLock()
	defer dt.rwMutex.RUnlock()
	if dt.state != StateInit {
		return newStateError(dt.state, StateInit)
	}
	return nil
}

func (dt *detector) startCmd(id, runTime string) error {
	dt.rwMutex.Lock()
	defer dt.rwMutex.Unlock()
	if dt.state != StateInit {
		return newStateError(dt.state, StateInit)
	}
	ctx, cancel := context.WithCancel(context.Background())
	dt.cancelFunc = cancel
	dt.cmd = exec.CommandContext(ctx, "python",
		poolSetting.TestScriptPath, id, runTime)
	//stdout,err:=dt.cmd.StdoutPipe()
	//if err != nil {
	//	cancel()
	//	return err
	//}
	if err := dt.cmd.Start(); err != nil {
		cancel()
		return err
	}
	dt.state = StateRunning
	go dt.waitCmd()

	//go func() {
	//	r := bufio.NewReader(stdout)
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			log.Print("exit pro")
	//			return
	//		default:
	//			s,err := r.ReadString('\n')
	//			if err!= nil || err==io.EOF {
	//				return
	//			}
	//			log.Print(s)
	//		}
	//	}
	//}()
	return nil
}

func (dt *detector) waitCmd() {
	err := dt.cmd.Wait()
	if err != nil {
		log.Print(err)
	}
	log.Print("exit")
	dt.rwMutex.Lock()
	defer dt.rwMutex.Unlock()
	dt.state = StateStopped
}

func (dt *detector) stopCmd() error {
	dt.rwMutex.Lock()
	defer dt.rwMutex.Unlock()
	if dt.state == StateRunning {
		dt.cancelFunc()
		log.Print("cancelFunc")
		dt.state = StateStopped
	} else if dt.state != StateStopped {
		return newStateError(dt.state, StateRunning, StateStopped)
	}
	return nil
	//dt.rwMutex.Lock()
	//defer dt.rwMutex.Unlock()
	//dt.state = StateStopped
}
