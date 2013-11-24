package tok3nsdkgoDemoAE
import (
		"time"
		)

type User struct {
    Creation time.Time
    Username string
    //With these key Tok3n identify the user in your system
    Tok3nKey string 
}