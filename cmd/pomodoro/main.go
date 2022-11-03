package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"time"
)

type PomodoroTimer struct {
	StartedAt          time.Time
	DurationToStop     time.Duration
	ShortRelaxDuration time.Duration
	LongRelaxDuration  time.Duration
	CountToLongSkip    int32
	CountOfSkipped     int32
}

func main() {
	pomodoroTimer := PomodoroTimer{
		StartedAt:          time.Now(),
		DurationToStop:     time.Minute,
		ShortRelaxDuration: time.Minute * 5,
		LongRelaxDuration:  time.Minute * 10,
		CountToLongSkip:    4,
		CountOfSkipped:     0,
	}
	a := app.New()
	w := a.NewWindow("Pomodoro")

	timerWidget := widget.NewLabel("")
	durationToDone := make(chan string)
	timerIsDone := make(chan bool)
	pomodoroTimer.UpdateTimer(timerWidget, durationToDone, timerIsDone)

	stopButton := widget.NewButton("Stop", func() {
		timerIsDone <- true
	})

	skipButton := widget.NewButton("Skip", func() {
		timerIsDone <- true
		pomodoroTimer.UpdateTimer(timerWidget, durationToDone, timerIsDone)
	})

	lineOne := container.New(layout.NewHBoxLayout(), timerWidget)
	lineTwo := container.New(layout.NewHBoxLayout(), stopButton, skipButton)
	w.SetContent(container.New(layout.NewVBoxLayout(), lineOne, lineTwo))
	w.ShowAndRun()
}

func (p *PomodoroTimer) StartTimerWidget(durationToDone chan string, timerIsDone chan bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		time.Sleep(p.DurationToStop)
		timerIsDone <- true
	}()

	for {
		select {
		case <-timerIsDone:
			durationToDone <- "Done!"
			p.IncrementCountOfSkipped()
			return
		case <-ticker.C:
			minutesBeforeTimerStart := time.Now().Sub(p.StartedAt)
			shortSkipDuration := p.DurationToStop - minutesBeforeTimerStart
			durationToDone <- fmt.Sprintf("%s minutes left", shortSkipDuration/time.Second)
		}
	}
}

func (p *PomodoroTimer) UpdateTimer(timerWidget *widget.Label, durationToDone chan string, timerIsDone chan bool) {
	p.StartedAt = time.Now()
	go p.StartTimerWidget(durationToDone, timerIsDone)
	go func() {
		for {
			select {
			case minutes := <-durationToDone:
				timerWidget.SetText(minutes)
			}
		}
	}()
}

func (p *PomodoroTimer) IncrementCountOfSkipped() {
	p.CountOfSkipped += 1
}
