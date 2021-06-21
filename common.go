package systemstat

func calcPercent(total, part int64) int {
	if total <= 0 {
		return 0
	}
	return int((part * 1000) / total)
}

type Hd struct {
	Unit       string
	Total      int64
	Used       int64
	Percent    int
	DetailList []*HdDetail
}

type HdDetail struct {
	Name    string
	Total   int64
	Used    int64
	Percent int
}

type Network struct {
	Unit       string
	Total      int64
	PerSecond  int64
	DetailList []*NetworkDetail
}

type NetworkDetail struct {
	Name      string
	Total     int64
	PerSecond int64
}

type Cpu struct {
	Percent    int
	DetailList []*CpuDetail
}

type CpuDetail struct {
	Name    string
	Percent int
}

type Memory struct {
	Unit    string
	Total   int64
	Used    int64
	Percent int
}
