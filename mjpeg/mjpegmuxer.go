package main

import (
	"fmt"
	"os/exec"
)

type MJPEGMuxer struct {
	cmd *exec.Cmd
}

func NewMJPEGMuxer(rtspStream string, path string) *MJPEGMuxer {
	muxer := &MJPEGMuxer{}
	muxer.cmd = exec.Command("ffmpeg",
		"-rtsp_transport",
		"tcp",
		"-i",
		rtspStream,
		"-c:v",
		"mjpeg",
		"-update",
		"1",
		"-f",
		"image2",
		"-y",
		path+"/stream.jpg",
	)

	return muxer
}

func (m *MJPEGMuxer) Start() error {
	// m.cmd.Stderr = os.Stderr
	return m.cmd.Start()
}

func (m *MJPEGMuxer) Stop() error {
	if m.cmd.Process == nil {
		return fmt.Errorf("you must start the process first")
	}

	return m.cmd.Process.Kill()
}
