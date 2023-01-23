package tui

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func setPasswd() {
	answers := struct {
		Password string
	}{}
	var qs = []*survey.Question{
		{
			Name:     "Password",
			Prompt:   &survey.Password{Message: "Please enter your admin API password："},
			Validate: survey.Required,
		},
	}
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	mac := hmac.New(sha512.New, []byte("T1mg1bl4bl0g"))
	mac.Write([]byte(answers.Password))
	expectedMAC := mac.Sum(nil)
	answers.Password = hex.EncodeToString(expectedMAC)
	marshal, _ := json.Marshal(answers)
	err = os.WriteFile("./config/control.json", marshal, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
