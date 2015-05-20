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
			result.Index(i).Set(iterateValue.Call([]reflect.Value{sliceValue.Index(i), reflect.ValueOf(i)})[0])
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
	result := reflect.MakeSlice(reflect.SliceOf(iterateValue.Type().In(0)), 0, length)
	for i := 0; i < length; i++ {
		go func(i int) {
			iValue := sliceValue.Index(i)
			shouldBeAdded := iterateValue.Call([]reflect.Value{iValue, reflect.ValueOf(i)})[0].Bool()
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
			iValue := sliceValue.Index(i)
			isOk := iterateValue.Call([]reflect.Value{iValue, reflect.ValueOf(i)})[0].Bool()
			if !isOk {
				result = false
			}
			wait.Done()
		}(i)
	}
	wait.Wait()
	return result
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
			iValue := sliceValue.Index(i)
			isOk := iterateValue.Call([]reflect.Value{iValue, reflect.ValueOf(i)})[0].Bool()
			if isOk {
				result = true
			}
			wait.Done()
		}(i)
	}
	wait.Wait()
	return result
}
