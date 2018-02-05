package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/nlopes/slack"
)

func main() ***REMOVED***
	var (
		verificationToken string
	)

	flag.StringVar(&verificationToken, "token", "YOUR_VERIFICATION_TOKEN_HERE", "Your Slash Verification Token")
	flag.Parse()

	http.HandleFunc("/slash", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		s, err := slack.SlashCommandParse(r)
		if err != nil ***REMOVED***
			w.WriteHeader(http.StatusInternalServerError)
			return
		***REMOVED***

		if !s.ValidateToken(verificationToken) ***REMOVED***
			w.WriteHeader(http.StatusUnauthorized)
			return
		***REMOVED***

		switch s.Command ***REMOVED***
		case "/echo":
			params := &slack.Msg***REMOVED***Text: s.Text***REMOVED***
			b, err := json.Marshal(params)
			if err != nil ***REMOVED***
				w.WriteHeader(http.StatusInternalServerError)
				return
			***REMOVED***
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		***REMOVED***
	***REMOVED***)
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
***REMOVED***
