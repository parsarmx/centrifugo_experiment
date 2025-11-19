package application

import "golang_template/api/handler"

func (a *application) InitHandler() handler.Handler {
	return handler.NewHandler()
}
