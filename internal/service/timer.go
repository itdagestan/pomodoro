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
	stopStartedTimer chan bool, stopUpdateStartedTimer chan bool, updateStarPauseButton chan bool) {
	if pomodoroTimer.IsTimerPause() {
		return
	}
	pomodoroTimer.Pause()
	stopStartedTimer <- true
	stopUpdateStartedTimer <- true
	updateStarPauseButton <- true
}

func (s *PomodoroTimerService) Stop(pomodoroTimer *entity.PomodoroTimer,
	stopStartedTimer chan bool, stopUpdateStartedTimer chan bool, updateStarPauseButton chan bool) {
	if pomodoroTimer.IsTimerStop() {
		return
	}
	pomodoroTimer.Stop()
	stopStartedTimer <- true
	stopUpdateStartedTimer <- true
	updateStarPauseButton <- true
}

func (s *PomodoroTimerService) Skip(pomodoroTimer *entity.PomodoroTimer, skipTimer chan bool,
	updateStarPauseButton chan bool) {
	pomodoroTimer.Skip()
	skipTimer <- true
	updateStarPauseButton <- true
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
