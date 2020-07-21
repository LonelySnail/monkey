package monkey

import (
	"github.com/LonelySnail/monkey/app"
	"github.com/LonelySnail/monkey/module"
)

const version = "1.0.0"

func NewDefaultApp(opts ...app.OptionFn) module.IDefaultApp {
	opts = append(opts, app.SetAppVersion(version))
	return app.NewDefaultApp(opts...)
}
