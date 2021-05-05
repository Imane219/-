package tester_pool

import (
	"context"
	"contrplatform/configs"
	"log"
	"os/exec"
	"sync"
)

type Tester struct {
	id string
	state TesterState
	cmd *exec.Cmd
	rwLock    sync.RWMutex
	stop context.CancelFunc
}

func (t *Tester) SetState(state TesterState)  {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()
	t.state=state
}

func (t *Tester) State() TesterState {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()
	return t.state
}

func (t *Tester) startCmd(runTime string) error {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()
	ctx,cancel := context.WithCancel(context.Background())
	t.cmd = exec.CommandContext(ctx, "python",configs.TestScriptPath,t.id,runTime)
	//stdout,err:=t.cmd.StdoutPipe()
	//if err != nil {
	//	cancel()
	//	return err
	//}
	if err := t.cmd.Start();err != nil {
		cancel()
		return err
	}
	t.stop=cancel
	t.state=StateRunning
	go t.WaitCmd()

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

func (t *Tester) WaitCmd()  {
	err := t.cmd.Wait()
	if err != nil {
		log.Print(err)
	}
	log.Print("exit")
	t.rwLock.Lock()
	defer t.rwLock.Unlock()
	t.state = StateStopped
}

func (t *Tester) StopCmd() {
	t.stop()
	log.Print("stop")
	//t.rwLock.Lock()
	//defer t.rwLock.Unlock()
	//t.state = StateStopped
}