package auth

type UnknownPrincipal struct{}

func (UnknownPrincipal) Error() string {
	return "Unknown principal"
}

type Unauthenticated struct{}

func (Unauthenticated) Error() string {
	return "Unauthenticated"
}

type UnknownSubject struct{}

func (UnknownSubject) Error() string {
	return "Unknown subject"
}

type Unauthorized struct{}

func (Unauthorized) Error() string {
	return "Unauthorized"
}
