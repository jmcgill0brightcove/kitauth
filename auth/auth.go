package auth

type Principal interface {
	PrincipalToken() interface{}
}

type Subject interface {
	PrincipalToken() interface{}
	SubjectTokens() []interface{}
}

type AuthNFunc func(p Principal) bool

type AuthZFunc func(s Subject) bool
