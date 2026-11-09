package enums

type LoginAction string

const (
	LoginActionRegister LoginAction = "register"
	LoginActionLogin    LoginAction = "login"
	LoginActionLogout   LoginAction = "logout"
	LoginActionFail     LoginAction = "fail"
	LoginActionLocked   LoginAction = "locked"
)
