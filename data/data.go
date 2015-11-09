package data

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/trenker/boxserver/log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var data *Data
var prefix string
var findComponents *regexp.Regexp

func init() {
	divider := `\` + string(os.PathSeparator)
	validKey := `[a-z0-9][a-z0-9_\-]*[a-z0-9]`
	providers := []string{
		string(Virtualbox),
		string(Vmware), "vmware", "vmware_fusion", "vmware_workstation",
		string(Docker),
		string(Hyperv),
	}
	allowedBoxes := `(` + strings.Join(providers, "|") + ")"
	validVersion := `[0-9]+\.[0-9]+\.[0-9]+`

	completeRegex := `^` +
		validKey +
		divider +
		validKey +
		divider +
		validVersion +
		divider +
		allowedBoxes +
		`\.box$`

	log.Debug("Boxes regex check is %s", completeRegex)

	findComponents = regexp.MustCompile(completeRegex)
}

func readFile(path string, info os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	AddFromPath(path, info)

	return nil
}

func AddFromPath(path string, info os.FileInfo) {

	if !info.IsDir() && len(path) > len(prefix) {

		path = strings.TrimPrefix(path, prefix)

		if findComponents.MatchString(path) {

			parts := strings.Split(strings.TrimSuffix(path, ".box"), "/")

			log.Debug("Found box %s", parts)

			providerName, _ := ProviderByName(parts[3])

			var p *Project
			var b *Box
			var v *Version

			log.Debug("Append %s", parts)

			p = data.getProject(parts[0])

			if p == nil {
				p = &Project{Name: parts[0], Boxes: make([]*Box, 0)}
				data.addProject(p)
			}

			b = p.getBox(parts[1])

			if b == nil {
				b = &Box{Name: parts[1], Versions: make([]*Version, 0)}
				p.addBox(b)
			}

			v = b.getVersion(parts[2])

			if v == nil {
				v = &Version{Version: parts[2], Providers: make([]*Provider, 0)}
				b.addVersion(v)
			}

			file := strings.TrimPrefix(path, prefix)

			log.Debug("Add provider %s for %s", providerName, file)
			v.addProvider(providerName, file, humanize.Bytes(uint64(info.Size())))
		}
	}
}

func LoadData() {
	data = new(Data)
	data.Projects = make([]*Project, 0)

	log.Debug("Search for box files in %s", prefix)

	filepath.Walk(prefix, readFile)
}

func Initialize(basePath string) *Data {
	prefix = basePath + string(os.PathSeparator)

	LoadData()

	return data
}
