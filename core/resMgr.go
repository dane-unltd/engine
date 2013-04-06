package core

import "strings"
import "encoding/json"
import "os"
import "log"

type ResMgr struct {
	path string
}

func NewResMgr(base string) *ResMgr {
	var resourcePath string
	{
		GOPATH := os.Getenv("GOPATH")
		GOPATH = GOPATH + ":."

		for _, gopath := range strings.Split(GOPATH, ":") {
			a := gopath + "/res/" + base
			_, err := os.Stat(a)
			if err == nil {
				resourcePath = a
				break
			}
		}
		if resourcePath == "" {
			log.Fatal("Failed to find resource directory")
		}
	}
	return &ResMgr{resourcePath}
}

func (r *ResMgr) ReadConfig(name string, out interface{}) error {
	file, err := os.Open(r.path + "/cfg/" + name) // For read access.
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(out)
	if err != nil {
		return err
	}
	return nil
}
