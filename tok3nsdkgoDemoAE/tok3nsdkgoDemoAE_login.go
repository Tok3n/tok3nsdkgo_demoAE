package tok3nsdkgoDemoAE
import (
		"net/http"
		"fmt"
		"html/template"
		"appengine"
		"appengine/datastore"
		"time"

		"github.com/gorilla/sessions"
		)

func registerLoginFunctions(){
	http.HandleFunc("/login.do", loginDo)
	http.HandleFunc("/login.docreate", loginDoCreate)
	http.HandleFunc("/login.verify", loginVerify)
	http.HandleFunc("/login.docreatenew", doCreateNew)
	
	
}

func loginDo(w http.ResponseWriter, r *http.Request){
	str := ReadString("static/templates/login.do.html")
	var homeTemplate = template.Must(template.New("login.do").Parse(str))
    
    if err := homeTemplate.Execute(w, ""); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
func loginDoCreate(w http.ResponseWriter, r *http.Request){
	str := ReadString("static/templates/login.docreate.html")
	var homeTemplate = template.Must(template.New("login.docreate").Parse(str))
    
    if err := homeTemplate.Execute(w, ""); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func loginVerify(w http.ResponseWriter, r *http.Request){
	user := r.FormValue("user")

	if user == "" {
		//empty values or variables not sended. Very basic prevention
		fmt.Fprintf(w, "<html>Error: error in arguments. Go <a href='/login.do' >Back</a></html>")
		return
	}

	c := appengine.NewContext(r)
	
	q := datastore.NewQuery("User").
		Filter("Username = ",user)// Yes we are not ussing passwords (very bad idea but for demo purposes is OK)

	if count,_ := q.Count(c);count>0 { //exists at least one row 
		t := q.Run(c)
		var x User
        _, err := t.Next(&x)
        if err != nil {
            fmt.Fprintf(w,"Error: something happend with the database")
           	return
        }

        session, _ := sessionsStore.Get(r, "logindata")
        session.Options = &sessions.Options{ //just valid for half an hour
		    Path:   "/",
		    MaxAge: 60*30,
		}
		session.Values["id"] = x.Username
		session.Save(r, w)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func doCreateNew(w http.ResponseWriter, r *http.Request){
	user := r.FormValue("user")

	if user == "" {
		//empty values or variables not sended. Very basic prevention
		fmt.Fprintf(w, "<html>Error: error in arguments. Go <a href='/login.do' >Back</a></html>")
		return
	}

	c := appengine.NewContext(r)
	q := datastore.NewQuery("User").
		Filter("Username = ",user)
	if count,_ := q.Count(c);count>0 { //Verify the previous existence of the user
		fmt.Fprintf(w,"<html>Username already exists. <a href='/login.docreate' >Go back</a></html>")
		return
	}

	var theuser User
	theuser.Username = user
	theuser.Creation = time.Now()
	theuser.Tok3nKey = ""

	key := datastore.NewIncompleteKey(c, "User", nil)
	key ,err := datastore.Put(c,key,&theuser)
	if err != nil{
		fmt.Fprintf(w,"Error: %v",err)
		return
	}

	fmt.Fprintf(w, "<html>New user (%s) added. <a href='/login.do'>Go login now</a>.</html>",user)
}