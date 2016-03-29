package auth

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

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
