package service

import (
	"context"
	"fyne.io/fyne/v2/widget"
	"github.com/itdagestan/pomodoro/internal/entity"
	"strconv"
	"time"
)

type GuiService struct {
}

func NewGuiService() *GuiService {
	return &GuiService{}
}

func (s *GuiService) UpdateSkipCounter(ctx context.Context, pomodoroTimer *entity.PomodoroTimer,
	counterWidget *widget.Label, skipTimer chan bool) {
	go func() {
		for {
			select {
			case <-skipTimer:
				counterWidget.SetText("Counter " + strconv.Itoa(pomodoroTimer.CountOfDone))
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *GuiService) UpdateTimer(ctx context.Context, pomodoroTimer *entity.PomodoroTimer,
	ticker *time.Ticker, timerWidget *widget.Label, stopUpdateTimer chan bool) {
	go func() {
		for {
			select {
			case <-ticker.C:
				pomodoroTimer.DurationToStop -= time.Second
				timerWidget.SetText(pomodoroTimer.GetTimerAsString())
			case <-ctx.Done():
				return
			case <-stopUpdateTimer:
				return
			}
		}
	}()
}

func (s *GuiService) UpdateStartPauseButton(ctx context.Context, pomodoroTimer *entity.PomodoroTimer,
	startPauseButton *widget.Button, updateStarPauseButton chan bool) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-updateStarPauseButton:
				if pomodoroTimer.IsTimerPause() || pomodoroTimer.IsTimerStop() {
					startPauseButton.SetText("Start")
				} else {
					startPauseButton.SetText("Pause")
				}
			}
		}
	}()
}
