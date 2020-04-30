package main

import (
	"fmt"
	"log"
	"os/exec"
)

// "-rtsp_transport", "tcp", "-i", this.url, '-f', 'mpegts', '-codec:v', 'mpeg1video', '-bf', '0', '-codec:a', 'mp2', '-r', '30', '-'

type Mpeg1Muxer struct {
	cmd     *exec.Cmd
	onData  func(data []byte)
	onError func(data []byte)
}

func NewMpeg1Muxer(rtspStream string) *Mpeg1Muxer {
	muxer := &Mpeg1Muxer{}
	muxer.cmd = exec.Command("ffmpeg",
		"-rtsp_transport",
		"tcp",
		"-i",
		rtspStream,
		"-f",
		"mpegts",
		"-codec:v",
		"mpeg1video",
		"-bf",
		"0",
		"-codec:a",
		"mp2",
		"-r",
		"30",
		"-",
	)

	return muxer
}

func (m *Mpeg1Muxer) OnData(fn func(data []byte)) {
	m.onData = fn
}

func (m *Mpeg1Muxer) OnError(fn func(data []byte)) {
	m.onError = fn
}

func (m *Mpeg1Muxer) Start() error {
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

func (m *Mpeg1Muxer) Stop() error {
	if m.cmd.Process == nil {
		return fmt.Errorf("you must start the process first")
	}

	return m.cmd.Process.Kill()
}
