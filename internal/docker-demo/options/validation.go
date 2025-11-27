package options

func (opts *RunOptions) Validate() []error {
	var errs []error

	errs = append(errs, opts.Log.Validate()...)

	return errs
}
