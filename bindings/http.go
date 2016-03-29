package bindings

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"math/big"

	"golang.org/x/net/context"

	"github.com/davecgh/go-spew/spew"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/jmc-audio/kitauth/auth"
	"github.com/jmc-audio/kitauth/consts"
	"github.com/jmc-audio/kitauth/log"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

type Servicer interface {
	Run(context.Context, interface{}) (interface{}, error)
}

type Request struct {
	params map[string]string
}

type Response struct {
	Status string
}

// Example of using a gorilla-mux path variable as an auth.Principal
func (r *Request) PrincipalToken() interface{} {
	if id, ok := r.params[consts.RequestPrincipalID]; ok {
		return &id
	}
	return nil
}

// Example of using HTTP Request Method as an auth.Subject
func (r *Request) SubjectTokens() []interface{} {
	if method, ok := r.params["HTTP_REQUEST_METHOD"]; ok {
		return []interface{}{&method}
	}
	return []interface{}{}
}

// Put a gorilla-mux path var into Go-Kit's Request scope for use an auth.Principal
// Put the HTTP Request Method into Go-Kit's Request scope for use as an auth.Subject
func decodeRequest(r *http.Request) (response interface{}, err error) {
	var (
		id string
		ok bool
	)
	urlParams := mux.Vars(r)

	if id, ok = urlParams[consts.RequestPrincipalID]; !ok {
		return nil, errors.New("No principal id in request")
	}
	return &Request{map[string]string{consts.RequestPrincipalID: id, "HTTP_REQUEST_METHOD": r.Method}}, nil
}

func encodeResponse(w http.ResponseWriter, i interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(i.(*Response))
}

type Endpoint struct{}

func NewEndpoint(ctx context.Context) Servicer {
	return &Endpoint{}
}

func (h *Endpoint) Run(ctx context.Context, i interface{}) (interface{}, error) {
	log.Debug(ctx, "ctx", spew.Sdump(ctx), "i", spew.Sdump(i))
	return &Response{"OK"}, nil
}

func StartHTTPListener(root context.Context) {
	go func() {
		ctx, cancel := context.WithCancel(root)
		defer cancel()

		errc := ctx.Value(consts.ContextErrorChannel).(chan error)

		router := createRouter(ctx, NewEndpoint(ctx))
		errc <- http.ListenAndServe(":6502", handlers.CombinedLoggingHandler(os.Stderr, router))
	}()
}

func createRouter(ctx context.Context, endpoint Servicer) *mux.Router {
	router := mux.NewRouter()

	// An example authn function that requires an auth.Principal token to (probably be a prime number
	authn := func(p auth.Principal) bool {
		if p == nil {
			return false
		}
		if p.PrincipalToken() == nil {
			return false
		}
		if token, ok := p.PrincipalToken().(*string); ok {
			log.Debug(ctx, consts.RequestPrincipalID, token)
			i, err := strconv.ParseInt(*token, 10, 64)
			if err != nil {
				log.Error(ctx, "principal", *token, "err", err.Error())
				return false
			}
			if big.NewInt(i).ProbablyPrime(20) {
				log.Error(ctx, "principal", *token, "probably_prime", true)
				return true
			} else {
				log.Error(ctx, "principal", *token, "probably_prime", false)
				return false
			}
		}
		return false
	}

	// An example authz function that requires (any) auth.Principal
	// and allows only HTTP GET represented as an auth.Subject
	authz := func(subjects auth.Subject) bool {
		if subjects.PrincipalToken() == nil {
			return false
		}
		if subjects == nil {
			return true
		}
		for _, s := range subjects.SubjectTokens() {
			if s != nil {
				if token, ok := s.(*string); ok {
					log.Debug(ctx, "subject", token)
					if *token == "GET" {
						return true
					}
				}
			}
		}
		return false
	}

	authenticator := auth.NewAuthenticator(authn, authz)

	Authenticated := authenticator.Authenticated()
	Authorized := authenticator.Authorized()

	router.Handle(fmt.Sprintf("/principal/{%s}", consts.RequestPrincipalID),
		kithttp.NewServer(
			ctx,
			Authorized(
				Authenticated(endpoint.Run)),
			decodeRequest,
			encodeResponse,
		)).Methods("GET", "POST")
	return router
}
