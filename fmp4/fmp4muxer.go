package main

import (
	"fmt"
	"log"
	"os/exec"
)

type FMP4Muxer struct {
	cmd     *exec.Cmd
	onData  func(data []byte)
	onError func(data []byte)
}

func NewFMP4Muxer(rtspStream string) *FMP4Muxer {
	muxer := &FMP4Muxer{}
	muxer.cmd = exec.Command("ffmpeg",
		"-rtsp_transport",
		"tcp",
		"-i",
		rtspStream,
		"-c:v",
		"copy",
		"-c:a",
		"copy",
		"-movflags",
		"frag_keyframe+empty_moov",
		"-f",
		"mp4",
		"-",
	)

	return muxer
}

func (m *FMP4Muxer) OnData(fn func(data []byte)) {
	m.onData = fn
}

func (m *FMP4Muxer) OnError(fn func(data []byte)) {
	m.onError = fn
}

func (m *FMP4Muxer) Start() error {
	stdOut, err := m.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stdErr, err := m.cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		for {
			buf := make([]byte, 1024, 1024)
			n, err := stdOut.Read(buf[:])
			if err != nil {
				log.Println(err)
				return
			}

			if n > 0 {
				m.onData(buf[:n])
			}
		}
	}()

	go func() {
		for {
			buf := make([]byte, 1024, 1024)
			n, err := stdErr.Read(buf[:])
			if err != nil {
				log.Println(err)
				return
			}

			if n > 0 {
				m.onError(buf[:n])
			}
		}
	}()

	return m.cmd.Start()
}

func (m *FMP4Muxer) Stop() error {
	if m.cmd.Process == nil {
		return fmt.Errorf("you must start the process first")
	}

	return m.cmd.Process.Kill()
}
