package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

type PomodoroTimer struct {
	StartedAt          time.Time
	EndedAt            time.Time
	DurationToStop     time.Duration
	ShortRelaxDuration time.Duration
	LongRelaxDuration  time.Duration
	CountToLongSkip    int
	CountOfSkipped     int
	isStop             bool
	isPause            bool
	PausedAt           time.Time
}

func main() {
	pomodoroTimer := PomodoroTimer{
		StartedAt:          time.Time{},
		EndedAt:            time.Time{},
		DurationToStop:     time.Minute * 30,
		ShortRelaxDuration: time.Minute * 10,
		LongRelaxDuration:  time.Minute * 60,
		CountToLongSkip:    4,
		CountOfSkipped:     0,
		isStop:             true,
		isPause:            true,
		PausedAt:           time.Time{},
	}
	a := app.New()
	w := a.NewWindow("Pomodoro")

	timerWidget := widget.NewLabel("")
	counterWidget := widget.NewLabel("Counter 0")
	durationToDone := make(chan string)
	stopUpdateTimerWidgetChan := make(chan bool)
	incrementCount := make(chan bool)
	timerIsStarted := make(chan bool)
	go pomodoroTimer.UpdateCounter(counterWidget, incrementCount)

	startPauseButton := widget.NewButton("Start", func() {
		if pomodoroTimer.canStart() {
			go pomodoroTimer.UpdateTimerWidget(timerWidget, durationToDone, stopUpdateTimerWidgetChan,
				incrementCount)
			timerIsStarted <- true
		}
		if pomodoroTimer.canPause() {
			go pomodoroTimer.pauseTimer()
			stopUpdateTimerWidgetChan <- true
			timerIsStarted <- false
		}
	})
	go pomodoroTimer.updateStartPauseButtonWidget(startPauseButton, timerIsStarted)

	stopButton := widget.NewButton("Stop", func() {
		if pomodoroTimer.canStop() {
			pomodoroTimer.isStop = true
			stopUpdateTimerWidgetChan <- true
			timerWidget.SetText("Stopped")
		}
	})

	skipButton := widget.NewButton("Skip", func() {
		stopUpdateTimerWidgetChan <- true
		incrementCount <- true
		go pomodoroTimer.UpdateTimerWidget(timerWidget, durationToDone, stopUpdateTimerWidgetChan,
			incrementCount)
	})

	lineOne := container.New(layout.NewHBoxLayout(), timerWidget, counterWidget)
	lineTwo := container.New(layout.NewHBoxLayout(), startPauseButton, stopButton, skipButton)
	w.SetContent(container.New(layout.NewVBoxLayout(), lineOne, lineTwo))
	w.ShowAndRun()
}

func (p *PomodoroTimer) StartTimerAfterStop(durationToDone chan string, stopUpdateTimerWidgetChan chan bool,
	stopTimer chan bool, incrementCount chan bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		time.Sleep(p.DurationToStop)
		incrementCount <- true
		stopUpdateTimerWidgetChan <- true
	}()

	for {
		select {
		case <-stopTimer:
			return
		case <-ticker.C:
			durationToDone <- p.GetTimerAsString()
		}
	}
}

func (p *PomodoroTimer) StartTimerAfterPause(durationToDone chan string, stopUpdateTimerWidgetChan chan bool,
	stopTimer chan bool, incrementCount chan bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		time.Sleep(p.DurationToStop)
		incrementCount <- true
		stopUpdateTimerWidgetChan <- true
	}()

	for {
		select {
		case <-stopTimer:
			return
		case <-ticker.C:
			durationToDone <- p.GetTimerAsString()
		}
	}
}

func (p *PomodoroTimer) UpdateTimerWidget(timerWidget *widget.Label, durationToDone chan string,
	stopUpdateTimerWidgetChan chan bool, incrementCount chan bool) {
	if p.canStart() == false {
		return
	}
	stopTimer := make(chan bool)
	if p.isStop {
		p.StartedAt = time.Now()
		p.EndedAt = time.Now().Add(p.DurationToStop)
		go p.StartTimerAfterStop(durationToDone, stopUpdateTimerWidgetChan, stopTimer, incrementCount)
	} else {
		dsd := p.EndedAt.Sub(time.Now())
		p.EndedAt.Add(dsd)
		go p.StartTimerAfterPause(durationToDone, stopUpdateTimerWidgetChan, stopTimer, incrementCount)
	}
	p.isStop = false
	p.isPause = false
	for {
		select {
		case minutes := <-durationToDone:
			timerWidget.SetText(minutes)
		case <-stopUpdateTimerWidgetChan:
			stopTimer <- true
			return
		}
	}
}

func (p *PomodoroTimer) pauseTimer() {
	p.isPause = true
	p.PausedAt = time.Now()
}

func (p *PomodoroTimer) updateStartPauseButtonWidget(stopUpdateTimerWidgetChan *widget.Button,
	timerIsStarted chan bool) {
	for {
		select {
		case c := <-timerIsStarted:
			if c {
				stopUpdateTimerWidgetChan.SetText("Pause")
			} else {
				stopUpdateTimerWidgetChan.SetText("Start")
			}
		}
	}
}

func (p *PomodoroTimer) UpdateCounter(counterWidget *widget.Label, incrementCount chan bool) {
	for {
		select {
		case <-incrementCount:
			p.IncrementCountOfSkipped()
			counterWidget.SetText("Counter " + strconv.Itoa(p.CountOfSkipped))
		}
	}
}

func (p *PomodoroTimer) IncrementCountOfSkipped() {
	p.CountOfSkipped += 1
}

func (p *PomodoroTimer) GetTimerAsString() string {
	minutesInHour := 60
	secondsInMinute := 60

	_, nowMinute, nowSecond := time.Now().Clock()
	_, minuteToEnd, secondToEnd := p.EndedAt.Clock()

	nowMinute = minutesInHour - nowMinute     // 60 - 15 = 45
	minuteToEnd = minutesInHour - minuteToEnd // 60 - 45 = 15
	minuteDifference := 0
	if nowMinute > minuteToEnd {
		minuteDifference = nowMinute - minuteToEnd
	} else {
		minuteDifference = minuteToEnd - nowMinute
	}

	secondDifference := secondToEnd - nowSecond
	if secondDifference < 0 {
		minuteDifference -= 1
		secondDifference = secondsInMinute - (-secondDifference)
	}

	viewSecond := ""
	if secondDifference < 10 {
		viewSecond = "0" + strconv.Itoa(secondDifference)
	} else {
		viewSecond = strconv.Itoa(secondDifference)
	}

	viewMinute := ""
	if minuteDifference < 10 {
		viewMinute = "0" + strconv.Itoa(minuteDifference)
	} else {
		viewMinute = strconv.Itoa(minuteDifference)
	}

	return viewMinute + ":" + viewSecond
}

func (p *PomodoroTimer) canStart() bool {
	return p.isStop == true || p.isPause == true
}

func (p *PomodoroTimer) canPause() bool {
	return p.isPause == false
}

func (p *PomodoroTimer) canStop() bool {
	return p.isStop == false
}
