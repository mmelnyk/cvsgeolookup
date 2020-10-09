package cvsgeolookup

type options struct {
	fieldNameBegin      string // required
	fieldNameEnd        string // required
	fieldNameLantitude  string // required
	fieldNameLongtitude string // required
	fieldNameSkip       string // optional
	skipValue           string // optional
	commaRune           rune
	commentRune         rune
	metrics             Metrics
}

type Option interface {
	apply(*options)
}

type fieldBeginOption string

func (o fieldBeginOption) apply(opts *options) {
	opts.fieldNameBegin = string(o)
}

func WithBeginName(name string) Option {
	return fieldBeginOption(name)
}

type fieldEndOption string

func (o fieldEndOption) apply(opts *options) {
	opts.fieldNameEnd = string(o)
}

func WithEndName(name string) Option {
	return fieldEndOption(name)
}

type fieldLantitudeOption string

func (o fieldLantitudeOption) apply(opts *options) {
	opts.fieldNameLantitude = string(o)
}

func WithLantitudeName(name string) Option {
	return fieldLantitudeOption(name)
}

type fieldLongtitudeOption string

func (o fieldLongtitudeOption) apply(opts *options) {
	opts.fieldNameLongtitude = string(o)
}

func WithLongtitudeName(name string) Option {
	return fieldLongtitudeOption(name)
}

type fieldSkipOption string

func (o fieldSkipOption) apply(opts *options) {
	opts.fieldNameSkip = string(o)
}

func WithSkipName(name string) Option {
	return fieldSkipOption(name)
}

type fieldSkipValueOption string

func (o fieldSkipValueOption) apply(opts *options) {
	opts.skipValue = string(o)
}

func WithSkipValue(value string) Option {
	return fieldSkipValueOption(value)
}

type commaValueOption rune

func (o commaValueOption) apply(opts *options) {
	opts.commaRune = rune(o)
}

func WithCommaValue(value rune) Option {
	return commaValueOption(value)
}

type commentValueOption rune

func (o commentValueOption) apply(opts *options) {
	opts.commentRune = rune(o)
}

func WithCommentValue(value rune) Option {
	return commentValueOption(value)
}

type metricsOption struct {
	metrics Metrics
}

func (o metricsOption) apply(opts *options) {
	opts.metrics = o.metrics
}

func WithMetrics(m Metrics) Option {
	return metricsOption{metrics: m}
}
