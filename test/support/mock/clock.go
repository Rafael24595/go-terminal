package mock

type TestClock struct {
	Time int64
}

func (c *TestClock) Now() int64 {
	return c.Time
}

func (c *TestClock) Advance(ms int64) {
	c.Time += ms
}

func FixedClock(t int64) func() int64 {
	return func() int64 {
		return t
	}
}
