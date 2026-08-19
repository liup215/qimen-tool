package interpretation

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"qimen-tool/internal/calendar"
	"qimen-tool/internal/models"
)

//go:embed doors.json
var doorsJSON []byte

//go:embed stars.json
var starsJSON []byte

//go:embed spirits.json
var spiritsJSON []byte

//go:embed yong_shen.json
var yongShenJSON []byte

//go:embed patterns.json
var patternsJSON []byte

// RuleSet 规则集
type RuleSet struct {
	Doors    map[string]DoorStarSpirit
	Stars    map[string]DoorStarSpirit
	Spirits  map[string]DoorStarSpirit
	YongShen map[string]YongShenRule
	Patterns PatternRule
}

type DoorStarSpirit struct {
	WuXing  string `json:"五行,omitempty"`
	JiXiong string `json:"吉凶"`
	XiangYi string `json:"象义"`
}

type YongShenRule struct {
	Description string   `json:"description"`
	YongShen    []string `json:"yong_shen"`
}

type PatternRule struct {
	WuXing          map[string]string       `json:"wu_xing"`
	RuMu            map[string]string       `json:"ru_mu"`
	JiXing          map[string]string       `json:"ji_xing"`
	SanqiRuMu       map[string]string       `json:"sanqi_ru_mu"`
	KongWangMeaning string                  `json:"kong_wang_meaning"`
	MaXingMeaning   string                  `json:"ma_xing_meaning"`
	FuYin           string                  `json:"fu_yin"`
	FanYin          string                  `json:"fan_yin"`
	WangXiang       WangXiangRule           `json:"wang_xiang"`
	NamedPatterns   map[string]NamedPattern `json:"named_patterns"`
	JiugongXiangyi  map[string]string       `json:"jiugong_xiangyi"`
	SanqiXiangyi    map[string]string       `json:"sanqi_xiangyi"`
}

type NamedPattern struct {
	Heaven      string `json:"heaven"`
	Earth       string `json:"earth"`
	Nature      string `json:"nature"`
	Description string `json:"description"`
}

type WangXiangRule struct {
	Description string              `json:"description"`
	Seasons     map[string][]string `json:"seasons"`
}

var rules *RuleSet

func init() {
	rules = &RuleSet{
		Doors:    make(map[string]DoorStarSpirit),
		Stars:    make(map[string]DoorStarSpirit),
		Spirits:  make(map[string]DoorStarSpirit),
		YongShen: make(map[string]YongShenRule),
	}
	_ = json.Unmarshal(doorsJSON, &rules.Doors)
	_ = json.Unmarshal(starsJSON, &rules.Stars)
	_ = json.Unmarshal(spiritsJSON, &rules.Spirits)
	_ = json.Unmarshal(yongShenJSON, &rules.YongShen)
	_ = json.Unmarshal(patternsJSON, &rules.Patterns)
}

// Rules 返回规则集（供 rules 子命令使用）
func Rules() *RuleSet {
	return rules
}

// YongShenTopics 返回所有用神主题 key
func YongShenTopics() []string {
	keys := make([]string, 0, len(rules.YongShen))
	for k := range rules.YongShen {
		keys = append(keys, k)
	}
	return keys
}

// Report 解释报告
type Report struct {
	Meta             ReportMeta      `json:"meta"`
	Self             SubjectState    `json:"self"`
	Target           SubjectState    `json:"target"`
	Relationship     Relationship    `json:"relationship"`
	Resources        []Force         `json:"resources"`
	Threats          []Force         `json:"threats"`
	DetectedPatterns []PatternMatch  `json:"detected_patterns"`
	VetoChecks       []VetoCheck     `json:"veto_checks"`
	Summary          string          `json:"summary"`
	Caution          string          `json:"caution"`
}

type PatternMatch struct {
	Name        string `json:"name"`
	Palace      int    `json:"palace"`
	Nature      string `json:"nature"`
	Description string `json:"description"`
}

type VetoCheck struct {
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Message string `json:"message"`
}

