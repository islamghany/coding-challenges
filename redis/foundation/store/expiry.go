package store

import "time"

type ExpireOptions struct {
	EX   time.Time // expire time in seconds
	PX   time.Time // expire time in milliseconds
	EXAT time.Time // expire timestamp-seconds at the specified time
	PXAT time.Time // expire tmestamp-milliseconds  at the specified time in milliseconds
}

func WithEX(ex time.Duration) Option {
	return func(d *DataItem) {
		d.EX = time.Now().Add(ex)
	}
}

func WithPX(px time.Duration) Option {
	return func(d *DataItem) {
		d.PX = time.Now().Add(px)
	}
}

func WithEXAT(exat time.Time) Option {
	return func(d *DataItem) {
		d.EXAT = exat
	}
}

func WithPXAT(pxat time.Time) Option {
	return func(d *DataItem) {
		d.PXAT = pxat
	}
}
