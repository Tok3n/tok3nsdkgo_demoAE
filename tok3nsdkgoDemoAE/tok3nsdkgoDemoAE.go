package  tok3nsdkgoDemoAE

import (
	"net/http"
	"fmt"
	"os"
	"bufio"

	"appengine"
	"appengine/datastore"
	"github.com/gorilla/sessions"
)

var sessionsStore = sessions.NewCookieStore([]byte("What ever you feel secure"))

func ReadString(filename string) string{
    
    var lines string

    f, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
        return ""
    }
    defer f.Close()
    r := bufio.NewReader(f)
    line, err := r.ReadString('\n')
    for err == nil {
        lines += line
        //fmt.Print(line)
        line, err = r.ReadString('\n')
    }
    
    return lines
    
}

func secureWebAccess(w http.ResponseWriter, r *http.Request) *User{
	session, _ := sessionsStore.Get(r, "logindata")
	value := session.Values["id"]

	fmt.Printf("%v",value)
	
	var u User

    if (value == nil){
    	http.Redirect(w, r, "/login.do", http.StatusTemporaryRedirect)
		return nil
	}else{
		c := appengine.NewContext(r)
		q := datastore.NewQuery("User").
			Filter("Username = ",value.(string))
		//k := datastore.NewKey(c,"User",,0,nil)

		if count,_ := q.Count(c);count==0{
			//what ever you have in the session is not in the datastore, these happends when you delete a user directly from the datastore, and the sesion is still active
			http.Redirect(w, r, "/login.do", http.StatusTemporaryRedirect)
			return nil
		}
    }

    return &u
}

func init() {
	http.HandleFunc("/_ah/warmup",warmup_method) //usefull for initing changes betwen versions
    http.HandleFunc("/", rootWS)

    registerLoginFunctions()
}

/**
Do nothing just print ok
**/
func warmup_method(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w,"OK")
}

func rootWS(w http.ResponseWriter, r *http.Request) {
	usr := secureWebAccess(w,r)
	addedtext := ""
	if usr.Tok3nKey == ""{
		addedtext = "<div>Did you want to add more security to your account. <a href=''>YES PLEAE!!!!</a></div>"
	}
	resp := fmt.Sprintf("<html>%v<br />Here comes the service</html>",addedtext)
	fmt.Fprintf(w, resp)
}