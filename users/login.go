package users

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"

	"totp/types"
)

type (
	LoginFormRecord struct {
		Name     string `schema:"username"`
		Password string `schema:"password"`
		Remember bool   `schema:"remember"`
		Commit   string `schema:"commit"`
	}

	signupRecord struct {
		Name      string `schema:"name"`
		Password1 string `schema:"pass1"`
		Password2 string `schema:"pass2"`
		Commit    string `schema:"commit"`
	}
)

// LoggedIn checks to see if a user is already logged in
func LoggedIn(w http.ResponseWriter, r *http.Request, theUser *types.LoginRec) bool {
	c, err := r.Cookie("session")
	if err != nil {
		DisplayWelcome(w, r)
		return false
	}

	_, ok := logins[c.Value]
	if !ok {
		return false
	}

	return true
}

// Signup gets the email address and password, creates the cookie and goes to the QRCode page
func Signup(w http.ResponseWriter, r *http.Request) {
	var (
		theSignup    signupRecord
		header       types.HeaderRecord
		errorMessage string
		letters      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	)

	switch r.Method {
	case http.MethodGet:
		// a get does nothing other than present the form

	case http.MethodPost:
		// get the details entered in the form
		err := r.ParseForm()
		err = decoder.Decode(&theSignup, r.PostForm)
		if err != nil {
			fmt.Println("Decode schema error", err)
		}
		fmt.Printf("%+v\r\n", theSignup)

		// Do the passwords match?
		if theSignup.Password1 != theSignup.Password2 {
			errorMessage = "Mismatched password"
			break
		}
		// check to see if the password has been found anywhere else
		err = checkPasswordSafe(theSignup.Password1)
		if err != nil {
			errorMessage = err.Error()
			break
		}

		// make a value for the session cookie
		rand.Seed(time.Now().UnixNano())
		tempStr := ""
		for i := 0; i < 16; i++ {
			num := rand.Intn(62)
			tempStr += string(letters[num])
		}

		logins[tempStr] = types.LoginRec{Mail: theSignup.Name, Pass: theSignup.Password1, Cookie: tempStr}

		cookie := &http.Cookie{
			Name:  "session",
			Value: tempStr,
			Path: "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		// Go to the page where the QRCode is displayed
		http.Redirect(w, r, "/qrcode", http.StatusFound)
		return
	}

	header.Title = "Edinburgh Go signup"
	data := struct {
		Header       types.HeaderRecord
		Name         string
		ErrorMessage string
		Menu         []types.NavItem
		Logout       bool
	}{
		header,
		theSignup.Name,
		errorMessage,
		types.GeneralMenu,
		true,
	}

	t, err := template.ParseFiles("templates/signup.html", types.ViewHeader, types.ViewMenuConstant)
	if err != nil {
		fmt.Println("login form err", err)
	}
	t.Execute(w, data)
}

// QRCode displays the totp image which can be scanned for the authenticator program
func QRCode(w http.ResponseWriter, r *http.Request) {
	var (
		header types.HeaderRecord
	)
	c, err := r.Cookie("session")
	if err != nil {
		DisplayWelcome(w, r)
		return
	}

	// Get the user using the cookie value
	theUser := logins[c.Value]
	fmt.Printf("%+v\r\n", theUser)

	// Generate the QR Code image
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "lordprotector.uk",
		AccountName: theUser.Mail,
	})
	// save the secret, in the real world this should be stored encrypted
	theUser.Secret = key.Secret()

	// add the user to the list of users
	logins[c.Value] = theUser

	// save the QRCode as a png
	f, err := os.Create("img/" + c.Value + ".png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, err := key.Image(200, 200)
	png.Encode(f, img)

	header.Title = "Edinburgh Go signup"
	data := struct {
		Header types.HeaderRecord
		Name   string
		Image  string
		Menu   []types.NavItem
		Logout bool
	}{
		header,
		theUser.Mail,
		c.Value,
		types.GeneralMenu,
		true,
	}

	t, err := template.ParseFiles("templates/qrcode.html", types.ViewHeader, types.ViewMenuConstant)
	if err != nil {
		fmt.Println("qrcode form err", err)
	}
	t.Execute(w, data)
}

// Validate asks for a code using the authenticator app
func Validate(w http.ResponseWriter, r *http.Request) {
	var (
		theCode      signupRecord
		header       types.HeaderRecord
		errorMessage string
	)

	c, err := r.Cookie("session")
	if err != nil {
		DisplayWelcome(w, r)
		return
	}

	// find the user's details
	theUser := logins[c.Value]
	fmt.Printf("Validate %+v\r\n", theUser)

	switch r.Method {
	case http.MethodGet:

	case http.MethodPost:
		// get the code they have entered
		err := r.ParseForm()
		err = decoder.Decode(&theCode, r.PostForm)
		if err != nil {
			fmt.Println("Decode schema error", err)
		}
		fmt.Printf("%+v\r\n", theCode)

		// validate the code using the code and the secret
		valid := totp.Validate(theCode.Name, theUser.Secret)
		fmt.Println("Validation:", valid)
		if !valid {
			errorMessage = "Wrong code"
			break
		}

		// if the code is good mark them as logged in
		theUser.LoggedIn = true
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	header.Title = "Edinburgh Go validate"
	data := struct {
		Header       types.HeaderRecord
		Name         string
		ErrorMessage string
		Menu         []types.NavItem
		Logout       bool
	}{
		header,
		theUser.Mail,
		errorMessage,
		types.GeneralMenu,
		true,
	}

	t, err := template.ParseFiles("templates/validate.html", types.ViewHeader, types.ViewMenuConstant)
	if err != nil {
		fmt.Println("login form err", err)
	}
	t.Execute(w, data)
}

func DisplayWelcome(w http.ResponseWriter, r *http.Request) {
	// start with a blank form
	OutputLoginForm("", "", false, w, r)
}

func OutputLoginForm(userName string, errorMessage string, remember bool, w http.ResponseWriter, r *http.Request) {
	var (
		header types.HeaderRecord
	)
	header.Title = "Edinburgh Go Login"
	data := struct {
		Header       types.HeaderRecord
		Name         string
		ErrorMessage string
		Remember     bool
		Menu         []types.NavItem
		Logout       bool
	}{
		header,
		userName,
		errorMessage,
		remember,
		types.GeneralMenu,
		true,
	}

	t, err := template.ParseFiles("templates/login.html", types.ViewHeader, types.ViewMenuConstant)
	if err != nil {
		fmt.Println("login form err", err)
	}
	t.Execute(w, data)
}

// checkPasswordSafe contacts haveibeenpwned and gets a list of known password with the same start
func checkPasswordSafe(password string) error {
	var (
		resp *http.Response
	)

	// create a new sha1 hash of the password
	h := sha1.New()
	io.WriteString(h, password)
	theStr := fmt.Sprintf("%X", h.Sum(nil))

	// cal the hibp api with the first 5 characters of the hash
	theURL := "https://api.pwnedpasswords.com/range/" + theStr[:5]
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodGet, theURL, nil)
	if err != nil {
		fmt.Println("NewRequest Error:", err)
		return err
	}

	// try for three times to get a response
	for i := 0; i < 3; i++ {
		resp, err = client.Do(r)
		if err == nil {
			break
		}
		i++
	}

	if err != nil {
		fmt.Println("Error: Pass check client", err)
		return err
	}

	// get the returned text and see if the suffix of the password's hash
	// appears in the list. If it does reject the password
	body, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(body), theStr[5:]) {
		return errPassUsed
	}

	return nil
}