type ReportMeta struct {
	Time         string `json:"time"`
	QuestionType string `json:"question_type"`
	Description  string `json:"description"`
}

type SubjectState struct {
	Symbol string   `json:"symbol"`
	Palace int      `json:"palace"`
	Level  string   `json:"level"`
	Notes  []string `json:"notes"`
}

type Relationship struct {
	Type        string `json:"type"`
	Effect      string `json:"effect"`
	Description string `json:"description"`
}

type Force struct {
	Symbol         string `json:"symbol"`
	Palace         int    `json:"palace"`
	RelationToSelf string `json:"relation_to_self"`
	Effect         string `json:"effect"`
	Note           string `json:"note"`
}

// Interpret 根据盘面和问题类型生成解释报告
func Interpret(plate *models.Plate, topic string) (*Report, error) {
	ys, ok := rules.YongShen[topic]
	if !ok {
		ys = rules.YongShen["general"]
		topic = "general"
	}

	monthZhi := ""
	if len(plate.MonthGanZhi) >= 2 {
		_, monthZhi = calendar.SplitGanZhi(plate.MonthGanZhi)
	}

	dayGan, _ := calendar.SplitGanZhi(plate.DayGanZhi)
	selfStem := dayGan
	selfSymbol := "日干" + dayGan
	if dayGan == "甲" {
		xunShou := calendar.XunShou(plate.DayGanZhi)
		selfStem = calendar.DunGan(xunShou)
		selfSymbol = fmt.Sprintf("日干甲（遁%s）", selfStem)
	}
	selfPalace := findStemPalace(plate, selfStem)
	if selfPalace < 0 {
		return nil, fmt.Errorf("日干 %s 未找到落宫", selfSymbol)
	}

	targetSymbol, targetPalace := resolveTarget(plate, topic)
	if targetPalace < 0 {
		return nil, fmt.Errorf("用神 %s 未找到落宫", targetSymbol)
	}

	selfP := plate.Palaces[selfPalace-1]
	targetP := plate.Palaces[targetPalace-1]

	selfState := analyzeSubject(selfP, selfSymbol, monthZhi)
	targetState := analyzeSubject(targetP, targetSymbol, monthZhi)

	rel := analyzeRelationship(selfP, targetP)

	resources := scanForces(plate, selfPalace, []string{"值符", "太阴", "六合", "生门", "九天"}, true)
	threats := scanForces(plate, selfPalace, []string{"白虎", "玄武", "螣蛇", "惊门", "伤门"}, false)

	patterns := detectPatterns(plate)
	vetos := detectVetoChecks(plate)

	report := &Report{
		Meta: ReportMeta{
			Time:         plate.SolarTime,
			QuestionType: topic,
			Description:  ys.Description,
		},
		Self:             selfState,
		Target:           targetState,
		Relationship:     rel,
		Resources:        resources,
		Threats:          threats,
		DetectedPatterns: patterns,
		VetoChecks:       vetos,
		Caution:          "本结果仅供学习参考，不替代专业判断。疾病、法律、财务等重大决策请咨询现实专业人士。",
	}

	report.Summary = buildSummary(report, ys)
	return report, nil
}

// resolveTarget 根据问题类型确定用神符号和落宫
func resolveTarget(plate *models.Plate, topic string) (string, int) {
	switch topic {
	case "career":
		return "开门", findDoorPalace(plate, "开")
	case "wealth":
		return "生门", findDoorPalace(plate, "生")
	case "marriage":
		return "六合", findSpiritPalace(plate, "六合")
	case "health":
		return "天芮", findStarPalace(plate, "天芮")
	case "lawsuit":
		return "惊门", findDoorPalace(plate, "惊")
	case "travel":
		return "景门", findDoorPalace(plate, "景")
	case "study":
		return "天辅", findStarPalace(plate, "天辅")
	case "lost", "general":
		return resolveHourStemTarget(plate)
	}
	// 默认用时干
	return resolveHourStemTarget(plate)
}

