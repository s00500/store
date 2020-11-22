package store

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Cat struct {
	Name string
	Big  bool
}

type Settings struct {
	Age          int
	Cats         []Cat
	RandomString string
}

func equal(a, b Settings) bool {
	if a.Age != b.Age {
		return false
	}

	if a.RandomString != b.RandomString {
		return false
	}

	if len(a.Cats) != len(b.Cats) {
		return false
	}

	for i, cat := range a.Cats {
		if cat != b.Cats[i] {
			return false
		}
	}

	return true
}

func TestSaveLoad(t *testing.T) {
	settings := Settings{
		Age: 42,
		Cats: []Cat{
			{"Rudolph", true},
			{"Patrick", false},
			{"Jeremy", true},
		},
		RandomString: "gophers are gonna conquer the world",
	}

	settingsFile := "preferences.toml"

	err := Load(settingsFile, &settings)
	if err != nil {
		t.Fatalf("failed to load preferences: %s\n", err)
		return
	}

	err = Save(settingsFile, &settings)
	if err != nil {
		t.Fatalf("failed to save preferences: %s\n", err)
		return
	}

	defer os.Remove(settingsFile)

	var newSettings Settings

	err = Load(settingsFile, &newSettings)
	if err != nil {
		t.Fatalf("failed to load preferences: %s\n", err)
		return
	}

	if !equal(settings, newSettings) {
		t.Fatalf("broken")
	}
}

func TestSaveLoadInvalid(t *testing.T) {
	settings := Settings{
		Age: 42,
		Cats: []Cat{
			{"Rudolph", true},
			{"Patrick", false},
			{"Jeremy", true},
		},
		RandomString: "gophers are gonna conquer the world",
	}

	settingsFile := "preferences.roml"
	assert.Panics(t, func() { Save(settingsFile, &settings) }, "Save should panic on unknown format")
	assert.Panics(t, func() { Load(settingsFile, &settings) }, "Load should panic on unknown format")

	settingsFile = "preferences"
	assert.Panics(t, func() { Save(settingsFile, &settings) }, "Save should panic without .")

}

type RomlMarshaller struct {
	mock.Mock
}

// MarshalFunc is any marshaler.
func (m *RomlMarshaller) MarshalFunc(v interface{}) ([]byte, error) {
	return []byte{}, fmt.Errorf("Failed on purpouse")

}

func (m *RomlMarshaller) UnmarshalFunc(data []byte, v interface{}) error {
	return fmt.Errorf("Failed on purpouse")
}

func TestCustomBrokenMarshaller(t *testing.T) {
	settings := Settings{
		Age: 42,
		Cats: []Cat{
			{"Rudolph", true},
			{"Patrick", false},
			{"Jeremy", true},
		},
		RandomString: "gophers are gonna conquer the world",
	}

	settingsFile := "preferences.roml"

	m := RomlMarshaller{}

	Register("roml", m.MarshalFunc, m.UnmarshalFunc)

	err := Save(settingsFile, &settings)
	assert.Error(t, err)
	err = nil

	err = Load(settingsFile, &settings)
	assert.Error(t, err)

	f, err := os.Create("preferences.roml")
	if err != nil {
		t.FailNow()
	}
	_, err = f.WriteString("invalid")
	if err != nil {
		t.FailNow()
	}
	f.Close()

	err = Load(settingsFile, &settings)
	fmt.Println(err)
	assert.Error(t, err)

	os.Remove(settingsFile)

	err = Load("/settingsFile.toml", &settings)
	assert.Error(t, err)

	err = Save("/settingsFile.toml", &settings)
	assert.Error(t, err)

	err = Save("/fuu/lqwjefnw/settingsFile.toml", &settings)
	assert.Error(t, err)
}
