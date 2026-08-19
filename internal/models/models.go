package models

// Plate 表示一个完整的奇门遁甲盘面
type Plate struct {
	SolarTime      string         `json:"solar_time"`      // 阳历时间，如 2024-08-19 14:30
	LunarTime      string         `json:"lunar_time"`      // 农历时间描述
	YearGanZhi     string         `json:"year_ganzhi"`     // 年柱
	MonthGanZhi    string         `json:"month_ganzhi"`    // 月柱
	DayGanZhi      string         `json:"day_ganzhi"`      // 日柱
	HourGanZhi     string         `json:"hour_ganzhi"`     // 时柱
	JieQi          string         `json:"jie_qi"`          // 当前节气
	Dun            string         `json:"dun"`             // 阳遁 / 阴遁
	Bureau         int            `json:"bureau"`          // 局数 1-9
	RuleSetVersion string         `json:"rule_set_version"`// 规则集版本，如 mainline-cn-v1
	CenterPalace   string         `json:"center_palace"`   // 中宫处理规则说明
	XunShou        string         `json:"xun_shou"`        // 旬首
	ZhiFuStar      string         `json:"zhi_fu_star"`     // 值符星
	ZhiFuPalace    int            `json:"zhi_fu_palace"`   // 值符落宫
	ZhiShiDoor     string         `json:"zhi_shi_door"`    // 值使门
	ZhiShiPalace   int            `json:"zhi_shi_palace"`  // 值使落宫
	KongWang       string         `json:"kong_wang"`       // 空亡
	MaXing         string         `json:"ma_xing"`         // 马星
	DayStemIndex   int            `json:"day_stem_index"`  // 日干落宫
	HourStemIndex  int            `json:"hour_stem_index"` // 时干落宫
	DoorIndex      map[string]int `json:"door_index"`      // 八门 → 宫
	StarIndex      map[string]int `json:"star_index"`      // 九星 → 宫
	SpiritIndex    map[string]int `json:"spirit_index"`    // 八神 → 宫
	StemIndex      map[string]int `json:"stem_index"`      // 天盘天干 → 宫
	Palaces        [9]Palace      `json:"palaces"`         // 九宫，索引 0 对应坎一宫，按后天八卦顺序
}

// Palace 表示一个宫位
type Palace struct {
	Number       int    `json:"number"`        // 宫号 1-9
	Gua          string `json:"gua"`           // 八卦
	Branch       string `json:"branch"`        // 地支
	Element      string `json:"element"`       // 宫位五行
	EarthStem    string `json:"earth_stem"`    // 地盘天干（三奇六仪）
	HeavenStem   string `json:"heaven_stem"`   // 天盘天干
	StemRelation string `json:"stem_relation"` // 天盘干与地盘干关系：天生地/地生天/天克地/地克天/比和
	Door         string `json:"door"`          // 八门
	Star         string `json:"star"`          // 九星
	Spirit       string `json:"spirit"`        // 八神
	IsKongWang   bool   `json:"is_kong_wang"`  // 是否空亡
	HasMaXing    bool   `json:"has_ma_xing"`   // 是否马星
}

// GanZhi 表示一个干支
type GanZhi struct {
	Gan string `json:"gan"`
	Zhi string `json:"zhi"`
}

func (gz GanZhi) String() string {
	return gz.Gan + gz.Zhi
}
