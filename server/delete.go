package server

import (
	"github.com/trenker/boxserver/conf"
	"github.com/trenker/boxserver/data"
	"github.com/trenker/boxserver/log"
	"github.com/trenker/boxserver/util"
	"net/http"
	"os"
)

func handleDelete(path []string) (util.Message, int) {
	log.Debug("Delete by path %s", path)

	if util.ValidKey(path[0]) {

		filePath := util.Join(conf.Get().Data, path[0])

		if len(path) > 1 {

			if util.ValidKey(path[1]) {

				filePath = util.Join(filePath, path[1])

				if len(path) > 2 {

					if util.ValidVersion(path[2]) {

						filePath = util.Join(filePath, path[2])

						if len(path) > 3 {

							if util.ValidProvider(path[3]) {

								filePath = util.Join(filePath, path[3]+".box")
								provider, _ := data.ProviderByName(path[3])

								err := data.DeleteProvider(path[0], path[1], path[2], provider)

								if err != nil {
									return util.Err(err), http.StatusBadRequest
								}

								return removePath(filePath)

							} else {
								return util.Str("Invalid provider"), http.StatusBadRequest
							}

						} else {

							err := data.DeleteVersion(path[0], path[1], path[2])

							if err != nil {
								return util.Err(err), http.StatusBadRequest
							}

							return removePath(filePath)
						}

					} else {
						return util.Str("Invalid version number"), http.StatusBadRequest
					}

				} else {

					err := data.DeleteBox(path[0], path[1])

					if err != nil {
						return util.Err(err), http.StatusBadRequest
					}

					return removePath(filePath)
				}

			} else {
				return util.Str("Invalid box key"), http.StatusBadRequest
			}

		} else {
			err := data.DeleteProject(path[0])

			if err != nil {
				return util.Err(err), http.StatusBadRequest
			}

			return removePath(filePath)
		}
	}

	return util.Str("Unknown resource path"), http.StatusNotFound
}

func removePath(path string) (util.Message, int) {
	err := os.RemoveAll(path)

	if err != nil {
		return util.Err(err), http.StatusInternalServerError
	}

	return util.Str("Removed"), http.StatusOK
}
