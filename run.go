package trinity

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var (
	contextType = reflect.ValueOf(context.Background()).Type()
)

func (t *Trinity) Run(instanceName string, f interface{}) error {
	if f == nil {
		return errors.New("func cannot be null")
	}
	if !t.container.CheckInstanceNameIfExist(instanceName) {
		return errors.New("instance name not exist")
	}
	ftype := reflect.TypeOf(f)
	if ftype.NumIn() > 2 || ftype.NumIn() < 0 {
		return errors.New("func in must between 0-2")
	}
	if ftype.NumOut() != 1 {
		return errors.New("func out must be 1")
	}
	if !IsFunc(ftype) {
		return fmt.Errorf("func type should be func, actual: %v", ftype.Kind())
	}
	injectMap := injectMapPool.Get().(map[string]interface{})
	instance := t.container.GetInstance(instanceName, injectMap)
	numsIn := ftype.NumIn()
	inParams := make([]reflect.Value, numsIn)
	var i = 0
	for i < numsIn {
		inType := ftype.In(i)
		inKind := inType.Kind()
		switch inKind {
		case reflect.Interface:
			if contextType.Implements(inType) {
				inParams[i] = reflect.ValueOf(context.Background())
			} else if reflect.TypeOf(instance).Implements(inType) {
				inParams[i] = reflect.ValueOf(instance)
			} else {
				return fmt.Errorf("failed to inject param")
			}
		default:
			return fmt.Errorf("failed to inject param")
		}
		i++
	}
	responseValue := reflect.ValueOf(f).Call(inParams)
	if len(responseValue) != 1 {
		return fmt.Errorf("func response has to be one")
	}
	if err, ok := responseValue[0].Interface().(error); ok {
		if err != nil {
			return err
		}
	}
	return nil
}

func IsFunc(handlerType reflect.Type) bool {
	return handlerType.Kind() == reflect.Func
}
