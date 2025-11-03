package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	fileName   string
	totalBytes int64
	bytesRead  int64
	checksum   []byte
	progress   float64
	errorMsg   string
}

type updateMsg struct {
	progress  float64
	bytesRead int64
}

type doneMsg struct {
	result []byte
}
type errMsg struct {
	errorMsg string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}

	totalBytes := fileInfo.Size()
	m := model{
		fileName:   filename,
		totalBytes: totalBytes,
	}

	p := tea.NewProgram(m)

	go func() {
		// Get the hash, read 4k chunks, send msgs to bubbletea

		hasher := sha256.New()
		const chunkSize = 4096
		buffer := make([]byte, chunkSize)

		var bytesRead int64 = 0
		for {
			n, err := file.Read(buffer)
			if n > 0 {
				hasher.Write(buffer[:n])
				bytesRead += int64(n)

				// Note: blocks until receiver is ready to process msgs
				p.Send(updateMsg{
					progress:  float64(m.bytesRead) / float64(m.totalBytes) * 100,
					bytesRead: bytesRead,
				})
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				p.Send(errMsg{
					errorMsg: fmt.Sprintf("Error reading file: %v\n", err),
				})
				return
			}
		}

		p.Send(tea.Msg(doneMsg{
			result: hasher.Sum(nil),
		}))
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
		return
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateMsg:
		m.bytesRead = msg.bytesRead
		m.progress = msg.progress
		return m, nil
	case doneMsg:
		m.checksum = msg.result
		return m, tea.Quit
	case errMsg:
		m.errorMsg = msg.errorMsg
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	if m.errorMsg != "" {
		return fmt.Sprintf("ERROR: %s\n", m.errorMsg)
	}
	if m.checksum != nil {
		return fmt.Sprintf("sha256: %x %s\n", m.checksum, m.fileName)
	}
	return fmt.Sprintf("Bytes read %d, progress: %.2f%%\n", m.bytesRead, m.progress)
}