// resolveHourStemTarget 解析时干用神；当时干为甲时按旬首遁干定位
func resolveHourStemTarget(plate *models.Plate) (string, int) {
	hourGan, _ := calendar.SplitGanZhi(plate.HourGanZhi)
	if hourGan == "甲" {
		xunShou := calendar.XunShou(plate.HourGanZhi)
		dun := calendar.DunGan(xunShou)
		return fmt.Sprintf("时干甲（遁%s）", dun), findStemPalace(plate, dun)
	}
	return "时干" + hourGan, findStemPalace(plate, hourGan)
}

// findStemPalace 查找某天干在天盘上的落宫
func findStemPalace(plate *models.Plate, stem string) int {
	for _, p := range plate.Palaces {
		if p.HeavenStem == stem {
			return p.Number
		}
	}
	return -1
}

// findDoorPalace 查找某门落宫
func findDoorPalace(plate *models.Plate, door string) int {
	for _, p := range plate.Palaces {
		if p.Door == door {
			return p.Number
		}
	}
	return -1
}

// findStarPalace 查找某星落宫
func findStarPalace(plate *models.Plate, star string) int {
	for _, p := range plate.Palaces {
		if p.Star == star || (star == "天芮" && p.Star == "天禽") {
			return p.Number
		}
	}
	return -1
}

// findSpiritPalace 查找某神落宫
func findSpiritPalace(plate *models.Plate, spirit string) int {
	for _, p := range plate.Palaces {
		if p.Spirit == spirit {
			return p.Number
		}
	}
	return -1
}

// analyzeSubject 分析某一宫的状态
func analyzeSubject(p models.Palace, symbol, monthZhi string) SubjectState {
	notes := []string{}

	// 门星神基础信息
	if p.Door != "" {
		notes = append(notes, fmt.Sprintf("临%s门", p.Door))
	}
	if p.Star != "" {
		notes = append(notes, fmt.Sprintf("临%s星", p.Star))
	}
	if p.Spirit != "" {
		notes = append(notes, fmt.Sprintf("临%s", p.Spirit))
	}

	// 特殊状态
	if p.IsKongWang {
		notes = append(notes, "空亡")
	}
	if p.HasMaXing {
		notes = append(notes, "时马")
	}
	if p.HasDayMaXing {
		notes = append(notes, "日马")
	}
	if palaceHasZhi(p, rules.Patterns.RuMu[p.HeavenStem]) {
		notes = append(notes, fmt.Sprintf("%s入墓", p.HeavenStem))
	}
	if palaceHasZhi(p, rules.Patterns.JiXing[p.HeavenStem]) {
		notes = append(notes, fmt.Sprintf("%s击刑", p.HeavenStem))
	}

	level := "平"
	if monthZhi != "" {
		season := seasonOf(monthZhi)
		element := rules.Patterns.WuXing[p.HeavenStem]
		level = wuXingStatus(element, season)
	}

	return SubjectState{
		Symbol: symbol,
		Palace: p.Number,
		Level:  level,
		Notes:  notes,
	}
}

// analyzeRelationship 分析日干宫与用神宫的生克关系
func analyzeRelationship(selfP, targetP models.Palace) Relationship {
	rel := wuXingRelation(selfP.Element, targetP.Element)
	var relType, effect, desc string
	switch rel {
	case "生":
		relType = "事生我"
		effect = "positive"
		desc = fmt.Sprintf("%s宫（%s）生%s宫（%s），所问之事对你有助益，机会偏正面。", targetP.Gua, targetP.Element, selfP.Gua, selfP.Element)
	case "被生":
		relType = "我生事"
		effect = "neutral"
		desc = fmt.Sprintf("%s宫（%s）生%s宫（%s），需要你主动投入、付出，事成在你。", selfP.Gua, selfP.Element, targetP.Gua, targetP.Element)
	case "克":
		relType = "事克我"
		effect = "negative"
		desc = fmt.Sprintf("%s宫（%s）克%s宫（%s），事情对你有压力、阻力或风险。", targetP.Gua, targetP.Element, selfP.Gua, selfP.Element)
	case "被克":
		relType = "我克事"
		effect = "neutral"
		desc = fmt.Sprintf("%s宫（%s）克%s宫（%s），你能主导此事，但需费心费力。", selfP.Gua, selfP.Element, targetP.Gua, targetP.Element)
	default:
		relType = "比和"
		effect = "neutral"
		desc = fmt.Sprintf("%s宫与%s宫同属%s，双方势均力敌，结果看各自旺衰。", selfP.Gua, targetP.Gua, selfP.Element)
	}

	// 修正：若用神空亡或入墓，生克力量减弱
	if targetP.IsKongWang {
		desc += " 用神逢空亡，力量减半，事情不实或待时机。"
	}

	return Relationship{
		Type:        relType,
		Effect:      effect,
		Description: desc,
	}
}

