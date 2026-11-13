package main

import (
	"fmt"
	"log"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	width               int
	height              int
	paused              bool
	currentTime         time.Time
	stopwatchRunning    bool
	stopwatchIntervalId int
	stopwatchStartTime  time.Time
	lastElapsed         time.Duration
	stopwatchElapsed    time.Duration
}

type tickMsg time.Time
type stopwatchMsg struct {
	id   int
	time time.Time
}

func tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func stopwatchTick(id int) tea.Cmd {
	return tea.Every(time.Millisecond*100, func(t time.Time) tea.Msg {
		return stopwatchMsg{time: t, id: id}
	})
}

func humanizeDuration(d time.Duration) string {
	d = d.Truncate(100 * time.Millisecond)

	hrs := d / time.Hour
	mins := (d % time.Hour) / time.Minute
	secs := (d % time.Minute).Seconds()

	return fmt.Sprintf("%02d:%02d:%04.1f", hrs, mins, secs)

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
		case "q":
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
		case "s":
			m.stopwatchRunning = !m.stopwatchRunning

			if !m.stopwatchRunning {
				// stopwatch stopped
				m.lastElapsed += time.Since(m.stopwatchStartTime)
				m.stopwatchElapsed = m.lastElapsed

				return m, nil
			} else {
				// stopwatch started
				// new run → new id; ticks still pending from older runs
				// will no longer match and get dropped
				m.stopwatchStartTime = time.Now()
				m.stopwatchIntervalId++
				return m, stopwatchTick(m.stopwatchIntervalId)
			}
		case "r":
			now := time.Now()
			m.lastElapsed = time.Duration(0)
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
		} else if m.stopwatchIntervalId == msg.id {
			m.stopwatchElapsed = m.lastElapsed + time.Since(m.stopwatchStartTime)
			// same run → keep the same id
			return m, stopwatchTick(msg.id)
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	v := m.currentTime.Format("15:04:05")
	if m.paused {
		v += " (paused)"
	}

	v += "\n\n"

	v += humanizeDuration(m.stopwatchElapsed)

	v += "\n\n"

	v += "press s to start/stop the stopwatch, r to reset"

	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, v))
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
