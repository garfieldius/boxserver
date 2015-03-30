package data

import "errors"

func DeleteProject(projectName string) error {
	newProjects := make([]Project, 0)
	found := false

	for _, p := range data.Projects {
		if p.Name != projectName {
			newProjects = append(newProjects, p)
		} else {
			found = true
		}
	}

	if found == true {
		data.Projects = newProjects
		return nil
	} else {
		return errors.New("No such project")
	}
}

func DeleteBox(projectName, boxName string) error {

	for i, p := range data.Projects {
		if p.Name == projectName {

			newBoxes := make([]Box, 0)
			found := false

			for _, b := range p.Boxes {
				if b.Name != boxName {
					newBoxes = append(newBoxes, b)
				} else {
					found = true
				}
			}

			if found == true {

				if len(newBoxes) == 0 {
					return DeleteProject(projectName)
				} else {
					data.Projects[i].Boxes = newBoxes
					return nil
				}
			}
		}
	}

	return errors.New("No such project/box found")
}

func DeleteVersion(projectName, boxName, version string) error {

	for i, p := range data.Projects {
		if p.Name == projectName {

			for j, b := range p.Boxes {
				if b.Name == boxName {
					found := false
					newVersions := make([]Version, 0)

					for _, v := range b.Versions {
						if v.Version == version {
							found = true
						} else {
							newVersions = append(newVersions, v)
						}
					}

					if found == true {
						if len(newVersions) > 0 {
							data.Projects[i].Boxes[j].Versions = newVersions
							return nil
						} else {
							return DeleteBox(projectName, boxName)
						}
					}
				}
			}
		}
	}

	return errors.New("No such project/box/version found")
}

func DeleteProvider(projectName, boxName, version string, provider VagrantProvider) error {

	for i, p := range data.Projects {
		if p.Name == projectName {

			for j, b := range p.Boxes {
				if b.Name == boxName {

					for k, v := range b.Versions {
						if v.Version == version {

							found := false
							newProviders := make([]Provider, 0)

							for _, pr := range v.Providers {
								if pr.Type == provider {
									found = true
								} else {
									newProviders = append(newProviders, pr)
								}
							}

							if found == true {

								if len(newProviders) > 0 {
									data.Projects[i].Boxes[j].Versions[k].Providers = newProviders
									return nil
								} else {
									return DeleteVersion(projectName, boxName, version)
								}
							}
						}
					}
				}
			}
		}
	}

	return errors.New("No such project/box/version/provider found")
}
