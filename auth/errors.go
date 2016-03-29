package auth

import "github.com/go-kit/kit/endpoint"

type Authenticator interface {
	Authenticated() endpoint.Middleware
	//Authorized() endpoint.Middleware
}

type UnknownPrincipal struct{}

func (UnknownPrincipal) Error() string {
	return "Unknown principal"
}

type Unauthenticated struct{}

func (Unauthenticated) Error() string {
	return "Unauthenticated"
}
