package utils

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/spf13/pflag"
)

var RPCTimeout int64

func IntiliaseLuminoUtils(optionsPackageStruct OptionsPackageStruct) Utils {
	UtilsInterface = optionsPackageStruct.UtilsInterface
	FlagSetInterface = optionsPackageStruct.FlagSetInterface
	return &UtilsStruct{}
}

func (u UtilsStruct) GetUint32(flagSet *pflag.FlagSet, name string) (uint32, error) {
	return flagSet.GetUint32(name)
}

func (f FLagSetStruct) GetLogFileName(flagSet *pflag.FlagSet) (string, error) {
	return flagSet.GetString("logFile")
}

func InvokeFunctionWithTimeout(interfaceName interface{}, methodName string, args ...interface{}) []reflect.Value {
	var functionCall []reflect.Value
	var gotFunction = make(chan bool)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(RPCTimeout)*time.Second)
	defer cancel()

	go func() {
		inputs := make([]reflect.Value, len(args))
		for i := range args {
			inputs[i] = reflect.ValueOf(args[i])
		}
		log.Debug("Blockchain function: ", methodName)
		functionCall = reflect.ValueOf(interfaceName).MethodByName(methodName).Call(inputs)
		gotFunction <- true
	}()
	for {
		select {
		case <-ctx.Done():
			log.Errorf("%s function timeout!", methodName)
			log.Debug("Kindly check your connection")
			return nil

		case <-gotFunction:
			return functionCall
		}
	}
}

func CheckIfAnyError(result []reflect.Value) error {
	if result == nil {
		return errors.New("RPC timeout error")
	}

	errorDataType := reflect.TypeOf((*error)(nil)).Elem()
	errorIndexInReturnedValues := -1

	for i := range result {
		returnedValue := result[i]
		returnedValueDataType := reflect.TypeOf(returnedValue.Interface())
		if returnedValueDataType != nil {
			if returnedValueDataType.Implements(errorDataType) {
				errorIndexInReturnedValues = i
			}
		}
	}
	if errorIndexInReturnedValues == -1 {
		return nil
	}
	returnedError := result[errorIndexInReturnedValues].Interface()
	if returnedError != nil {
		return returnedError.(error)
	}
	return nil
}
