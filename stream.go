package stream

import (
	"reflect"
	"sync"
)

func Each(slice, iterate interface{}) {
	sliceValue, iterateValue := reflect.ValueOf(slice), reflect.ValueOf(iterate)
	wait := sync.WaitGroup{}
	length := sliceValue.Len()
	wait.Add(length)
	for i := 0; i < length; i++ {
		go func(i int) {
			iterateValue.Call([]reflect.Value{sliceValue.Index(i), reflect.ValueOf(i)})
			wait.Done()
		}(i)
	}
	wait.Wait()
}

func Map(slice, iterate interface{}) interface{} {
	sliceValue := reflect.ValueOf(slice)
	iterateValue := reflect.ValueOf(iterate)
	wait := sync.WaitGroup{}
	length := sliceValue.Len()
	wait.Add(length)
	result := reflect.MakeSlice(reflect.SliceOf(iterateValue.Type().Out(0)), length, length)
	for i := 0; i < length; i++ {
		go func(i int) {
			arguments := []reflect.Value{sliceValue.Index(i), reflect.ValueOf(i)}
			callResult := iterateValue.Call(arguments)[0]
			result.Index(i).Set(callResult)
			wait.Done()
		}(i)
	}
	wait.Wait()
	return result.Interface()
}

func Filter(slice, iterate interface{}) interface{} {
	sliceValue := reflect.ValueOf(slice)
	iterateValue := reflect.ValueOf(iterate)
	wait := sync.WaitGroup{}
	length := sliceValue.Len()
	wait.Add(length)
	result := reflect.MakeSlice(reflect.SliceOf(iterateValue.Type().In(0)), 0, length/2)
	for i := 0; i < length; i++ {
		go func(i int) {
			iValue := sliceValue.Index(i)
			arguments := []reflect.Value{iValue, reflect.ValueOf(i)}
			shouldBeAdded := iterateValue.Call(arguments)[0].Bool()
			if shouldBeAdded {
				result = reflect.Append(result, iValue)
			}
			wait.Done()
		}(i)
	}
	wait.Wait()
	return result.Interface()
}

func Every(slice, iterate interface{}) (result bool) {
	sliceValue := reflect.ValueOf(slice)
	iterateValue := reflect.ValueOf(iterate)
	wait := sync.WaitGroup{}
	length := sliceValue.Len()
	wait.Add(length)
	result = true
	for i := 0; i < length; i++ {
		go func(i int) {
			defer wait.Done()
			if !result {
				return
			}
			iValue := sliceValue.Index(i)
			arguments := []reflect.Value{iValue, reflect.ValueOf(i)}
			isOk := iterateValue.Call(arguments)[0].Bool()
			if !isOk {
				result = false
			}
		}(i)
	}
	wait.Wait()
	return
}

func Some(slice, iterate interface{}) (result bool) {
	sliceValue := reflect.ValueOf(slice)
	iterateValue := reflect.ValueOf(iterate)
	wait := sync.WaitGroup{}
	length := sliceValue.Len()
	wait.Add(length)
	result = false
	for i := 0; i < length; i++ {
		go func(i int) {
			defer wait.Done()
			if result {
				return
			}
			iValue := sliceValue.Index(i)
			arguments := []reflect.Value{iValue, reflect.ValueOf(i)}
			isOk := iterateValue.Call(arguments)[0].Bool()
			if isOk {
				result = true
			}
		}(i)
	}
	wait.Wait()
	return
}
