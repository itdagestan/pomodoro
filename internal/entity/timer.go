package entity

import (
	"time"
)

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
		DurationToStop:     time.Minute * 30,
		ShortRelaxDuration: time.Minute * 10,
		LongRelaxDuration:  time.Minute * 60,
		CountToLongSkip:    4,
		CountOfDone:        0,
		IsStop:             true,
		IsPause:            true,
		PausedAt:           time.Time{},
	}
}

func (p *PomodoroTimer) StartAfterPause() {
	diff := p.PausedAt.Sub(p.StartedAt)
	p.StartedAt = time.Now()
	p.EndedAt = time.Now().Add(p.DurationToStop - diff)
	p.IsStop = false
	p.IsPause = false
	p.PausedAt = time.Time{}
}

func (p *PomodoroTimer) StartAfterStop() {
	p.StartedAt = time.Now()
	p.EndedAt = time.Now().Add(p.DurationToStop)
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
	diff := p.EndedAt.Sub(time.Now())
	out := time.Time{}.Add(diff)

	return out.Format("04:05")
}
