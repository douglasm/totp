# totp

Time based one time password demo app. Feel free to use and modify to your heart's content.

### Description

The app will allow users to sign up, login and validate using a time based one time password. To run the app use the command
```
go run main.go
```
All the signups are held in a map which gets deleted every time the app restarts. The app listens on port 62222, ie localhost:62222

There are four pages:  
- /
- /signup
- /qrcode
- /validate

#### /

The home page displayes a standard login page with a signup button. The signup button takes you to /signup

#### /signup

The signup page allows you to create a new user by entering an email address and a password. The email's validity isn't checked, but the password is 
checked against the HaveIBeenPwned breached passwords list. Successful signup will take you to the /qrcode page.

#### /qrcode

The qrcode page displays the QRCode used in your authenticator app. The QRCode is created using the domain, your email and a secret key. After
scanning the code you go to the /validate page.

#### /validate

The validate page allows you to enter the time based code from the authenticator app and check whether it is valid or not.
