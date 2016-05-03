package input

import "github.com/dcw303/crenshaw-go/util"

// Look is a Lookahead Character
var Look rune

// GetChar Reads New Character from Input Stream
func GetChar() {
	Look = util.Read()
}

func init() {
	GetChar()
}
