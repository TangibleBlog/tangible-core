package tui

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
)

func getGravatar(answers *system) {
	qs := []*survey.Question{
		{
			Name:     "Email",
			Prompt:   &survey.Input{Message: "Please enter your email address:"},
			Validate: survey.Required,
		},
	}
	result := struct {
		Email string
	}{}
	err := survey.Ask(qs, &result)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	resp, err := http.Get(fmt.Sprintf("https://en.gravatar.com/%s.json", result.Email))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	var respResult GravatarResp
	err = jsoniter.Unmarshal(body, &respResult)
	if err != nil {
		log.Println(err)
		return
	}
	answers.AuthorAvatar = respResult.Entry[0].ThumbnailURL
}

type GravatarResp struct {
	Entry []struct {
		ID                string `json:"id"`
		Hash              string `json:"hash"`
		RequestHash       string `json:"requestHash"`
		ProfileURL        string `json:"profileUrl"`
		PreferredUsername string `json:"preferredUsername"`
		ThumbnailURL      string `json:"thumbnailUrl"`
		Photos            []struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"photos"`
		DisplayName string        `json:"displayName"`
		Urls        []interface{} `json:"urls"`
	} `json:"entry"`
}
