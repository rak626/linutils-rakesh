package config

import (
	"github.com/rak626/linutils-rakesh/internal/pkgmanager"
)

type EnvConfigurator interface {
	Setup(manager pkgmanager.PackageManager) error
}
