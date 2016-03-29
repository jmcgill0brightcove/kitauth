package auth

type UnknownPrincipal struct{}

func (UnknownPrincipal) Error() string {
	return "Unknown principal"
}

type Unauthenticated struct{}

func (Unauthenticated) Error() string {
	return "Unauthenticated"
}

type Principal interface {
	PrincipalToken() *string
}

type Subject interface {
	SubjectToken() *string
}

type AuthNFunc func(p Principal) bool

type AuthZFunc func(p Principal, s []Subject) []Subject
