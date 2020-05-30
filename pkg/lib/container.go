package lib

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	TypeNotRegisteredError = errors.New("TypeNotRegistered")
	CircularReferenceError = errors.New("CircularReferenceDetected")
)

type threadSafeContainer struct {
	newFunctionRegistry map[reflect.Type]*NewFunction
	singletonRegistry   map[reflect.Type]interface{}
	mu                  *sync.RWMutex
}

func newThreadSafeContainer() *threadSafeContainer {
	return &threadSafeContainer{
		newFunctionRegistry: map[reflect.Type]*NewFunction{},
		singletonRegistry:   map[reflect.Type]interface{}{},
		mu:                  &sync.RWMutex{},
	}
}

type NewFunction struct {
	F         func() (interface{}, error)
	Singleton bool
}

func (c *threadSafeContainer) Provide(t interface{}) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	typeNeeded := reflect.TypeOf(t).Elem()
	return c.provide(typeNeeded, map[reflect.Type]string{})
}

const (
	Visiting string = "visiting"
	Done     string = "done"
)

func (c *threadSafeContainer) provide(typeNeeded reflect.Type, visited map[reflect.Type]string) (interface{}, error) {
	if ret, ok := c.singletonRegistry[typeNeeded]; ok {
		return ret, nil
	}

	if visited[typeNeeded] == Visiting {
		return nil, CircularReferenceError
	}

	visited[typeNeeded] = Visiting

	if newFunction, ok := c.newFunctionRegistry[typeNeeded]; ok {
		ret, err := newFunction.F()
		if err != nil {
			return nil, err
		}
		retType := reflect.TypeOf(ret).Elem()
		retValue := reflect.ValueOf(ret).Elem()
		for i := 0; i < retValue.NumField(); i++ {
			fieldType := retValue.Field(i).Type().Elem()
			// fmt.Println("get field type", retType.Name(), fieldType.Name(), retValue.Field(i).CanAddr(), retValue.Field(i).CanSet())
			autowired := retType.Field(i).Tag.Get("autowired")
			if autowired == "true" {
				fieldValue, err := c.provide(fieldType, visited)
				if err != nil {
					return nil, err
				}

				if retValue.Field(i).CanSet() {
					retValue.Field(i).Set(reflect.ValueOf(fieldValue))
				}
			}
		}

		if newFunction.Singleton {
			c.singletonRegistry[typeNeeded] = ret
		}

		visited[typeNeeded] = Done
		return ret, nil
	} else {
		return nil, TypeNotRegisteredError
	}
}

func (c *threadSafeContainer) RegisterType(t interface{}, f *NewFunction) {
	if f == nil || f.F == nil {
		panic(fmt.Sprintf("register type %v with invalid function\n", t))
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	typeToRegister := reflect.TypeOf(t).Elem()
	if _, ok := c.newFunctionRegistry[typeToRegister]; !ok {
		c.newFunctionRegistry[typeToRegister] = f
	} else {
		panic(fmt.Sprintf("type: %s already registered\n", typeToRegister))
	}
}

var defaultContainer = newThreadSafeContainer()

func Provide(t interface{}) (interface{}, error) {
	return defaultContainer.Provide(t)
}

func RegisterType(t interface{}, f *NewFunction) {
	defaultContainer.RegisterType(t, f)
}
