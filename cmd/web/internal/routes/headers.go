package routes

import (
	"bytes"

	"golang.org/x/exp/slog"
)

func leftNavbar() string {
	var buff bytes.Buffer
	t, ok := templates["navbar-set.html"]
	if !ok {
		slog.Error("Error fetching template: navbar-set")
		return ""
	}
	err := t.Execute(&buff, dataArray{
		dataMap{
			"URL":           "/home",
			"Name":          "Home",
			"ShouldPushURL": true,
			"PushURL":       "/",
		},
	})
	if err != nil {
		slog.Error("Error parsing template: ", err)
		return ""
	}
	return buff.String()
}

func rightNavbar() string {
	var buff bytes.Buffer
	t, ok := templates["navbar-set.html"]
	if !ok {
		slog.Error("Error fetching template: navbar-set")
		return ""
	}
	err := t.Execute(&buff, dataArray{
		dataMap{
			"URL":  "/auth/login",
			"Name": "Login",
		},
		dataMap{
			"URL":  "/auth/register",
			"Name": "Register",
		},
	})
	if err != nil {
		slog.Error("Error parsing template: ", err)
		return ""
	}
	return buff.String()
}
