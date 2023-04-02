package memory

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func (m *Module) loadFromFile() error {
	_, err := os.Stat("memories.json")
	if os.IsNotExist(err) {
		return nil
	}

	b, err := ioutil.ReadFile("memories.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &m.memories)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) writeToFile() error {
	b, err := json.Marshal(m.memories)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("memories.json", b, 0660)
	if err != nil {
		return err
	}
	return nil
}

func (m *Module) appendMemory(mem memory) error {
	m.memories = append(m.memories, mem)
	return m.writeToFile()
}
