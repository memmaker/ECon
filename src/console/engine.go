package console

import (
	"github.com/memmaker/ECon/input"
)

type Engine interface {
	GetInput() input.GridInput
}
