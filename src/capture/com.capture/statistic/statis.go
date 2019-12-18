package statistic

import (
	"container/list"
	"sync"
)

var ReceivePackageNum uint32 = 0

var NetCardList = new(list.List)
var PIP sync.Map
var FromToIP sync.Map
