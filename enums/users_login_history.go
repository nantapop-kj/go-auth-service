package enums

type LoginAction string

const (
	LoginActionLogin  LoginAction = "login"
	LoginActionLogout LoginAction = "logout"
	LoginActionFail   LoginAction = "fail"
)
