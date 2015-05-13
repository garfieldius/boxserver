package data

import (
	"errors"
)

func ProjectNames() []string {
	r := make([]string, len(data.Projects))
	for i, p := range data.Projects {
		r[i] = p.Name
	}
	return r
}

func BoxesOfProject(project string) ([]string, error) {

	p := data.getProject(project)

	if p == nil {
		return nil, errors.New("No such project")
	}

	r := make([]string, len(p.Boxes))

	for i, b := range p.Boxes {
		r[i] = b.Name
	}

	return r, nil
}

func VersionsOfBox(project, box string) ([]*Version, error) {
	p := data.getProject(project)

	if p == nil {
		return nil, errors.New("No such project")
	}

	b := p.getBox(box)

	if b == nil {
		return nil, errors.New("No such box")
	}

	return b.Versions, nil
}
