package hlog

// copy from logrus hooks

// A hook to be fired when logging on the logging severitys returned from
// `Severitys()` on your implementation of the interface. this must be in
// severitys.
//
// Note that this is not fired in a goroutine or a channel with workers,
// you should handle such functionality yourself if your call is non-blocking
// and you don't wish for the logging calls for severitys returned from
// `Severitys()` to block.
type Hook interface {
	Severitys() []int
	Fire(*Entry) error
}

// Internal type for storing the hooks on a logger instance.
type SeverityHooks map[int][]Hook

// Add a hook to an instance of logger. This is called with
// `log.Hooks.Add(new(MyHook))` where `MyHook` implements the `Hook` interface.
func (hooks SeverityHooks) Add(hook Hook) {
	for _, severity := range hook.Severitys() {
		hooks[severity] = append(hooks[severity], hook)
	}
}

// Fire all the hooks for the passed severity. Used by `entry.log` to fire
// appropriate hooks for a log entry.
func (hooks SeverityHooks) Fire(severity int, entry *Entry) error {
	for _, hook := range hooks[severity] {
		if err := hook.Fire(entry); err != nil {
			return err
		}
	}

	return nil
}