// scanForces 扫描全盘的资源或威胁符号
func scanForces(plate *models.Plate, selfPalace int, symbols []string, isResource bool) []Force {
	var forces []Force
	selfP := plate.Palaces[selfPalace-1]

	for _, symbol := range symbols {
		var p models.Palace
		found := false
		for _, palace := range plate.Palaces {
			if (isDoor(symbol) && palace.Door == symbol) ||
				(isStar(symbol) && (palace.Star == symbol || (symbol == "天芮" && palace.Star == "天禽"))) ||
				(isSpirit(symbol) && palace.Spirit == symbol) {
				p = palace
				found = true
				break
			}
		}
		if !found {
			continue
		}

		rel := wuXingRelation(selfP.Element, p.Element)
		var relation, effect, note string
		switch rel {
		case "生":
			relation = "生我"
			if isResource {
				effect = "positive"
				note = fmt.Sprintf("%s落%s宫，生扶日干宫，有助力可借助", symbol, p.Gua)
			} else {
				effect = "negative"
				note = fmt.Sprintf("%s落%s宫，生扶日干宫，表面缓和但隐患暗藏", symbol, p.Gua)
			}
		case "被生":
			relation = "我生"
			effect = "neutral"
			note = fmt.Sprintf("%s落%s宫，被日干宫所生", symbol, p.Gua)
		case "克":
			relation = "克我"
			if isResource {
				effect = "negative"
				note = fmt.Sprintf("%s落%s宫，克日干宫，助力带有压力", symbol, p.Gua)
			} else {
				effect = "negative"
				note = fmt.Sprintf("%s落%s宫，克日干宫，需防范", symbol, p.Gua)
			}
		case "被克":
			relation = "我克"
			effect = "positive"
			if isResource {
				note = fmt.Sprintf("%s落%s宫，被日干宫所制，助力需主动争取", symbol, p.Gua)
			} else {
				note = fmt.Sprintf("%s落%s宫，被日干宫所制，威胁可控", symbol, p.Gua)
			}
		default:
			relation = "比和"
			effect = "neutral"
			note = fmt.Sprintf("%s落%s宫，与日干宫比和", symbol, p.Gua)
		}

		if p.IsKongWang {
			note += "，但该宫逢空亡，力量不实"
		}

		forces = append(forces, Force{
			Symbol:         symbol,
			Palace:         p.Number,
			RelationToSelf: relation,
			Effect:         effect,
			Note:           note,
		})
	}

	return forces
}

func isDoor(s string) bool {
	_, ok := rules.Doors[s]
	return ok
}

func isStar(s string) bool {
	_, ok := rules.Stars[s]
	return ok
}

func isSpirit(s string) bool {
	_, ok := rules.Spirits[s]
	return ok
}

// detectPatterns 检测具名格局
func detectPatterns(plate *models.Plate) []PatternMatch {
	var matches []PatternMatch
	for name, pat := range rules.Patterns.NamedPatterns {
		for _, palace := range plate.Palaces {
			if palace.HeavenStem == pat.Heaven && palace.EarthStem == pat.Earth {
				matches = append(matches, PatternMatch{
					Name:        name,
					Palace:      palace.Number,
					Nature:      pat.Nature,
					Description: pat.Description,
				})
			}
		}
	}
	return matches
}

