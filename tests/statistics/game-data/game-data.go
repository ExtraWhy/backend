package tests

import (
	"encoding/json"
	"fmt"
	"os"
)

type Iarr struct {
	Data []int32
}

type GamePlay struct {
	Name    string
	Pattern []int32
}

type GameData struct {
	Data     []Iarr
	GamePlay []GamePlay
}

func MakeGameData(fname string) (*GameData, error) {
	//preallocate
	gd := GameData{Data: make([]Iarr, 10),
		GamePlay: make([]GamePlay, 10)}

	file, err := os.ReadFile("kst-data.json")
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	for key, _ := range data {
		if key == "reels" {
			tmp := data[key]
			switch t := tmp.(type) {
			case []any:
				for p, h := range t {
					for _, vv := range h.([]interface{}) {
						jjj := int32(vv.(float64))
						gd.Data[p].Data = append(gd.Data[p].Data, jjj)
					}
				}
			}
		} else if key == "lines" {
			tmp := data[key]
			switch t := tmp.(type) {
			case []any:
				for pp, vv := range t {
					fmt.Println(pp, vv)
					for k1, v1 := range vv.(map[string]interface{}) {
						if k1 == "name" {
							gd.GamePlay[pp].Name = string(v1.(string))
						} else if k1 == "pattern" {
							for _, ll := range v1.([]interface{}) {
								lll := int32(ll.(float64))
								gd.GamePlay[pp].Pattern = append(gd.GamePlay[pp].Pattern, lll)
							}
						}
					}
				}
			}
		} else {
			// none
		}
	}

	return &gd, nil
}
