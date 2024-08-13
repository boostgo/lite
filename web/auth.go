package web

type basicAuth struct {
	username string
	password string
}

func initBasicAuth(username, password string) basicAuth {
	return basicAuth{
		username: username,
		password: password,
	}
}
