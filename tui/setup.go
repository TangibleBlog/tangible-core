package tui

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"tangible-core/public/common"
	PublicStruct "tangible-core/public/globalstruct"

	"github.com/AlecAivazis/survey/v2"
)

func SetupSystem() {
	_ = os.Mkdir("./config", 0755)
	_ = os.Mkdir("./documents", 0755)
	_ = os.Mkdir("./include", 0755)
	_ = os.Mkdir("./templates", 0755)
	_ = os.WriteFile("./documents/index.json", []byte("[]"), 0755)
	_ = os.WriteFile("./config/menu.json", []byte("[]"), 0755)
	_ = os.WriteFile("./include/comment.html", []byte(""), 0755)
	_ = os.WriteFile("./include/header.html", []byte(""), 0755)
	_ = os.WriteFile("./include/footer.html", []byte(""), 0755)
	answers := system{}
	var qs = []*survey.Question{
		{
			Name:     "ProjectName",
			Prompt:   &survey.Input{Message: "Please enter your blog title："},
			Validate: survey.Required,
		},
		{
			Name:     "ProjectURL",
			Prompt:   &survey.Input{Message: "Please enter your blog URL："},
			Validate: survey.Required,
		},
		{
			Name:     "ProjectDescription",
			Prompt:   &survey.Input{Message: "Please enter your blog description："},
			Validate: survey.Required,
		},
		{
			Name:     "AuthorName",
			Prompt:   &survey.Input{Message: "Please enter author name："},
			Validate: survey.Required,
		},
		{
			Name:     "AuthorIntroduction",
			Prompt:   &survey.Input{Message: "Please enter author introduction："},
			Validate: survey.Required,
		},
	}
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	qs = []*survey.Question{
		{
			Name:     "UseGravatar",
			Prompt:   &survey.Confirm{Message: "Whether to use Gravatar"},
			Validate: survey.Required,
		},
	}
	resultUseGravatar := struct {
		UseGravatar bool
	}{}
	err = survey.Ask(qs, &resultUseGravatar)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if resultUseGravatar.UseGravatar {
		getGravatar(&answers)
	}
	qs = []*survey.Question{
		{
			Name:     "AuthorAvatar",
			Prompt:   &survey.Input{Message: "Please enter the author avatar URL：", Default: answers.AuthorAvatar},
			Validate: survey.Required,
		},
		{
			Name:     "Paging",
			Prompt:   &survey.Input{Message: "Please enter the number of single-page articles：", Default: "10"},
			Validate: survey.Required,
		},
		{
			Name:     "TimeFormat",
			Prompt:   &survey.Input{Message: "Please enter a time conversion format:", Default: "%Y-%m-%d"},
			Validate: survey.Required,
		},
		{
			Name:     "Control",
			Prompt:   &survey.Confirm{Message: "Are you using the management API?"},
			Validate: survey.Required,
		},
		{
			Name:     "Comment",
			Prompt:   &survey.Confirm{Message: "Are you using the Comment?"},
			Validate: survey.Required,
		},
	}

	err = survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	themeOption := loadThemeDir()
	if len(themeOption) != 0 {
		qs = []*survey.Question{
			{
				Name: "Theme",
				Prompt: &survey.Select{
					Message: "Choose a theme:",
					Options: themeOption,
				},
			},
		}
		err = survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	if answers.Control {
		setPasswd()
	}
	var systemConfig PublicStruct.SystemConfigStruct
	project := make(map[string]interface{})
	project["Name"] = answers.ProjectName
	project["URL"] = answers.ProjectURL
	project["Description"] = answers.ProjectDescription
	systemConfig.Project = project
	author := make(map[string]interface{})
	author["Name"] = answers.AuthorName
	author["Introduction"] = answers.AuthorIntroduction
	author["Avatar"] = answers.AuthorAvatar
	systemConfig.Author = author
	systemConfig.Paging = answers.Paging
	systemConfig.Theme = answers.Theme
	systemConfig.TimeFormat = answers.TimeFormat
	systemConfig.Service.Control = answers.Control
	systemConfig.Service.Comment = answers.Comment
	marshal, _ := json.Marshal(systemConfig)
	err = os.WriteFile("./config/system.json", marshal, 0755)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func loadThemeDir() []string {
	entries, err := os.ReadDir("./templates/")
	if err != nil || len(entries) == 0 {
		log.Println(err)
	}
	var result []string
	for _, entry := range entries {
		info, err := entry.Info()
		common.HandleError(err)
		if info.IsDir() {
			result = append(result, info.Name())
		}

	}

	return result

}

type system struct {
	ProjectName        string
	ProjectURL         string
	ProjectDescription string
	AuthorName         string
	AuthorIntroduction string
	AuthorAvatar       string
	Paging             int
	Theme              string
	Control            bool
	Comment            bool
	TimeFormat         string
}
