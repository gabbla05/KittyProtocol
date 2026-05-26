package protection

import "time"

// AuthTimer wraps a time.Timer used for the AUTH timeout.
// It is started after a successful HELLO and stopped after AUTH/REGISTER.
type AuthTimer struct {
	timer *time.Timer
}

// StartAuthTimer starts an AUTH timeout timer that calls onTimeout when it fires.
// The timeout duration is controlled by DefaultAuthTimeout.
func StartAuthTimer(onTimeout func()) *AuthTimer {
	return &AuthTimer{
		timer: time.AfterFunc(DefaultAuthTimeout, onTimeout),
	}
}

// Stop cancels the AUTH timer if it is still running.
func (at *AuthTimer) Stop() {
	if at.timer != nil {
		at.timer.Stop()
	}
}
