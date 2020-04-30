package main

import (
	"fmt"
	"os/exec"
)

type HLSMuxer struct {
	cmd *exec.Cmd
}

func NewHLSMuxer(rtspStream string, path string) *HLSMuxer {
	muxer := &HLSMuxer{}
	muxer.cmd = exec.Command("ffmpeg",
		"-rtsp_transport",
		"tcp",
		"-i",
		rtspStream,
		"-fflags",
		"flush_packets",
		"-max_delay",
		"0",
		"-reset_timestamps",
		"1",
		"-flags",
		"-global_header",
		"-hls_time",
		"1",
		"-hls_list_size",
		"10",
		"-f",
		"hls",
		"-hls_flags",
		"delete_segments",
		"-use_localtime",
		"1",
		"-vcodec",
		"copy",
		"-y",
		path+"/stream.m3u8",
	)

	return muxer
}

func (m *HLSMuxer) Start() error {
	return m.cmd.Start()
}

func (m *HLSMuxer) Stop() error {
	if m.cmd.Process == nil {
		return fmt.Errorf("you must start the process first")
	}

	return m.cmd.Process.Kill()
}
