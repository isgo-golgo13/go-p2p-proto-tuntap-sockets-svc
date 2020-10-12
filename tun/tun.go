package tun

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/songgao/water"
)

type IfcCtx struct {
	Ifc *water.Interface
	Err error
	Ctx context.Context
}

func NewIfcCtx(tunDev string) *IfcCtx {
	ctx := &IfcCtx{}
	ctx.Ifc, ctx.Err = water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: tunDev,
		},
	})
	if ctx.Err != nil {
		log.Fatalf("%v", ctx.Err)
		return nil
	}
	return ctx
}
