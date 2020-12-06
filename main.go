package main

import (
	"github.com/leighmacdonald/lurkr/cmd"
	_ "github.com/leighmacdonald/lurkr/internal/tracker/broadcasthenet"
	_ "github.com/leighmacdonald/lurkr/internal/tracker/redacted"
	_ "github.com/leighmacdonald/lurkr/internal/tracker/revolutiontt"
	_ "github.com/leighmacdonald/lurkr/internal/tracker/torrentleech"
)

func main() {
	cmd.Execute()
}
