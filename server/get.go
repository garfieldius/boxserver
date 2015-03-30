package server

import (
	c "github.com/trenker/boxserver/conf"
	d "github.com/trenker/boxserver/data"
	"github.com/trenker/boxserver/log"
	"github.com/trenker/boxserver/util"
	"net/http"
	"strings"
)

func handleGet(parts []string) (interface{}, int) {

	log.Debug("Find data for [%s] %s", len(parts), parts)

	var s int
	var r interface{}

	switch len(parts) {
	case 1:
		p := parts[0]
		if p == "" {
			log.Debug("Sending all projects")
			r = d.ProjectNames()
			s = http.StatusOK
			break
		}
		if util.ValidKey(p) {
			log.Debug("Sending boxes of project %s", p)
			boxes, err := d.BoxesOfProject(p)
			if err != nil {
				log.Warn("Not found %s", p)
				r = util.Err(err)
				s = http.StatusNotFound
			} else {
				r = boxes
				s = http.StatusOK
			}
		} else {
			s = http.StatusNotAcceptable
			r = util.Str("Not a vaild project ID")
		}
		break

	case 2:

		p := parts[0]
		b := parts[1]
		log.Debug("Requested versions of %s/%s", p, b)

		if !util.ValidKey(p) || !util.ValidKey(b) {
			log.Warn("%s/%s is not valid", p, b)
			r = util.Str("Not a valid project/box key")
			s = http.StatusNotAcceptable
		} else {

			l, err := d.VersionsOfBox(p, b)

			if err != nil || l == nil {
				log.Error("%s", err)
				r = util.Str("Box not found")
				s = http.StatusNotFound
			} else {
				r = buildVagrantData(p, b, l)
				s = http.StatusOK
			}
		}
		break

	default:
		s = http.StatusNotFound
		r = util.Str("Resource not found")
	}

	return r, s
}

type VagrantProvider struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type VagrantVersion struct {
	Version   string            `json:"version"`
	Status    string            `json:"status"`
	Html      string            `json:"description_html"`
	Md        string            `json:"description_markdown"`
	Providers []VagrantProvider `json:"providers"`
}

type VagrantBox struct {
	Description string           `json:"description"`
	Short       string           `json:"short_description"`
	Name        string           `json:"name"`
	Versions    []VagrantVersion `json:"versions"`
}

func buildVagrantData(project, box string, versions []d.Version) *VagrantBox {
	vagrant := &VagrantBox{
		Name:        project + "/" + box,
		Short:       "",
		Description: "",
		Versions:    make([]VagrantVersion, 0),
	}

	for _, v := range versions {

		nV := VagrantVersion{
			Version:   v.Version,
			Status:    "active",
			Providers: make([]VagrantProvider, 0),
		}

		for _, p := range v.Providers {
			providerName := string(p.Type)
			nP := VagrantProvider{
				Name: providerName,
				Url:  makeUrl(p.File),
			}
			nV.Providers = append(nV.Providers, nP)
		}

		vagrant.Versions = append(vagrant.Versions, nV)
	}

	return vagrant
}

func makeUrl(file string) string {
	url := c.Get().BaseUrl

	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	url += file

	return url
}
