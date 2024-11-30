package helper

import "github.com/aarioai/AaGo/internal/app"

type Job interface {
	Run(*app.App) error
}

func Run(app *app.App, jobs ...Job) {
	for _, serve := range jobs {
		if err := serve.Run(app); err != nil {
			panic(err)
		}
	}
}
