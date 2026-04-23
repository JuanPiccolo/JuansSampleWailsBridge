package main

import (
    "fmt"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
    instance *application.App
}

func NewApp() *App {
    return &App{}
}

func (a *App) OnStartup(ctx *application.Context) {
    a.instance = application.Get() 
}

// --- The Dispatcher ---
// This replaces your reflect logic. It's explicit, safe, and fast.
func (a *App) CallGo(cmd string, data map[string]any) (any, error) {
    switch cmd {
    case "hello":
        name, _ := data["name"].(string)
        return a.SayHello(name), nil
        
    case "test":
        return "Bridge connection successful!", nil

    default:
        return nil, fmt.Errorf("command '%s' not recognized", cmd)
    }
}

// --- Your Business Logic ---
// Now these are just normal Go functions. No reflect needed.
func (a *App) SayHello(name string) string {
    if name == "" {
        return "Hello, mystery person!"
    }
    return "Hello, " + name + "!"
}
