package config

import (
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

type EnvConfigurator interface {
	Setup(manager pkgmanager.PackageManager) error
}
