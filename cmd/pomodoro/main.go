package main

import (
	"context"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/itdagestan/pomodoro/internal/entity"
	"github.com/itdagestan/pomodoro/internal/service"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pomodoroTimer := entity.NewPomodoroTimer()
	services := service.NewService()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	stopStartedTimer := make(chan bool, 1)
	stopUpdateTimer := make(chan bool, 1)
	skipTimer := make(chan bool, 1)
	updateStarPauseButton := make(chan bool, 1)

	a := app.New()
	w := a.NewWindow("Pomodoro")
	timerWidget := widget.NewLabel("00:00")
	counterWidget := widget.NewLabel("Counter 0")

	startPauseButton := widget.NewButton("Start", func() {
		if pomodoroTimer.CanTimerStart() {
			go func() {
				services.PomodoroTimer.Start(ctx, pomodoroTimer, stopStartedTimer, updateStarPauseButton)
				timerWidget.SetText(pomodoroTimer.GetTimerAsString())
				services.Gui.UpdateTimer(ctx, pomodoroTimer, ticker, timerWidget, stopUpdateTimer)
			}()
		}
		if pomodoroTimer.IsTimerStart() {
			services.PomodoroTimer.Pause(pomodoroTimer, stopStartedTimer, stopUpdateTimer,
				updateStarPauseButton)
		}
	})

	stopButton := widget.NewButton("Stop", func() {
		if pomodoroTimer.IsTimerStart() || pomodoroTimer.IsTimerPause() {
			services.PomodoroTimer.Stop(pomodoroTimer, stopStartedTimer, stopUpdateTimer,
				updateStarPauseButton)
			timerWidget.SetText("Stopped")
		}
	})

	skipButton := widget.NewButton("Skip", func() {
		go func() {
			if pomodoroTimer.IsTimerStart() {
				services.PomodoroTimer.Stop(pomodoroTimer, stopStartedTimer, stopUpdateTimer,
					updateStarPauseButton)
			}
			if pomodoroTimer.CanTimerStart() {
				services.PomodoroTimer.Start(ctx, pomodoroTimer, stopStartedTimer, updateStarPauseButton)
				services.Gui.UpdateTimer(ctx, pomodoroTimer, ticker, timerWidget, stopUpdateTimer)
			}
		}()
		services.PomodoroTimer.Skip(pomodoroTimer, skipTimer, updateStarPauseButton)
	})

	services.Gui.UpdateSkipCounter(ctx, pomodoroTimer, counterWidget, skipTimer)
	services.Gui.UpdateStartPauseButton(ctx, pomodoroTimer, startPauseButton, updateStarPauseButton)

	lineOne := container.New(layout.NewHBoxLayout(), timerWidget, counterWidget)
	lineTwo := container.New(layout.NewHBoxLayout(), startPauseButton, stopButton, skipButton)
	w.SetContent(container.New(layout.NewVBoxLayout(), lineOne, lineTwo))
	w.ShowAndRun()
}
