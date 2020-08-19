package gopaperapi

type PaperProject string

const (
	Paper      PaperProject = "paper"
	Waterfall  PaperProject = "waterfall"
	Travertine PaperProject = "travertine"
)

const apiRoot = "https://papermc.io/api/v1"

type ProjectInformation struct {
	Name     string   `json:"project"`
	Versions []string `json:"versions"`
}

type ProjectVersion struct {
	Project string               `json:"project"`
	Version string               `json:"version"`
	Builds  ProjectVersionBuilds `json:"builds"`
}

type ProjectVersionBuilds struct {
	Latest string   `json:"latest"`
	All    []string `json:"all"`
}
