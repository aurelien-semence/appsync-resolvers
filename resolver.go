package resolvers

import (
	"encoding/json"
	"reflect"
)

type resolver struct {
	function interface{}
}

func (r *resolver) hasArguments() bool {
	return reflect.TypeOf(r.function).NumIn() >= 1
}

func (r *resolver) hasIdentity() bool {
	return reflect.TypeOf(r.function).NumIn() >= 2
}

func (r *resolver) call(p json.RawMessage, i *json.RawMessage) (interface{}, error) {
	var args []reflect.Value

	if r.hasArguments() {
		pld := payload{p}
		arguments, err := pld.parse(reflect.TypeOf(r.function).In(0))

		if err != nil {
			return nil, err
		}
		args = append(args, *arguments)
	}

	if r.hasIdentity() && i != nil {
		idt := payload{*i}
		identity, err := idt.parse(reflect.TypeOf(r.function).In(1))

		if err != nil {
			return nil, err
		}
		args = append(args, *identity)
	}

	returnValues := reflect.ValueOf(r.function).Call(args)
	var returnData interface{}
	var returnError error

	if len(returnValues) == 2 {
		returnData = returnValues[0].Interface()
	}

	if err := returnValues[len(returnValues)-1].Interface(); err != nil {
		returnError = returnValues[len(returnValues)-1].Interface().(error)
	}

	return returnData, returnError
}
