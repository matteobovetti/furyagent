package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	getter "github.com/hashicorp/go-getter"
)

const (
	configFile     = "Furyfile"
	protocol       = "git::ssh://"
	modulesRepo    = "git@git.sighup.io/sighup/fury-modules.git"
	rolesRepo      = "git@git.sighup.io/sighup/fury-roles.git"
	katalogRepo    = "git@git.sighup.io/sighup/fury-katalog.git"
	modulesRepoDev = "git@git.incubator.sh/sighup/fury-modules.git"
	rolesRepoDev   = "git@git.incubator.sh/sighup/fury-roles.git"
	katalogRepoDev = "git@git.incubator.sh/sighup/fury-katalog.git"

	defaultVendorFolderName = "vendor"
)

// Furyconf is reponsible for the structure of the Furyfile
type Furyconf struct {
	VendorFolderName string    `yaml:"vendorFolderName"`
	SSHKeyPath       string    `yaml:"sshKeyPath"`
	Roles            []Package `yaml:"roles"`
	Modules          []Package `yaml:"modules"`
	Bases            []Package `yaml:"bases"`
}

// Package is the type to contain the definition of a single package
type Package struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`

	url string
	dit string
}

// Download is the main function to put all the files in vendor folder
func (f *Furyconf) Download(dev bool) error {
	var rolesRepoURL string
	var modulesRepoURL string
	var katalogRepoURL string

	if dev {
		rolesRepoURL = rolesRepoDev
		modulesRepoURL = modulesRepoDev
		katalogRepoURL = katalogRepoDev
	} else {
		rolesRepoURL = rolesRepo
		modulesRepoURL = modulesRepo
		katalogRepoURL = katalogRepo
	}

	for _, v := range f.Roles {
		url := fmt.Sprintf("%s%s//%s?ref=%s", protocol, rolesRepoURL, v.Name, v.Version)
		dir := fmt.Sprintf("%s/%s/%s", f.VendorFolderName, "roles", v.Name)
		err := get(url, dir)
		if err != nil {
			return err
		}
	}
	for _, v := range f.Modules {
		url := fmt.Sprintf("%s%s//%s?ref=%s", protocol, modulesRepoURL, v.Name, v.Version)
		dir := fmt.Sprintf("%s/%s/%s", f.VendorFolderName, "modules", v.Name)
		err := get(url, dir)
		if err != nil {
			return err
		}
	}
	for _, v := range f.Bases {
		url := fmt.Sprintf("%s%s//%s?ref=%s", protocol, katalogRepoURL, v.Name, v.Version)
		dir := fmt.Sprintf("%s/%s/%s", f.VendorFolderName, "bases", v.Name)
		err := get(url, dir)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validate is used for validation of configuration and initization of default paramethers
func (f *Furyconf) Validate() error {
	if f.VendorFolderName == "" {
		f.VendorFolderName = defaultVendorFolderName
	}
	return nil
}

func get(src, dest string) error {
	fmt.Println("DOWNLOADING...\nSRC: ", src, "\nDST: ", dest)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	client := &getter.Client{
		Src:  src,
		Dst:  dest,
		Pwd:  pwd,
		Mode: getter.ClientModeDir,
	}
	return client.Get()
}

func addSSHKey(url, sshKeyPath string) string {
	if sshKeyPath != "" {
		sshKeyData, err := ioutil.ReadFile(sshKeyPath)
		if err != nil {
			log.Println("couldn't find or read provided SSHKEY, ignoring error and continuing")
			log.Println("ERR:", err)
			return url
		}
		url = fmt.Sprintf("%s&sshkey=%s", url, base64.StdEncoding.EncodeToString(sshKeyData))
	}
	return url
}