package types

type SearchInfo struct {
	StartTime int64
	StopTime  int64
	Depth     int
	DepthSet  bool
	TimeSet   bool
	MovesToGo int
	Infinite  bool

	Nodes int64

	Quit    bool
	Stopped bool

	Fh  float64
	Fhf float64

	NullCut int
}
