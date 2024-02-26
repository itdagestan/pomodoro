package entity

import (
	"time"
)

const defaultDurationToStop = time.Minute * 1
const defaultShortRelaxDuration = time.Minute * 10
const defaultLongRelaxDuration = time.Minute * 30
const defaultCountToLongSkip = 4

type PomodoroTimer struct {
	StartedAt          time.Time
	EndedAt            time.Time
	DurationToStop     time.Duration
	ShortRelaxDuration time.Duration
	LongRelaxDuration  time.Duration
	CountToLongSkip    int
	CountOfDone        int
	IsStop             bool
	IsPause            bool
	PausedAt           time.Time
}

func NewPomodoroTimer() *PomodoroTimer {
	return &PomodoroTimer{
		StartedAt:          time.Time{},
		EndedAt:            time.Time{},
		DurationToStop:     defaultDurationToStop,
		ShortRelaxDuration: defaultShortRelaxDuration,
		LongRelaxDuration:  defaultLongRelaxDuration,
		CountToLongSkip:    defaultCountToLongSkip,
		CountOfDone:        0,
		IsStop:             true,
		IsPause:            false,
		PausedAt:           time.Time{},
	}
}

func (p *PomodoroTimer) StartAfterPause() {
	diff := p.PausedAt.Sub(p.StartedAt)
	now := time.Now()
	p.StartedAt = now
	p.DurationToStop -= diff
	p.IsStop = false
	p.IsPause = false
	p.PausedAt = time.Time{}
}

func (p *PomodoroTimer) StartAfterStop() {
	p.StartedAt = time.Now()
	p.DurationToStop = defaultDurationToStop
	p.IsStop = false
	p.IsPause = false
	p.PausedAt = time.Time{}
}

func (p *PomodoroTimer) Pause() {
	p.IsPause = true
	p.IsStop = false
	p.PausedAt = time.Now()
}

func (p *PomodoroTimer) Stop() {
	p.IsStop = true
	p.IsPause = false
}

func (p *PomodoroTimer) Skip() {
	p.IncrementCountOfDone()
}

func (p *PomodoroTimer) CanTimerStart() bool {
	return p.IsStop == true || p.IsPause == true
}

func (p *PomodoroTimer) IsTimerStart() bool {
	return p.IsStop == false && p.IsPause == false
}

func (p *PomodoroTimer) IsTimerPause() bool {
	return p.IsPause == true
}

func (p *PomodoroTimer) IsTimerStop() bool {
	return p.IsStop == true
}

func (p *PomodoroTimer) IncrementCountOfDone() {
	p.CountOfDone += 1
}

func (p *PomodoroTimer) GetTimerAsString() string {
	return time.Time{}.Add(p.DurationToStop).Format("04:05")
}
