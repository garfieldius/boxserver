package server

import (
	"github.com/trenker/boxserver/conf"
	"github.com/trenker/boxserver/data"
	"github.com/trenker/boxserver/log"
	"github.com/trenker/boxserver/util"
	"io"
	"net/http"
	"os"
	"strings"
)

type pendingFile struct {
	targetPath string
	srcPath    string
	name       string
}

func handlePut(parts []string, req *http.Request) (util.Message, int) {

	result := util.Str("Invalid Request")
	status := http.StatusBadRequest

	log.Debug("Handle PUT with path %s", parts)

	if len(parts) < 3 {
		return util.Str("Cannot post a box without ID"), status
	}

	if !util.ValidKey(parts[0]) || !util.ValidKey(parts[1]) {
		return util.Str("Not a valid box name"), status
	}

	if !util.ValidVersion(parts[2]) {
		return util.Str("Not a valid box version"), status
	}

	if len(parts) == 3 {

		log.Debug("Parts has three, expecting a from src copy")

		src := strings.Split(req.FormValue("source"), "/")

		log.Debug("copy from given source %s => %s", req.FormValue("source"), src)

		if len(src) != 3 {
			return util.Str("Not a valid source ID"), status
		}

		if !util.ValidKey(src[0]) || !util.ValidKey(src[1]) || !util.ValidVersion(src[2]) {
			return util.Str("Not a valid source ID"), status
		}

		versionsOfBox, err := data.VersionsOfBox(src[0], src[1])
		log.Debug("Got versions for %s/%s: [%s] %s", src[0], src[1], len(versionsOfBox), versionsOfBox)

		if err != nil {
			log.Error("Cannot copy %s/%s, not found %s", src[0], src[1], err)
			return util.Str("Source box not found"), http.StatusNotFound
		}

		var srcVersion data.Version
		foundVersion := false

		for _, b := range versionsOfBox {
			log.Debug("Check version %s to match %s", b.Version, src[2])
			if b.Version == src[2] {
				log.Debug("Found %s", b)
				srcVersion = b
				foundVersion = true
				break
			}
		}

		if foundVersion == false {
			return util.Str("Source version not found"), http.StatusNotFound
		}

		files := make([]pendingFile, len(srcVersion.Providers))

		for i, p := range srcVersion.Providers {
			providerfile := (string)(p) + ".box"

			f := pendingFile{
				targetPath: util.Join(conf.Get().Data, parts[0], parts[1], parts[2]),
				srcPath:    util.Join(conf.Get().Data, src[0], src[1], src[2]),
				name:       providerfile,
			}

			log.Debug("Need to copy %s to %s", util.Join(f.srcPath, f.name), util.Join(f.targetPath, f.name))

			if util.FileExists(util.Join(f.targetPath, f.name)) {
				return util.Str("Box file already exists"), status
			}

			if !util.FileExists(util.Join(f.srcPath, f.name)) {
				return util.Str("Source Box file not found"), status
			}

			files[i] = f
		}

		log.Debug("Copy schedule is %s", files)

		for _, file := range files {
			err := os.MkdirAll(file.targetPath, (os.FileMode)(0755))

			if err != nil {
				log.Error("Creating directory %s: %s", file.targetPath, err)
				return util.Str("Cannot create target directory"), http.StatusInternalServerError
			}

			srcFile := util.Join(file.srcPath, file.name)
			dstFile := util.Join(file.targetPath, file.name)

			src, err := os.Open(srcFile)
			defer src.Close()

			if err != nil {
				log.Error("Open Source file %s: %s", srcFile, err)
				return util.Str("Cannot open source file"), http.StatusInternalServerError
			}

			dst, err := os.Create(dstFile)
			defer dst.Close()

			if err != nil {
				log.Error("Create target file %s: %s", dstFile, err)
				return util.Str("Cannot open target file"), http.StatusInternalServerError
			}

			_, err = io.Copy(dst, src)

			if err != nil {
				log.Error("Copy from %s to %s: %s", srcFile, dstFile, err)
				return util.Str("Cannot copy to target file"), http.StatusInternalServerError
			}

			info, err := os.Stat(dstFile)

			if err != nil {
				log.Error("Cannot stat %s: %s", dstFile, err)
				return util.Str("Target box file error"), http.StatusInternalServerError
			}

			data.AddFromPath(dstFile, info)
		}

		return util.Str("Box copied"), http.StatusOK
	}

	if len(parts) == 4 {

		log.Debug("parts has a length of four, expecting a file upload for %s", parts)

		if !util.ValidKey(parts[0]) || !util.ValidKey(parts[1]) || !util.ValidVersion(parts[2]) {
			return util.Str("Not a valid source ID"), status
		}

		if !util.ValidProvider(parts[3]) {
			return util.Str("No such provider"), status
		}

		targetDir := util.Join(conf.Get().Data, parts[0], parts[1], parts[2])
		targetFile := parts[3] + ".box"

		_, srcH, err := req.FormFile("box")

		if err != nil {
			log.Error("Cannot read box upload file: %s", err)
			return util.Str("Cannot read box upload file"), status
		}

		err = os.MkdirAll(targetDir, os.FileMode(0755))

		if err != nil {
			log.Error("Cannot create target dir %s: %s", targetDir, err)
			return util.Str("Cannot create target directory"), http.StatusInternalServerError
		}

		dstFile := util.Join(targetDir, targetFile)
		dst, err := os.Create(dstFile)

		if err != nil {
			log.Error("Cannot create target file %s: %s", util.Join(targetDir, targetFile), err)
			return util.Str("Cannot create target file"), http.StatusInternalServerError
		}
		defer dst.Close()

		src, err := srcH.Open()

		if err != nil {
			log.Error("Cannot open uploaded file %s: %s", srcH, err)
			return util.Str("Cannot open uploaded file"), http.StatusInternalServerError
		}
		defer src.Close()

		_, err = io.Copy(dst, src)

		if err != nil {
			log.Error("Copy from %s to %s: %s", srcH, dstFile, err)
			return util.Str("Cannot copy to target file"), http.StatusInternalServerError
		}

		info, err := os.Stat(dstFile)

		if err != nil {
			log.Error("Cannot stat %s: %s", dstFile, err)
			return util.Str("Target box file error"), http.StatusInternalServerError
		}

		data.AddFromPath(dstFile, info)

		return util.Str("Box added"), http.StatusOK
	}

	log.Error("Invalid count of params, falling back to not acceptable state")

	return result, status
}
