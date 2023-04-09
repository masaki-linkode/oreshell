package myvariables

import (
	"oreshell/infra"
	"oreshell/log"
)

var shellVariables = map[string]string{}

func Variables() variables {
	return variables{
		osService:      infra.MyOSService{},
		shellVariables: shellVariables,
	}
}

type variables struct {
	osService      infra.OSService
	shellVariables map[string]string
}

func (me variables) GetValueFromShellVariables(variable_name string) (string, bool) {
	value, ok := me.shellVariables[variable_name]
	return value, ok
}

func (me variables) GetValue(variable_name string) string {
	value := me.osService.Getenv(variable_name)
	if len(value) > 0 {
		return value
	}
	return me.shellVariables[variable_name]
}

type kv struct {
	VariableName string
	Value        string
}

func GetKVIterator() <-chan kv {
	ch := make(chan kv)
	go func() {
		for k, val := range shellVariables {
			ch <- kv{
				VariableName: k,
				Value:        val,
			}
		}
		close(ch)
	}()
	return ch
}

func (me variables) AssignValueToShellVariable(variable_name string, value string) error {
	log.Logger.Printf("AssignValueToShellVariable start: %+v, %+v\n", variable_name, value)
	if me.osService.Hasenv(variable_name) {
		return me.osService.Setenv(variable_name, value)
	} else {
		if len(value) > 0 {
			me.shellVariables[variable_name] = value
		} else {
			log.Logger.Printf("AssignValueToShellVariable map delete: %+v\n", variable_name)
			delete(me.shellVariables, variable_name)
		}
		return nil
	}
}

func (me variables) AssignValuesToShellVariables(variables map[string]string) error {
	log.Logger.Printf("AssignValuesToShellVariables start: %+v\n", variables)
	for variable_name, value := range variables {
		err := me.AssignValueToShellVariable(variable_name, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (me variables) AssignValueToEnvironmentVariable(variable_name string, value string) error {
	return me.osService.Setenv(variable_name, value)
}
