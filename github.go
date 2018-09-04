package main

type GithubConfig struct {
	Password string `json:"password"`
}

type GithubWebHook struct {
	Repository Repository2 `json:"repository"`
	C GithubConfig `json:"config"`
}


type Repository2 struct {
	Url  string `json:"git_url"`
	Name string `json:"name"`
}

func (g GithubWebHook) getUrl() string{
	return g.Repository.Url
}

func (g GithubWebHook) getName() string{
	return g.Repository.Name
}

func (g GithubWebHook) getPwd() string{
	return g.C.Password
}