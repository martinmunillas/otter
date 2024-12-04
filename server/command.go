package server

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/martinmunillas/otter/server/tools"
)

type CommandInputField struct {
	Name string
	Type string
}

type Command[T any] struct {
	ID      string
	Handler func(r *http.Request, input *T, t tools.Tools)
	fields  []CommandInputField
}

func NewCommand[T any](
	id string,
	handler func(r *http.Request, input *T, t tools.Tools)) Command[T] {
	return Command[T]{
		ID:      id,
		Handler: handler,
	}
}

func (c Command[T]) Handle(r *http.Request, t tools.Tools) {

	if err := r.ParseForm(); err != nil {
		t.Send.BadRequest.JSON("Invalid form data")
		return
	}
	input := new(T)
	err := parseFormIntoInput(r, input)
	if err != nil {
		t.Send.BadRequest.JSON("Invalid form data")
		return
	}
	c.Handler(r, input, t)

}
func (c Command[T]) GetID() string {
	return c.ID
}

func CommandHref(id string) string {
	return fmt.Sprintf("/commands/%s", id)
}

func (c *Command[T]) GetFields() []CommandInputField {
	if len(c.fields) > 0 {
		return c.fields
	}
	val := reflect.ValueOf(new(T)).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		c.fields = append(c.fields, CommandInputField{
			Name: field.Name,
			Type: field.Type.Name(),
		})
	}

	return nil
}

type Commander interface {
	Handle(r *http.Request, t tools.Tools)
	GetID() string
}

func (s *Server) HandleCommands(commands ...Commander) *Server {
	for _, command := range commands {
		s.mux.Handle(fmt.Sprintf("POST %s", CommandHref(command.GetID())), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			command.Handle(r, tools.Make(w, r))
		}))
	}
	return s
}

// parseFormIntoInput populates a struct from form values in the request.
func parseFormIntoInput[T any](r *http.Request, input *T) error {
	val := reflect.ValueOf(input).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		formValue := r.FormValue(field.Name)
		if formValue == "" {
			continue
		}

		structField := val.FieldByName(field.Name)
		if !structField.IsValid() {
			continue // Skip fields that don't exist in the struct
		}

		err := setFieldValue(structField, formValue)
		if err != nil {
			return fmt.Errorf("error setting field %s: %v", field.Name, err)
		}
	}

	return nil
}

// setFieldValue sets a value in a reflect.Value based on its type.
func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("cannot set field")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
