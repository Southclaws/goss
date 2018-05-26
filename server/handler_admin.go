package server

import (
	"encoding/json"
	"io"

	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/Southclaws/ScavengeSurviveCore/types"
)

func (app *App) adminRoutes() []Route {
	return []Route{
		{
			"adminSet",
			"POST",
			"/store/adminSet",
			true,
			types.ExampleAdmin(),
			types.ExampleStatus(nil, true),
			app.adminSet,
		},
		{
			"adminGetList",
			"GET",
			"/store/adminGetList",
			false,
			nil,
			[]types.Admin{types.ExampleAdmin(), types.ExampleAdmin()},
			app.adminGetList,
		},
	}
}

func (app *App) adminSet(r io.Reader) (status types.Status, err error) {
	var admin types.Admin
	err = json.NewDecoder(r).Decode(&admin)
	if err != nil {
		return
	}
	err = app.validator.Struct(admin)
	if err != nil {
		return types.NewStatusValidationError(err.(validator.ValidationErrors)), nil
	}

	logger.Debug("received request adminSet",
		zap.Any("admin", admin))

	return types.NewStatus(nil, true, ""), app.store.AdminSetLevel(admin.PlayerID, *admin.Level)
}

func (app *App) adminGetList(r io.Reader) (status types.Status, err error) {
	list, err := app.store.AdminGetList()
	return types.NewStatus(list, true, ""), err
}
