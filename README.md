# Store
>Store is a dead simple configuration manager for Go applications.

[![GoDoc](https://godoc.org/github.com/tucnak/store?status.svg)](https://godoc.org/github.com/tucnak/store)

I didn't like existing configuration management solutions like [globalconf](https://github.com/rakyll/globalconf), [tachyon](https://github.com/vektra/tachyon) or [viper](https://github.com/spf13/viper). First two just don't feel right and viper, imo, a little overcomplicated—definitely offering too much for small things. Store supports either JSON, TOML or YAML out-of-the-box and lets you register practically any other configuration format. 

Look, when I say it's dead simple, I actually mean it:
```go
package main

import (
	"log"
	"time"

	"github.com/tucnak/store"
)

type Cat struct {
	Name   string `toml:"naym"`
	Clever bool   `toml:"ayy"`
}

type Hotel struct {
	Name string
	Cats []Cat `toml:"guests"`

OpeningHours
}

func main() {
  hotel := Hotel {
    Name: "Grand Budapest Hotel",
		Cats: []Cat{
			{"Rudolph", true},
			{"Patrick", false},
			{"Jeremy", true},
		},
    OpeningHours: store.Duration(time.Hour*8)
  }

// Load from some path, if file is not found it will be created with the values currently in the struct
	if err := store.Load("config/hotel.toml", &hotel); err != nil {
		log.Println("failed to load the cat hotel:", err)
		return
	}

	// ...

	// Save to some path
	if err := store.Save("config/hotel.toml", &hotel); err != nil {
		log.Println("failed to save the cat hotel:", err)
		return
	}
}
```

Store supports any other formats via the handy registration system: register the format once and you'd be able to Load and Save files in it afterwards:
```go
store.Register("ini", ini.Marshal, ini.Unmarshal)

err := store.Load("configuration.ini", &object)
// ...
```

## Duration

Store provides a special Duration type. Use it instead of time.Duration in your configs to make them marshall to strings in your config files
```go

type Cat struct {
	Name   string `toml:"naym"`
	Clever bool   `toml:"ayy"`
}

type Hotel struct {
	Name string
	Cats []Cat `toml:"guests"`

  OpeningHours config.Duration
}


func main() {
  hotel := Hotel {
    OpeningHours: store.Duration(time.Hour*8)
  }


  // ....
}

```
## Concurrency

In a multithreading scenario it is recommended that you use the **[atomic.Value](https://golang.org/pkg/sync/atomic/#Value)** interface to handle the config itself in a sane way



