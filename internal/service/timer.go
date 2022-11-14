package service

import (
	"context"
	"github.com/itdagestan/pomodoro/internal/entity"
	"time"
)

type PomodoroTimerService struct {
}

func NewPomodoroTimerService() *PomodoroTimerService {
	return &PomodoroTimerService{}
}

func (s *PomodoroTimerService) Start(ctx context.Context, pomodoroTimer *entity.PomodoroTimer,
	stopStartedTimer chan bool, updateStarPauseButton chan bool) {
	if pomodoroTimer.IsTimerStart() {
		return
	}
	if pomodoroTimer.IsTimerPause() {
		pomodoroTimer.StartAfterPause()
	} else {
		pomodoroTimer.StartAfterStop()
	}
	updateStarPauseButton <- true

	s.Cancel(ctx, pomodoroTimer, stopStartedTimer)
}

func (s *PomodoroTimerService) Pause(pomodoroTimer *entity.PomodoroTimer,
	stopStartedTimer chan bool, stopUpdateTimer chan bool, updateStarPauseButton chan bool) {
	if pomodoroTimer.IsTimerPause() {
		return
	}
	if pomodoroTimer.IsTimerStart() {
		stopStartedTimer <- true
		stopUpdateTimer <- true
	}
	pomodoroTimer.Pause()
	updateStarPauseButton <- true
}

func (s *PomodoroTimerService) Stop(pomodoroTimer *entity.PomodoroTimer,
	stopStartedTimer chan bool, stopUpdateTimer chan bool, updateStarPauseButton chan bool) {
	if pomodoroTimer.IsTimerStop() {
		return
	}
	if pomodoroTimer.IsTimerStart() {
		stopStartedTimer <- true
		stopUpdateTimer <- true
	}
	pomodoroTimer.Stop()
	updateStarPauseButton <- true
}

func (s *PomodoroTimerService) Skip(pomodoroTimer *entity.PomodoroTimer, skipTimer chan bool,
	stopUpdateTimer chan bool) {
	pomodoroTimer.Skip()
	skipTimer <- true
	if pomodoroTimer.IsTimerStart() {
		stopUpdateTimer <- true
	}
}

func (s *PomodoroTimerService) Cancel(ctx context.Context, pomodoroTimer *entity.PomodoroTimer, cancelTimer chan bool) {
	go func() {
		select {
		case <-time.After(pomodoroTimer.DurationToStop):
			pomodoroTimer.IncrementCountOfDone()
			return
		case <-cancelTimer:
			return
		case <-ctx.Done():
			return
		}
	}()
}