// detectVetoChecks 一票否决检测
func detectVetoChecks(plate *models.Plate) []VetoCheck {
	var checks []VetoCheck

	// 五不遇时：时干克日干，且阴阳相同
	dayGan, _ := calendar.SplitGanZhi(plate.DayGanZhi)
	hourGan, _ := calendar.SplitGanZhi(plate.HourGanZhi)
	dayElement := rules.Patterns.WuXing[dayGan]
	hourElement := rules.Patterns.WuXing[hourGan]
	if dayElement != "" && hourElement != "" {
		overcomes := map[string]string{"木": "土", "土": "水", "水": "火", "火": "金", "金": "木"}
		dayYin := isYinGan(dayGan)
		hourYin := isYinGan(hourGan)
		if overcomes[hourElement] == dayElement && dayYin == hourYin {
			checks = append(checks, VetoCheck{
				Name:    "五不遇时",
				Active:  true,
				Message: "时干克日干且阴阳相同，主事多阻碍，宜静不宜动。",
			})
		}
	}

	// 三奇入墓
	for _, palace := range plate.Palaces {
		if mu, ok := rules.Patterns.SanqiRuMu[palace.HeavenStem]; ok {
			if palaceHasZhi(palace, mu) {
				checks = append(checks, VetoCheck{
					Name:    "三奇入墓",
					Active:  true,
					Message: fmt.Sprintf("%s奇入墓于%d宫，希望受阻，难以发挥。", palace.HeavenStem, palace.Number),
				})
			}
		}
	}

	// 天网四张：六癸临某宫（癸+癸）
	for _, palace := range plate.Palaces {
		if palace.HeavenStem == "癸" && palace.EarthStem == "癸" {
			checks = append(checks, VetoCheck{
				Name:    "天网四张",
				Active:  true,
				Message: fmt.Sprintf("癸加癸临%d宫，主困顿、陷阱、难以脱身。", palace.Number),
			})
		}
	}

	// 全局伏吟
	allFuYin := true
	for _, palace := range plate.Palaces {
		if palace.Number == 5 {
			continue
		}
		if palace.HeavenStem != palace.EarthStem {
			allFuYin = false
			break
		}
	}
	if allFuYin {
		checks = append(checks, VetoCheck{
			Name:    "全局伏吟",
			Active:  true,
			Message: "全局伏吟，事体迟滞反复，宜守不宜攻。",
		})
	}

	// 全局反吟：天地盘相冲（地支相隔六位）
	allFanYin := true
	for _, palace := range plate.Palaces {
		if palace.Number == 5 {
			continue
		}
		if palace.StemRelation != "地克天" && palace.StemRelation != "天克地" {
			allFanYin = false
			break
		}
	}
	if allFanYin {
		checks = append(checks, VetoCheck{
			Name:    "全局反吟",
			Active:  true,
			Message: "全局反吟，事体来去迅速、变动剧烈。",
		})
	}

	if len(checks) == 0 {
		checks = append(checks, VetoCheck{
			Name:    "一票否决检查",
			Active:  false,
			Message: "未检测到明显一票否决格局。",
		})
	}

	return checks
}

func isYinGan(gan string) bool {
	switch gan {
	case "乙", "丁", "己", "辛", "癸":
		return true
	}
	return false
}

