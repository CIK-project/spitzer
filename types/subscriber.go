package types

type Subscriber interface {
	NextHeight() int64
	Commit(*BlockResult) error
	Stop() error
}

type TestSubscriber struct {
	Height int64
}

var _ Subscriber = &TestSubscriber{}

func NewTestSubscriber(from int64) *TestSubscriber {
	return &TestSubscriber{
		Height: from,
	}
}

func (test *TestSubscriber) NextHeight() int64 {
	return test.Height + 1
}

func (test *TestSubscriber) Commit(*BlockResult) error {
	test.Height++
	return nil
}

func (test *TestSubscriber) Stop() error {
	return nil
}
