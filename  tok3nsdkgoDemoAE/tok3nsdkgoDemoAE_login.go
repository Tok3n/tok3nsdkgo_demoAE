package tok3nsdkgoDemoAE
import (
		"net/http"
		)

func registerLoginFunctions(){
	http.HandleFunc("/login.do", )
}

func loginDo(w http.ResponseWriter, r *http.Request)
{
	fmt.Fprintf(w,"try login")
}