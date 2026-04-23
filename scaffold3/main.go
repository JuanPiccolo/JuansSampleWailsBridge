package main

import (
	"embed"
	_ "embed"
	"log"
	//"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	// application.RegisterEvent[string]("time")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// 1. Read the icon data
	iconData, _ := assets.ReadFile("frontend/dist/icon.png")

	// 2. Initialize the Application with the Icon here
	app := application.New(application.Options{
		Name:        "Scaffold",
		Description: "A demo of using raw HTML & CSS",
		Icon:        iconData, // <--- Move it here
		Services: []application.Service{
			application.NewService(NewApp()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// 3. Create the window without the Icon field
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Main Window",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		URL: "/",
	})

	// 4. Run the app
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}