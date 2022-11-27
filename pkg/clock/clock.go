package clock

import "time"

type Clock interface {
	Now() time.Time
}

type System struct{}

func (System) Now() time.Time {
	return time.Now()
}

type Mock struct {
	defaultValue time.Time
	queue        []time.Time
}

func NewMock(defaultValue time.Time) *Mock {
	return &Mock{
		defaultValue: defaultValue,
	}
}

func (c *Mock) Now() time.Time {
	if len(c.queue) == 0 {
		return c.defaultValue
	}

	t := c.queue[0]
	c.queue = c.queue[1:]

	return t
}

func (c *Mock) Add(t time.Time) {
	c.queue = append(c.queue, t)
}
