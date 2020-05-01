package main

import (
	"fmt"
	"os"
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
		"-i",
		"dash-with-logo/logo.png",
		"-filter_complex",
		"[1:v]scale=80:80[wat];[0:v][wat]overlay=x=(main_w-overlay_w):y=(main_h-overlay_h)[outv]",
		"-map",
		"[outv]",
		"-map",
		"0:a?",
		"-an",
		"-c:v",
		"libx264",
		"-preset",
		"ultrafast",
		"-b:v",
		"2000k",
		"-f",
		"dash",
		"-window_size",
		"4",
		"-extra_window_size",
		"0",
		"-seg_duration",
		"0.5",
		"-remove_at_exit",
		"1",
		path+"/manifest.mpd",
	)

	return muxer
}

func (m *DASHMuxer) Start() error {
	m.cmd.Stderr = os.Stderr
	return m.cmd.Start()
}

func (m *DASHMuxer) Stop() error {
	if m.cmd.Process == nil {
		return fmt.Errorf("you must start the process first")
	}

	return m.cmd.Process.Kill()
}
