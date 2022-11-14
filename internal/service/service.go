package service

import (
	"context"
	"fyne.io/fyne/v2/widget"
	"github.com/itdagestan/pomodoro/internal/entity"
	"time"
)

type PomodoroTimer interface {
	Start(ctx context.Context, pomodoroTimer *entity.PomodoroTimer, stopStartedTimer chan bool,
		updateStarPauseButton chan bool)
	Pause(pomodoroTimer *entity.PomodoroTimer, stopStartedTimer chan bool, stopUpdateStartedTimer chan bool,
		stopTicker chan bool, updateStarPauseButton chan bool)
	Stop(pomodoroTimer *entity.PomodoroTimer,
		stopStartedTimer chan bool, stopUpdateStartedTimer chan bool, stopTicker chan bool,
		updateStarPauseButton chan bool)
	Skip(pomodoroTimer *entity.PomodoroTimer, skipTimer chan bool, updateStarPauseButton chan bool)
	Cancel(ctx context.Context, pomodoroTimer *entity.PomodoroTimer, stopTimer chan bool)
}

type Gui interface {
	UpdateSkipCounter(ctx context.Context, pomodoroTimer *entity.PomodoroTimer, counterWidget *widget.Label,
		skipTimer chan bool)
	UpdateTimer(ctx context.Context, pomodoroTimer *entity.PomodoroTimer,
		ticker *time.Ticker, timerWidget *widget.Label, stopUpdateStartedTimer chan bool)
	UpdateStartPauseButton(ctx context.Context, pomodoroTimer *entity.PomodoroTimer, startPauseButton *widget.Button,
		updateStarPauseButton chan bool)
}

type Service struct {
	PomodoroTimer PomodoroTimer
	Gui           Gui
}

func NewService() *Service {
	return &Service{
		PomodoroTimer: NewPomodoroTimerService(),
		Gui:           NewGuiService(),
	}
}
