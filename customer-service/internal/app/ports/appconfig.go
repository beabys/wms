// Package ports holds shared interfaces to avoid import cycles
// when interfaces reference types from their own package.
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

// AppConfig is the configuration interface for the customer service.
type AppConfig interface {
	LoadConfig(path string) error
	GetGRPCPort() int
	GetDBDSN() string
	GetLogLevel() string
	GetJWTPublicKeyPath() string
}
