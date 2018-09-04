package main


type GiteeWebhook struct {
	HookName   string     `json:"hook_name"`
	Password   string     `json:"password"`
	Ref        string     `json:"ref"` //是否是 "refs/heads/master" 目前只需要支持master即可
	Repository Repository `json:"repository"`
}

type Repository struct {
	GitHttpUrl string `json:"git_http_url"`
	Path       string `json:"path"`
}

func (g GiteeWebhook) getUrl() string{
	return g.Repository.GitHttpUrl
}

func (g GiteeWebhook) getName() string{
	return g.Repository.Path
}

func (g GiteeWebhook) getPwd() string{
	return g.Password
}