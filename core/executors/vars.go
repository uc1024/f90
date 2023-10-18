package executors

import "time"

const defaultFlushInterval = time.Second

type placeholderType struct{}

// Execute defines the method to execute tasks.
type Execute func(tasks []any)
