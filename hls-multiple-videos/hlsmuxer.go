package main

import (
	"fmt"
	"os"
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
		"-i",
		rtspStream,
		"-i",
		rtspStream,
		"-i",
		rtspStream,
		"-filter_complex",
		"[0:v][0:v]hstack,split[top][bottom];[top][bottom]vstack,format=yuv420p[v]",
		"-map",
		"[v]",
		"-map",
		"0:a?",
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
		"-hls_wrap",
		"15",
		"-f",
		"hls",
		"-hls_flags",
		"delete_segments+second_level_segment_size+second_level_segment_index",
		"-hls_segment_filename",
		path+"/segment_%Y%m%d%H%M%S_%%s_%%d.ts",
		"-use_localtime",
		"1",
		"-vcodec",
		"libx264",
		"-preset",
		"ultrafast",
		"-s",
		"960x540",
		"-r",
		"30",
		"-y",
		path+"/stream.m3u8",
	)

	return muxer
}

func (m *HLSMuxer) Start() error {
	m.cmd.Stderr = os.Stderr
	return m.cmd.Start()
}

func (m *HLSMuxer) Stop() error {
	if m.cmd.Process == nil {
		return fmt.Errorf("you must start the process first")
	}

	return m.cmd.Process.Kill()
}
