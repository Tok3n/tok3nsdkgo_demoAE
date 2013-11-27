package  tok3nsdkgoDemoAE

import (
	"net/http"
	"fmt"
	"os"
	"bufio"
	"strings"

	"appengine"
	"appengine/datastore"
	"github.com/gorilla/sessions"
	"github.com/Tok3n/tok3nsdkgo"
)

var sessionsStore = sessions.NewCookieStore([]byte("What ever you feel secure"))
var tok3nConfig = tok3nsdkgo.GetTok3nConfigWithSecretPublic("a76f7b1e-c741-5c2d-5a35-cda0d08ddda7","c0c0fa82-6e4f-5040-5004-262c22acfa48")

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

	c := appengine.NewContext(r)
	c.Infof("The Value is: %v",value)
	
	var u User

    if (value == nil){
    	c.Infof("Redirecting to login")
    	http.Redirect(w, r, "/login.do", http.StatusTemporaryRedirect)
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
		t := q.Run(c)
		t.Next(&u)
    }

    return &u
}

func validateSecureAccessWithTok3n(w http.ResponseWriter, r *http.Request, u *User){
	session, _ := sessionsStore.Get(r, "logindata")
	tok3n := session.Values["tok3niced"]
	if u.Tok3nKey != "" && tok3n == nil{
			http.Redirect(w, r, "/login.tok3n", http.StatusTemporaryRedirect)
	}
}

var mydomain = "http://thegoapp.appspot.com"

func init() {


	http.HandleFunc("/_ah/warmup",warmup_method) //usefull for initing changes betwen versions
    http.HandleFunc("/", rootWS)

    registerLoginFunctions()
    registerTok3nFunctions()
}

/**
Do nothing just print ok
**/
func warmup_method(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w,"OK")
}

func rootWS(w http.ResponseWriter, r *http.Request) {
	usr := secureWebAccess(w,r)
	validateSecureAccessWithTok3n(w,r,usr)
	c := appengine.NewContext(r)
	tok3nInstance := tok3nsdkgo.GetAppEngineTok3nInstance(c,tok3nConfig)
	c.Infof("aver: %v\n",tok3nInstance)
	fmt.Fprintf(w, "<html>")
	if usr.Tok3nKey == ""{
		accessurl, err := tok3nInstance.GetAccessUrl(fmt.Sprintf("%s/tok3nreturn",mydomain),usr.Username)
		c.Infof("tok3n access url %s. user %v",accessurl,usr)
		if err != nil{
			c.Infof("Error getting the Tok3n Access Url: '%v'",err)
		}else{
			s := []string{"<div>Did you want to add more security to your account. <a href='", accessurl, "'>YES PLEAE!!!!</a></div>"}
			//addedtext = strings.Join(s,"")
			fmt.Fprint(w, strings.Join(s,""))
		}
	}
	/*
	

	addedtext := ""
	/*if usr.Tok3nKey == ""{
		accessurl, err := tok3nInstance.GetAccessUrl(fmt.Sprintf("%s/tok3ncallback",mydomain),usr.Username)
		if err != nil{
			c.Infof("Error getting the Tok3n Access Url: '%v'",err)
		}else{
			addedtext = fmt.Sprintf("<div>Did you want to add more security to your account. <a href='%s'>YES PLEAE!!!!</a></div>",accessurl)
		}
	}*/
	
	resp := fmt.Sprintf("<br />Here comes the service</html>")
	fmt.Fprintf(w, resp)
}