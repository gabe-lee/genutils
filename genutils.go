package genutils

import (
	"errors"
	"strings"
)

type DeepError struct {
	IsErr    bool
	Text     string
	Children []DeepError
	Total    uint
}

func NewDeepError(text string) DeepError {
	return DeepError{
		IsErr: true,
		Text:  text,
		Total: 1,
	}
}

func (e *DeepError) Error() string {
	if !e.IsErr {
		return ""
	}
	builder := strings.Builder{}
	e.BuildError(&builder, 0)
	return builder.String()
}

func (e *DeepError) FlatError() error {
	if !e.IsErr {
		return nil
	}
	builder := strings.Builder{}
	e.BuildError(&builder, 0)
	return errors.New(builder.String())
}

func (e *DeepError) AddChildError(err error) {
	if err == nil {
		return
	}
	if e.Children == nil {
		e.Children = make([]DeepError, 0, 1)
	}
	e.IsErr = true
	e.Total += 1
	e.Children = append(e.Children, NewDeepError(err.Error()))
}

func (e *DeepError) AddChildDeepError(err DeepError) {
	if !err.IsErr {
		return
	}
	if e.Children == nil {
		e.Children = make([]DeepError, 0, 1)
	}
	e.IsErr = true
	e.Total += err.Total
	e.Children = append(e.Children, err)
}

func (e *DeepError) BuildError(builder *strings.Builder, depth int) {
	if !e.IsErr {
		return
	}
	tabs := strings.Repeat("\t", depth)
	builder.WriteString("\n")
	builder.WriteString(tabs)
	builder.WriteString(e.Text)
	if e.Children != nil {
		for _, ec := range e.Children {
			ec.BuildError(builder, depth+1)
		}
	}
}
