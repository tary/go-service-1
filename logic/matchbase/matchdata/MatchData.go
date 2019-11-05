package matchdata

//MatchMode 匹配模式
type MatchMode uint8

const (
	// MatchModeResult 只关心最终匹配结果，不需要中间过程
	MatchModeResult MatchMode = 1
	// MatchModeProgress 需要知道匹配过程（人员进出情况）
	MatchModeProgress MatchMode = 2
)

// Matcher 匹配成员，单人或者队伍
type Matcher struct {
	MatchMode    MatchMode          // 匹配模式（区分是否需要知道匹配过程）
	Key          string             // 匹配成员的唯一标识，要确保不重复
	Num          uint32             // 人数，单人匹配时值为1
	MatchType    string             // 用户自定义的匹配类型，每种类型均有独立的匹配池
	StartTime    int64              // 进入匹配服的时间戳，发起者可以不用设置，匹配服会自动设置
	DoubleParams map[string]float64 // 数值型
	StringParams map[string]string  // 字符串型
	Extension    []byte             // 自定义数据
}

// MatchResult 匹配结果
type MatchResult struct {
	Matchers  []*Matcher //匹配者列表
	Extension []byte     //自定义数据
}
