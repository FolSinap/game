package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type sessionLoginContainer struct {
	Login string `json:"login"`
	UserAgent string `json:"user_agent"`
}

var users = map[string]string{
	"user": "pass",
	"user2": "pass2",
}
var sessions = make(map[string]string)
var htmlLoginForm = `
<html>
	<body>
	<form action="/login" method="post">
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login">
	</form>
	</body>
</html>
`

func initServer() {
	http.HandleFunc("/logs", func(writer http.ResponseWriter, request *http.Request) {
		sessionId, err := request.Cookie("session_id")
		if err == http.ErrNoCookie {
			http.Redirect(writer, request, "/login", http.StatusFound)
			return
		} else if err != nil {
			panicOnError(err)
		}
		jsonData, ok := sessions[sessionId.Value]
		if !ok {
			http.SetCookie(writer, &http.Cookie{Name: sessionId.Name, MaxAge: -1})
			http.Redirect(writer, request, "/login", http.StatusFound)
			return
		}
		data := sessionLoginContainer{}
		panicOnError(json.Unmarshal([]byte(jsonData), &data))
		writer.Write([]byte("Welcome, " + data.Login + " from " + data.UserAgent))
		logs := getLogs(100)
		logsHtml := "\n\n"
		for _, log := range logs {
			logsHtml += "- " + log.From + ": " + log.Command + "\n" + log.Result + "\n" +
				time.Unix(int64(log.Time), 0).Format(time.RFC850) + "\n\n"
		}
		writer.Write([]byte(logsHtml))
	})
	http.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			request.ParseForm()
			login := request.Form["login"][0]
			password := request.Form["password"][0]
			if pasw, ok := users[login]; ok {
				if pasw != password {
					writer.Write([]byte("wrong password"))
					return
				}
			}
			session := sessionLoginContainer{}
			session.Login = login
			session.UserAgent = request.Header.Get("User-Agent")
			jsonData, err := json.Marshal(session)
			panicOnError(err)
			sessionId := RandStringRunes(10)
			sessions[sessionId] = string(jsonData)
			cookie := http.Cookie{Name: "session_id", Value: sessionId, MaxAge: 15 * 60}
			http.SetCookie(writer, &cookie)
			http.Redirect(writer, request, "/logs", http.StatusFound)
			return
		}
		writer.Write([]byte(htmlLoginForm))
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/login", http.StatusFound)
	})
	panicOnError(http.ListenAndServe(":8000", nil))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
