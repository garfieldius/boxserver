package data

import (
	"errors"
)

type VagrantProvider string

const (
	Virtualbox VagrantProvider = "virtualbox"
	Vmware     VagrantProvider = "vmware_desktop"
	Docker     VagrantProvider = "docker"
	Hyperv     VagrantProvider = "hyperv"
)

type Data struct {
	Projects []Project
}

type Version struct {
	Version   string
	Providers []VagrantProvider
}

type Box struct {
	Name     string
	Versions []Version
}

type Project struct {
	Name  string
	Boxes []Box
}

func ProviderByName(name string) (VagrantProvider, error) {
	switch name {
		case "virtualbox":
			return Virtualbox, nil
		case "vmware":
		case "vmware_desktop":
		case "vmware_workstation":
		case "vmware_fusion":
			return Vmware, nil
		case "docker":
			return Docker, nil
		case "hyperv":
			return Hyperv, nil
	}

	return Virtualbox, errors.New("No such provider")
}

func (d *Data) addProject(project Project) *Project {
	d.Projects = append(d.Projects, project)
	return &d.Projects[d.Len()-1]
}

func (d *Data) getProject(name string) *Project {
	for i, p := range d.Projects {
		if p.Name == name {
			return &d.Projects[i]
		}
	}
	return nil
}

func (d *Data) Len() int {
	return len(d.Projects)
}

func (p *Project) addBox(box Box) *Box {
	p.Boxes = append(p.Boxes, box)
	return &p.Boxes[len(p.Boxes)-1]
}

func (p *Project) getBox(box string) *Box {
	for i, b := range p.Boxes {
		if b.Name == box {
			return &p.Boxes[i]
		}
	}
	return nil
}

func (b *Box) Len() int {
	return len(b.Versions)
}

func (b *Box) addVersion(version Version) *Version {
	b.Versions = append(b.Versions, version)
	return &b.Versions[len(b.Versions)-1]
}

func (b *Box) getVersion(version string) *Version {
	for i, v := range b.Versions {
		if v.Version == version {
			return &b.Versions[i]
		}
	}
	return nil
}

func (v *Version) Len() int {
	return len(v.Providers)
}

func (v *Version) addProvider(provider VagrantProvider) *VagrantProvider {
	v.Providers = append(v.Providers, provider)
	return &v.Providers[len(v.Providers)-1]
}