// buildSummary 生成综合小结
func buildSummary(r *Report, ys YongShenRule) string {
	selfStrong := r.Self.Level == "旺" || r.Self.Level == "相"
	targetStrong := r.Target.Level == "旺" || r.Target.Level == "相"

	parts := []string{
		fmt.Sprintf("问事主题：%s。", ys.Description),
		fmt.Sprintf("自身（%s）落%d宫，状态%s；所问之事（%s）落%d宫，状态%s。", r.Self.Symbol, r.Self.Palace, r.Self.Level, r.Target.Symbol, r.Target.Palace, r.Target.Level),
	}

	if r.Relationship.Effect == "positive" && targetStrong {
		parts = append(parts, "事体旺相对你有利，机会较好。")
	} else if r.Relationship.Effect == "negative" && targetStrong {
		parts = append(parts, "事体旺相且克你，阻力较大，需谨慎。")
	} else {
		parts = append(parts, r.Relationship.Description)
	}

	if len(r.Resources) > 0 {
		parts = append(parts, fmt.Sprintf("资源符号%d个，需关注是否有力。", len(r.Resources)))
	}
	if len(r.Threats) > 0 {
		parts = append(parts, fmt.Sprintf("威胁符号%d个，需关注是否克我。", len(r.Threats)))
	}

	// 格局提醒
	for _, pat := range r.DetectedPatterns {
		if pat.Nature == "吉" {
			parts = append(parts, fmt.Sprintf("%s临%d宫（%s），可资利用。", pat.Name, pat.Palace, pat.Nature))
		} else {
			parts = append(parts, fmt.Sprintf("%s临%d宫（%s），需加留意。", pat.Name, pat.Palace, pat.Nature))
		}
	}

	// 一票否决提醒
	for _, veto := range r.VetoChecks {
		if veto.Active {
			parts = append(parts, fmt.Sprintf("【一票否决】%s：%s", veto.Name, veto.Message))
		}
	}

	if !selfStrong {
		parts = append(parts, "自身状态不旺，宜稳守，不宜妄动。")
	}

	return joinStrings(parts, " ")
}

func joinStrings(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}

// palaceHasZhi 判断宫位是否包含某地支
func palaceHasZhi(p models.Palace, zhi string) bool {
	if zhi == "" {
		return false
	}
	for _, r := range []rune(p.Branch) {
		if string(r) == zhi {
			return true
		}
	}
	return false
}

// seasonOf 根据月支判断季节
func seasonOf(zhi string) string {
	for season, zhis := range rules.Patterns.WangXiang.Seasons {
		for _, z := range zhis {
			if z == zhi {
				return season
			}
		}
	}
	return ""
}

// wuXingStatus 根据季节判断五行旺衰，返回旺/相/休/囚/死
func wuXingStatus(element, season string) string {
	switch season {
	case "春":
		switch element {
		case "木":
			return "旺"
		case "火":
			return "相"
		case "水":
			return "休"
		case "金":
			return "囚"
		case "土":
			return "死"
		}
	case "夏":
		switch element {
		case "火":
			return "旺"
		case "土":
			return "相"
		case "木":
			return "休"
		case "水":
			return "囚"
		case "金":
			return "死"
		}
	case "秋":
		switch element {
		case "金":
			return "旺"
		case "水":
			return "相"
		case "土":
			return "休"
		case "火":
			return "囚"
		case "木":
			return "死"
		}
	case "冬":
		switch element {
		case "水":
			return "旺"
		case "木":
			return "相"
		case "金":
			return "休"
		case "土":
			return "囚"
		case "火":
			return "死"
		}
	}
	return "平"
}

// wuXingRelation 返回 self 对 target 的五行关系：生/被生/克/被克/同
func wuXingRelation(selfElement, targetElement string) string {
	if selfElement == targetElement {
		return "同"
	}
	// 五行相生：木->火->土->金->水->木
	generates := map[string]string{
		"木": "火",
		"火": "土",
		"土": "金",
		"金": "水",
		"水": "木",
	}
	if generates[selfElement] == targetElement {
		return "被生" // 我生它
	}
	if generates[targetElement] == selfElement {
		return "生" // 它生我
	}
	// 五行相克：木->土->水->火->金->木
	overcomes := map[string]string{
		"木": "土",
		"土": "水",
		"水": "火",
		"火": "金",
		"金": "木",
	}
	if overcomes[selfElement] == targetElement {
		return "被克" // 我克它
	}
	if overcomes[targetElement] == selfElement {
		return "克" // 它克我
	}
	return "同"
}
