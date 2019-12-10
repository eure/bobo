package command

import (
	"fmt"
	"io"

	"github.com/eure/bobo/engine"
)

type uploadEngineTask struct {
	engine   engine.Engine
	channel  string
	file     io.Reader
	filename string
	filetype string
}

// NewUploadEngineTask is a task to upload a file.
func NewUploadEngineTask(e engine.Engine, channel string, file io.Reader, filename string) *uploadEngineTask {
	return &uploadEngineTask{
		engine:   e,
		channel:  channel,
		file:     file,
		filename: filename,
	}
}

func (t *uploadEngineTask) SetFileType(typ string) {
	t.filetype = typ
}

func (uploadEngineTask) GetName() string {
	return "upload_engine_task"
}

func (t uploadEngineTask) Run() error {
	err := t.engine.FileUploadWithType(t.channel, t.file, t.filename, t.filetype)
	if err != nil {
		fmt.Printf("ERROR uploiad: %v\n", err)
	}
	return err
}
