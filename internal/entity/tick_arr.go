package entity

import (
	"time"
)

type RealTimeFutureTickArr []*RealTimeFutureTick

func (c RealTimeFutureTickArr) GetTotalTime() time.Duration {
	if len(c) < 2 {
		return 0
	}

	return c[len(c)-1].TickTime.Sub(c[0].TickTime)
}

func (c RealTimeFutureTickArr) GetLastTwoTickGapTime() time.Duration {
	if len(c) < 2 {
		return 0
	}

	return c[len(c)-1].TickTime.Sub(c[len(c)-2].TickTime)
}

func (c RealTimeFutureTickArr) GetOutInRatioAndRate(duration time.Duration) (float64, float64) {
	if len(c) < 2 {
		return 0, 0
	}

	takenTime := c[len(c)-1].TickTime.Sub(c[0].TickTime)
	if takenTime < duration {
		return 0, 0
	}

	var outVolume, inVolume int64
	for _, v := range c {
		switch v.TickType {
		case 1:
			outVolume += v.Volume
		case 2:
			inVolume += v.Volume
		default:
			outVolume += v.Volume
		}
	}

	return 100 * float64(outVolume) / float64(outVolume+inVolume), float64(outVolume+inVolume) / takenTime.Seconds()
}
