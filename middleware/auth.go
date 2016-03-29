package middleware

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/jmc-audio/kitauth/auth"
	"github.com/jmc-audio/kitauth/log"
	"golang.org/x/net/context"
)

type Authenticator interface {
	Authenticated() endpoint.Middleware
	//Authorized() endpoint.Middleware
}

type authenticator struct {
	authN auth.AuthNFunc
	authZ auth.AuthZFunc
}

func NewAuthenticator(authN auth.AuthNFunc, authZ auth.AuthZFunc) Authenticator {
	return &authenticator{authN: authN, authZ: authZ}
}

func (a *authenticator) Authenticated() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, i interface{}) (interface{}, error) {
			var ok bool
			var p auth.Principal
			if p, ok = i.(auth.Principal); ok {
				log.Debug(ctx, "principal", p.PrincipalToken())
				if p == nil || p.PrincipalToken() == nil {
					return nil, &auth.UnknownPrincipal{}
				}
				if a.authN(p) {
					return next(ctx, i)
				}
				return nil, &auth.Unauthenticated{}
			}
			return func(ctx context.Context, i interface{}) (interface{}, error) {
				return nil, &auth.UnknownPrincipal{}
			}(ctx, i)
		}
	}
}
