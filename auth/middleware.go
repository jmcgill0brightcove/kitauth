package auth

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type Authenticator interface {
	Authenticated() endpoint.Middleware
	Authorized() endpoint.Middleware
}

type authenticator struct {
	authN AuthNFunc
	authZ AuthZFunc
}

func NewAuthenticator(authN AuthNFunc, authZ AuthZFunc) Authenticator {
	return &authenticator{authN: authN, authZ: authZ}
}

func (a *authenticator) Authenticated() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, i interface{}) (interface{}, error) {
			var ok bool
			var p Principal

			if p, ok = i.(Principal); ok {
				if p == nil || p.PrincipalToken() == nil {
					return nil, &UnknownPrincipal{}
				}
				if a.authN(p) {
					return next(ctx, i)
				}
				return nil, &Unauthenticated{}
			}
			return func(ctx context.Context, i interface{}) (interface{}, error) {
				return nil, &UnknownPrincipal{}
			}(ctx, i)
		}
	}
}

func (a *authenticator) Authorized() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, i interface{}) (interface{}, error) {
			var ok bool
			var s Subject

			if s, ok = i.(Subject); ok {
				if s == nil {
					return nil, &UnknownSubject{}
				}
				if a.authZ(s) {
					return next(ctx, i)
				}
				return nil, &Unauthorized{}
			}
			return func(ctx context.Context, i interface{}) (interface{}, error) {
				return nil, &UnknownSubject{}
			}(ctx, i)
		}
	}
}
