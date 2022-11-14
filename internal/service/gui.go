package service

import (
	"context"
	"fyne.io/fyne/v2/widget"
	"github.com/itdagestan/pomodoro/internal/entity"
	"log"
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
	log.Println("update timer")
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("ticker")
				timerWidget.SetText(pomodoroTimer.GetTimerAsString())
			case <-ctx.Done():
				return
			case <-stopUpdateTimer:
				log.Println("stop update timer")
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
