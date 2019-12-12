package command

import (
	"github.com/eure/bobo/engine"
)

type reloadEngineTask struct {
	engine engine.Engine
}

// NewReloadEngineTask is a task to reload engine.
func NewReloadEngineTask(e engine.Engine) reloadEngineTask {
	return reloadEngineTask{
		engine: e,
	}
}

func (reloadEngineTask) GetName() string {
	return "reload_engine_task"
}

func (t reloadEngineTask) Run() error {
	t.engine.Reload()
	return nil
}
