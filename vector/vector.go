package vector

import (
	"hmruntime/vector/in_mem"
	"hmruntime/vector/index"
)

var GlobalInMemIndexFactory index.NamedFactory[float64]

func InitializeInMemFactory() {
	GlobalInMemIndexFactory = in_mem.CreateFactory()
}
