package gopaperapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Project represents and delivers the methods associated with a PaperMC.io project.
type Project struct {
	client             *http.Client
	project            PaperProject
	projectInformation *ProjectInformation
	projectVersions    map[string]ProjectVersion
}

// NewProject initializes and returns a Project given a PaperProject and http Client.
func NewProject(paperProject PaperProject, httpClient *http.Client) Project {
	return Project{
		client:          httpClient,
		project:         paperProject,
		projectVersions: map[string]ProjectVersion{},
	}
}

// GetProjectInformation gets the information about a project including the versions available.
func (p Project) GetProjectInformation() (info ProjectInformation, err error) {
	if p.projectInformation == nil {
		url := fmt.Sprintf("%s/%s", apiRoot, p.project)
		res, err := p.client.Get(url)
		if err != nil {
			return info, err
		}
		if res.Body != nil {
			defer res.Body.Close()
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return info, err
		}
		projectInfo := &ProjectInformation{}
		err = json.Unmarshal(body, projectInfo)
		if err != nil {
			return info, err
		}
		p.projectInformation = projectInfo
	}
	return *p.projectInformation, nil
}

// GetProjectVersion gets the builds of a version for this project if it exists.
func (p Project) GetProjectVersion(versionString string) (version ProjectVersion, err error) {
	if val, ok := p.projectVersions[versionString]; ok {
		return val, nil
	}

	projectInfo, err := p.GetProjectInformation()
	if err != nil {
		return version, err
	}

	if !IsStringInSlice(versionString, projectInfo.Versions) {
		return version, errors.New("version does not exist for that project")
	}
	url := fmt.Sprintf("%s/%s/%s", apiRoot, p.project, versionString)
	res, err := p.client.Get(url)
	if err != nil {
		return version, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return version, err
	}
	projectVersion := &ProjectVersion{}
	err = json.Unmarshal(body, projectVersion)
	if err != nil {
		return version, err
	}
	p.projectVersions[versionString] = *projectVersion
	return p.projectVersions[versionString], nil
}

// GetLatestVersion gets and returns the latest version of the project.
func (p Project) GetLatestVersion() (version string, err error) {
	projectInfo, err := p.GetProjectInformation()
	if err != nil {
		return version, err
	}

	return projectInfo.Versions[0], nil
}

// GetLatestBuild gets and returns the latest version and build of the project.
func (p Project) GetLatestBuild() (version, build string, err error) {
	version, err = p.GetLatestVersion()
	if err != nil {
		return "", "", err
	}
	builds, err := p.GetProjectVersion(version)
	if err != nil {
		return "", "", err
	}
	return version, builds.Builds.Latest, nil
}

// GetVersionLatestBuild gets and returns the latest version and build of the project.
func (p Project) GetVersionLatestBuild(versionString string) (build string, err error) {
	projectInfo, err := p.GetProjectInformation()
	if err != nil {
		return build, err
	}

	if !IsStringInSlice(versionString, projectInfo.Versions) {
		return "", errors.New("version doesn't exist")
	}
	builds, err := p.GetProjectVersion(versionString)
	if err != nil {
		return "", err
	}
	return builds.Builds.Latest, nil
}
