package main

import (
	"fmt"
	//"os"
	"reflect"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	instance *application.App
}

func NewApp() *App {
	return &App{}
}

// In v3, Services can have an OnStartup method that receives the app instance
func (a *App) OnStartup(ctx *application.Context) {
    // In some alpha versions of v3, you just use the context or:
    a.instance = application.Get() 
}

// Your "Super Vanilla" Dispatcher (same logic as v2)
func (a *App) CallGo(funcName string, args []interface{}, types []string) (interface{}, error) {
    method := reflect.ValueOf(a).MethodByName(funcName)
    if !method.IsValid() {
        return nil, fmt.Errorf("Backend Error: Method '%s' not found", funcName)
    }

    // 1. Validate array lengths
    if len(args) != len(types) {
        return nil, fmt.Errorf("Type Error: Args (%d) and Types (%d) count mismatch", len(args), len(types))
    }

    in := make([]reflect.Value, len(args))
    for i := 0; i < len(args); i++ {
        switch types[i] {
        case "int":
            // JavaScript sends all numbers as float64. 
            // We must handle the type conversion carefully.
            if val, ok := args[i].(float64); ok {
                in[i] = reflect.ValueOf(int(val))
            } else if val, ok := args[i].(int); ok {
                in[i] = reflect.ValueOf(val)
            } else {
                return nil, fmt.Errorf("Position [%d]: Expected number, got %T", i, args[i])
            }
        case "string":
            if val, ok := args[i].(string); ok {
                in[i] = reflect.ValueOf(val)
            } else {
                return nil, fmt.Errorf("Position [%d]: Expected 'string', got %T", i, args[i])
            }
        case "bool":
            if val, ok := args[i].(bool); ok {
                in[i] = reflect.ValueOf(val)
            } else {
                return nil, fmt.Errorf("Position [%d]: Expected 'bool', got %T", i, args[i])
            }
        case "[]string":
            if rawSlice, ok := args[i].([]interface{}); ok {
                stringSlice := make([]string, len(rawSlice))
                for j, v := range rawSlice {
                    stringSlice[j] = v.(string)
                }
                in[i] = reflect.ValueOf(stringSlice)
            } else {
                return nil, fmt.Errorf("Position [%d]: Expected array of strings, got %T", i, args[i])
            }
        default:
            return nil, fmt.Errorf("Position [%d]: Unsupported type '%s'", i, types[i])
        }
    }

    // 2. Invoke the Go Method
    results := method.Call(in)

    // 3. Handle Return Values Safely
    if len(results) == 0 {
        return nil, nil
    }

    // If Go returns (result, error)
    if len(results) == 2 {
        // Safe check for the error return: is it actually an error and not nil?
        if !results[1].IsNil() {
            return nil, results[1].Interface().(error)
        }
        return results[0].Interface(), nil
    }

    // If Go returns a single value
    return results[0].Interface(), nil
}


// --- Your Logic Commands ---

func (a *App) SayHello(name string) string {
	return "Hello " + name
}

