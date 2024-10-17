package controller

import "go.uber.org/fx"

func RunFx(lc fx.Lifecycle, ctrl Controller) {
	lc.Append(fx.Hook{
		OnStart: ctrl.Start,
		OnStop:  ctrl.Stop,
	})
}
