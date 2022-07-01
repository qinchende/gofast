package executors

import "time"

const defaultBulkTasks = 1000

type (
	BulkOption func(options *bulkOptions)

	BulkExecutor struct {
		executor  *IntervalExecutor
		container *bulkContainer
	}

	bulkOptions struct {
		cachedTasks   int
		flushInterval time.Duration
	}
)

func NewBulkExecutor(execute Execute, opts ...BulkOption) *BulkExecutor {
	options := newBulkOptions()
	for _, opt := range opts {
		opt(&options)
	}

	container := &bulkContainer{
		execute:  execute,
		maxTasks: options.cachedTasks,
	}
	executor := &BulkExecutor{
		executor:  NewIntervalExecutor(options.flushInterval, container),
		container: container,
	}

	return executor
}

func (be *BulkExecutor) Add(task any) error {
	be.executor.Add(task)
	return nil
}

func (be *BulkExecutor) Flush() {
	be.executor.Flush()
}

func (be *BulkExecutor) Wait() {
	be.executor.Wait()
}

func WithBulkTasks(tasks int) BulkOption {
	return func(options *bulkOptions) {
		options.cachedTasks = tasks
	}
}

func WithBulkInterval(duration time.Duration) BulkOption {
	return func(options *bulkOptions) {
		options.flushInterval = duration
	}
}

func newBulkOptions() bulkOptions {
	return bulkOptions{
		cachedTasks:   defaultBulkTasks,
		flushInterval: defaultFlushInterval,
	}
}

type bulkContainer struct {
	tasks    []any
	execute  Execute
	maxTasks int
}

func (bc *bulkContainer) AddItem(task any) bool {
	bc.tasks = append(bc.tasks, task)
	return len(bc.tasks) >= bc.maxTasks
}

func (bc *bulkContainer) Execute(tasks any) {
	vals := tasks.([]any)
	bc.execute(vals)
}

func (bc *bulkContainer) RemoveAll() any {
	tasks := bc.tasks
	bc.tasks = nil
	return tasks
}
