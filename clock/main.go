package main

import (
	"fmt"
	"log"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	width              int
	height             int
	paused             bool
	currentTime        time.Time
	stopwatchRunning   bool
	stopwatchStartTime time.Time
	stopwatchBuffer    time.Duration
	stopwatchElapsed   time.Duration
}

type tickMsg time.Time
type stopwatchMsg time.Time

func tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
func stopwatchTick() tea.Cmd {
	return tea.Every(time.Millisecond*100, func(t time.Time) tea.Msg {
		return stopwatchMsg(t)
	})
}

func humanizeDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	secs := totalSeconds % 60
	mins := totalSeconds / 60 % 60
	hrs := totalSeconds / 60 / 60
	pt := d.Seconds() - float64(totalSeconds)
	floatSecs := float64(secs) + pt
	return fmt.Sprintf("%02d:%02d:%04.1f", hrs, mins, floatSecs)
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		msgStr := msg.String()
		switch msgStr {
		case "ctrl+c":
			return m, tea.Quit
		case "s":
			stopped := m.stopwatchRunning == true

			if stopped {
				// stopwatch stopped
				m.stopwatchBuffer += time.Since(m.stopwatchStartTime)
				m.stopwatchRunning = false

				return m, nil
			} else {
				// stopwatch started
				m.stopwatchStartTime = time.Now()
				m.stopwatchRunning = true
				return m, stopwatchTick()
			}
		case "r":
			now := time.Now()
			m.stopwatchBuffer = time.Duration(0)
			m.stopwatchElapsed = time.Duration(0) // reset stop watch buffer
			m.stopwatchStartTime = now            // reset start time
			return m, nil

		case "space":
			m.paused = !m.paused
			if !m.paused {
				m.currentTime = time.Now()
			}
			return m, nil

		}

	case tickMsg:
		if !m.paused {
			m.currentTime = time.Now()
		}
		return m, tick()

	case stopwatchMsg:
		if !m.stopwatchRunning {
			return m, nil
		} else {
			m.stopwatchElapsed = m.stopwatchBuffer + time.Since(m.stopwatchStartTime)
			return m, stopwatchTick()
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center)

	v := m.currentTime.Format("15:04:05")
	if m.paused {
		v += " (paused)"
	}

	v += "\n\n"

	v += humanizeDuration(m.stopwatchElapsed)

	v += "\n\n"

	v += "press s to start/stop the stopwatch, r to reset"

	return tea.NewView(style.Render(v))
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer f.Close()

	p := tea.NewProgram(model{
		currentTime: time.Now(),
	})

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
