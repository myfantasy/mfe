package mfe

import (
	"testing"
	"time"
)

func doTask() {
	_ = 5 * 5
}

func Test_Start(t *testing.T) {

	j := JobCreate(doTask, time.Millisecond)
	j.Start()
	j.Stop()
}
