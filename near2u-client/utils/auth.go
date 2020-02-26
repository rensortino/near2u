package utils

func Login(responseMsg, token chan string, email, password string) {

	data := struct {
		Email string
		Password string
	} {
		email,
		password,
	}

	rx := make(chan map[string]interface{})
	go SocketCommunicate("login", "", data, rx)

	res := <- rx // res has type map[string]interface{}

	if res["status"] == "Succesfull" {
		// Accesses nested json
		token <- res["data"].(map[string]interface{})["auth"].(string)
		responseMsg <- "User Authenticated"
	}else {
		token <- "NULL"
		responseMsg <- res["error"].(string)
	}

}

func Register(responseMsg chan string, name, surname, email, password string)  {

	newUser := struct {
		Name string `json:"name"`
		Surname string `json:"surname"`
		Email string `json:"email"`
		Password string `json:"password"`
	} {
		name,
		surname,
		email,
		password,
	}

	rx := make(chan map[string]interface{})

	go SocketCommunicate("register", "", newUser, rx)

	res := <- rx

	if res["status"] == "Succesfull" {
		responseMsg <- "User Registered"
	}else {
		responseMsg <- res["error"].(string)
	}
}
