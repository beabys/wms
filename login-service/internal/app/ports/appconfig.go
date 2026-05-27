// Package ports holds shared interfaces to avoid import cycles.
//
// The pattern: if interface A in package P uses type B from package P
// in its method signature, generating mocks for A requires importing P,
// which creates a cycle when tests in P import those mocks.
//
// Moving such interfaces to a sibling ports/ package breaks the cycle
// while keeping the consumer-defines-interface principle intact.
//
// Use ports/ only when a same-package type appears in an interface
// method signature. Interfaces referencing only external types
// (domain, config, stdlib, etc.) stay in their consumer package.
package ports

import (
	"github.com/beabys/wms/login-service/pkg/config"
)

// AppConfig defines the configuration interface.
type AppConfig interface {
	LoadConfigs() error
	GetConfigs() *config.Config
}
