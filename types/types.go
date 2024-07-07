package types

import (
)

const (
	ViewMenu         = "templates/menu.html"
	ViewHeader       = "templates/header.html"
	ViewErr          = "templates/error.html"
	ViewNavbar       = "templates/navbar.html"
	ViewMenuConstant = "templates/navbar_constant.html"
	ViewMenuUs       = "templates/navbarus.html"
	ViewNavButtons   = "templates/navbuttons.html"
	KListLimit = 20

	KLoginFormID = 789
)

type (
	HeaderRecord struct {
		Title          string
		Angular        bool
		JSEditor       bool
		FooTable       bool
		DatePicker     bool
		DocReady       bool
		ChatSizeDelete bool
		NewStyle       bool
		Moderator      bool
		Moderating     bool
		HasScript1     bool
		Script1        string
		HasScript2     bool
		Script2        string
		HasScript3     bool
		Script3        string
	}

	LoginRec struct {
		Mail     string
		Pass     string
		Cookie   string
		Secret   string
		LoggedIn bool
	}

	NavItem struct {
		Text  string
		Title string
		Link  string
	}
)

var (
	GeneralMenu = []NavItem{
		{"Index", "/", "Go to the Agamik home page"},
		{"Programs", "/barcoder", "Information about our barcode programs, download the latest release"},
		{"Fonts", "/fonts", "Information about barcode fonts, download the latest release"},
		{"Creation", "/create", "We can supply your barcodes as files"},
		{"Downloads", "/download", "Download working versions and demos of our products"},
		{"Types", "/symbols", "Information about barcode types and how to identify different types"},
		{"Explained", "/explain", "Answers to the common questions. Information about barcoding. What types to use for which jobs"},
		{"Buying", "/buying", "Information about buying our products. How to buy and how to pay"},
		{"Contact", "/contact", "Contact information: e-mail, phone, mail, Skype, MSN Messenger and Yahoo Messenger details"},
	}
)
