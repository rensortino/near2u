package utils

type User struct {
	Auth string
	IsAdmin bool
}

func Login(resCh, errCh chan string, loggedUser chan * User, email, password string) {

	data := struct {
		Email    string
		Password string
	}{
		email,
		password,
	}

	res := SocketCommunicate("login", "", data)

	if res["status"] == "Successful" {
		// Accesses nested json
		usr := &User{
			res["data"].(map[string]interface{})["auth"].(string),
			res["data"].(map[string]interface{})["admin"].(bool),
		}
		if usr.IsAdmin {
			resCh <- "Admin Authenticated"
		} else {
			resCh <- "User Authenticated"
		}
		loggedUser <- usr
	} else {
		errCh <- res["error"].(string)
	}

}

func (u * User) Logout(resCh, errCh chan string) {

	res := SocketCommunicate("logout", u.Auth, nil)

	if res["status"] == "Successful" {
		// Accesses nested json
		resCh <- "User Logout"
	} else {
		errCh <- res["error"].(string)
	}
}

func Register(responseMsg chan string, name, surname, email, password string) {

	newUser := struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		name,
		surname,
		email,
		password,
	}

	res := SocketCommunicate("register", "", newUser)

	if res["status"] == "Successful" {
		responseMsg <- "User Registered"
	} else {
		responseMsg <- res["error"].(string)
	}
}
