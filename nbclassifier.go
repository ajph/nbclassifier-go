package nbclassifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Model struct {
	Classes []*Class `json:"c"`
	Total   int      `json:"t"`
}

type ClassItem struct {
	Id    string `json:"n"`
	Count int    `json:"c"`
}

type Class struct {
	Id    string       `json:"n"`
	Items []*ClassItem `json:"i"`
	Total int          `json:"t"`
}

type ScoreResult struct {
	Class *Class
	Score float64
}

func (m *Model) SaveToFile(path string) error {
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fo.Close()
	json.NewEncoder(fo).Encode(m)
	return nil
}

func (m *Model) FindClass(id string) (*Class, bool) {
	for i, _ := range m.Classes {
		if m.Classes[i].Id == id {
			return m.Classes[i], true
		}
	}
	return nil, false
}

func (m *Model) NewClass(id string) (*Class, error) {
	if _, ok := m.FindClass(id); ok {
		return nil, errors.New("class already exists")
	}
	c := &Class{
		Id:    id,
		Total: 0,
	}
	m.Classes = append(m.Classes, c)
	return c, nil
}

func (m *Model) Score(t []string) ([]*ScoreResult, error) {
	if len(m.Classes) < 2 {
		return nil, errors.New("need more than 1 class")
	}

	res := make([]*ScoreResult, len(m.Classes))
	var scoreSum float64 = 0
	for i, _ := range m.Classes {
		var prior float64 = float64(m.Classes[i].Total) / float64(m.Total)
		res[i] = &ScoreResult{m.Classes[i], prior}
		for _, s := range t {
			if item, ok := m.Classes[i].FindItem(s); ok {
				res[i].Score *= (float64(item.Count) / float64(m.Classes[i].Total))
			} else {
				res[i].Score *= 0.00000000001
			}
		}
		scoreSum += res[i].Score
	}

	for i, _ := range res {
		res[i].Score /= scoreSum
	}

	return res, nil
}

func (m *Model) Classify(t ...string) (*Class, bool, error) {
	res, err := m.Score(t)
	if err != nil {
		return nil, false, err
	}
	unsure := false
	winner := res[0]
	for i := 1; i < len(res); i++ {
		if res[i].Score > winner.Score {
			winner = res[i]
			unsure = false
		} else if res[i].Score == winner.Score {
			unsure = true
		}
	}
	return winner.Class, unsure, nil
}

func (c *Class) FindItem(id string) (*ClassItem, bool) {
	for i, _ := range c.Items {
		if c.Items[i].Id == id {
			return c.Items[i], true
		}
	}
	return nil, false
}

func (m *Model) Learn(class string, id ...string) error {
	c, ok := m.FindClass(class)
	if !ok {
		return fmt.Errorf("cannot find class %s", class)
	}
	for _, v := range id {
		if item, ok := c.FindItem(v); ok {
			item.Count++
		} else {
			c.Items = append(c.Items, &ClassItem{v, 1})
		}
		c.Total++
		m.Total++
	}
	return nil
}

func New() *Model {
	return &Model{}
}

func LoadFromFile(path string) (*Model, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	m := &Model{}
	json.NewDecoder(fi).Decode(m)
	return m, nil
}
