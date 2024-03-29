package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"gochain/constcoe"
	"gochain/utils"

	// "io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type RefList map[string]string

func (r *RefList) Save(){
	filename := constcoe.WalletsRefList + "ref_list.data"
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(r)
	utils.Handle(err)
	err = os.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

func (r *RefList) Update() {
	err := filepath.Walk(constcoe.Wallets, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileName := f.Name()
		if strings.Compare(fileName[len(fileName)-4:], ".wlt") == 0 {
			_, ok := (*r)[fileName[:len(fileName)-4]]
			if !ok {
				(*r)[fileName[:len(fileName)-4]] = ""
			}
		}
		return nil
	})
	utils.Handle(err)
}

func LoadRefList() *RefList {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var reflist RefList
	if utils.FileExists(filename) {
		fileContent, err := os.ReadFile(filename)
		utils.Handle(err)
		decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
		err = decoder.Decode(&reflist)
		utils.Handle(err)
	} else {
		reflist = make(RefList)
		reflist.Update()
	}
	return &reflist
}


func (r *RefList) BindRef(address, refname string) {
	(*r)[address] = refname
}

func (r *RefList) FindRef(refname string) (string, error) {
	tmp := ""
	for key, val:= range *r{
		if val == refname{
			tmp = key
			break
		}
	}

	if tmp == ""{
		err := errors.New("the refname is not found")
		return "", err
	}else{
		return tmp, nil
	}

}