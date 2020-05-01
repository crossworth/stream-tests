package main

import (
	"fmt"
	"os/exec"
)

type DASHMuxer struct {
	cmd *exec.Cmd
}

func NewDASHMuxer(rtspStream string, path string) *DASHMuxer {
	muxer := &DASHMuxer{}
	muxer.cmd = exec.Command("ffmpeg",
		"-rtsp_transport",
		"tcp",
		"-i",
		rtspStream,
		"-an",
		"-c:v",
		"copy",
		"-b:v",
		"2000k",
		"-f",
		"dash",
		"-window_size",
		"4",
		"-extra_window_size",
		"0",
		"-min_seg_duration",
		"2000000",
		"-remove_at_exit",
		"1",
		path+"/manifest.mpd",
	)

	return muxer
}

func (m *DASHMuxer) Start() error {
	return m.cmd.Start()
}

func (m *DASHMuxer) Stop() error {
	if m.cmd.Process == nil {
		return fmt.Errorf("you must start the process first")
	}

	return m.cmd.Process.Kill()
}
