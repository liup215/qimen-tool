package plate

import (
	"fmt"
	"strings"
	"time"

	"qimen-tool/internal/calendar"
	"qimen-tool/internal/models"
)

// BuildPlate 根据时间构建完整奇门盘面
func BuildPlate(t time.Time) (*models.Plate, error) {
	yearGZ, monthGZ, dayGZ, hourGZ := calendar.SiZhu(t)
	jieQi := calendar.CurrentJieQi(t)
	lunar := calendar.LunarString(t)

	bureau, err := DetermineBureau(t)
	if err != nil {
		return nil, err
	}

	yin := bureau.Dun == "阴遁"

	// 地盘
	earth := BuildEarthPlate(bureau.Ju, yin)

	// 值符值使相关的旬首（时柱旬首）
	dayXunShou := XunShou(dayGZ)
	hourXunShou := XunShou(hourGZ)
	xunStem := xunShouStem(hourXunShou)
	xunShouPalace := palaceOfStem(earth, xunStem)
	if xunShouPalace < 0 {
		return nil, fmt.Errorf("旬首天干 %s 未找到", xunStem)
	}

	// 中五宫寄坤二宫，旬首在 5 宫时按 2 宫处理值符值使
	if xunShouPalace == 5 {
		xunShouPalace = 2
	}

	// 时干所在地盘宫（值符落宫）；中五寄坤
	hourStem, _ := calendar.SplitGanZhi(hourGZ)
	lookupHourStem := hourStem
	if hourStem == "甲" {
		// 甲遁于六仪，用时柱旬首所遁之干定位
		lookupHourStem = xunShouStem(hourXunShou)
	}
	hourStemPalace := palaceOfStem(earth, lookupHourStem)
	if hourStemPalace < 0 {
		return nil, fmt.Errorf("时干 %s（旬首遁%s）未找到", hourStem, lookupHourStem)
	}
	if hourStemPalace == 5 {
		hourStemPalace = 2
	}

	// 时支宫
	_, hourZhi := calendar.SplitGanZhi(hourGZ)
	hourZhiPalace := palaceOfZhi(hourZhi)
	if hourZhiPalace < 0 {
		return nil, fmt.Errorf("时支 %s 未找到", hourZhi)
	}

	// 九星
	stars := BuildStars(xunShouPalace, hourStemPalace, yin)

	// 天盘（基于九星携带原始宫地盘干）
	heaven := BuildHeavenPlate(earth, stars)

	// 八门
	doors, zhiShiPalace := BuildDoors(xunShouPalace, hourXunShou, hourZhi, yin)

	// 八神
	spirits := BuildSpirits(hourStemPalace, yin)

	// 空亡（用时柱旬首）
	kongWang := KongWang(hourXunShou)
	kongWangSet := make(map[string]bool)
	for _, r := range []rune(kongWang) {
		kongWangSet[string(r)] = true
	}

	// 马星：时支马为主流，日支马为辅助参考
	_, dayZhi := calendar.SplitGanZhi(dayGZ)
	_, hourZhi = calendar.SplitGanZhi(hourGZ)
	maXing := MaXing(hourZhi)
	dayMaXing := MaXing(dayZhi)

	plate := &models.Plate{
		SolarTime:      t.Format("2006-01-02 15:04:05"),
		LunarTime:      lunar,
		YearGanZhi:     yearGZ,
		MonthGanZhi:    monthGZ,
		DayGanZhi:      dayGZ,
		HourGanZhi:     hourGZ,
		JieQi:          jieQi,
		Dun:            bureau.Dun,
		Bureau:         bureau.Ju,
		RuleSetVersion: "mainline-cn-v1",
		CenterPalace:   "中五宫寄坤二宫",
		XunShou:        hourXunShou,
		ZhiFuStar:      ZhiFuStar(xunShouPalace),
		ZhiFuPalace:    hourStemPalace,
		ZhiShiDoor:     ZhiShiDoor(xunShouPalace),
		ZhiShiPalace:   zhiShiPalace,
		KongWang:       kongWang,
		MaXing:         maXing,
		DayMaXing:      dayMaXing,
		DoorIndex:      make(map[string]int),
		StarIndex:      make(map[string]int),
		SpiritIndex:    make(map[string]int),
		StemIndex:      make(map[string]int),
	}

	for i := 0; i < 9; i++ {
		pnum := i + 1
		branch := calendar.PalaceBranches[i]
		isKongWang := false
		for _, r := range []rune(branch) {
			if kongWangSet[string(r)] {
				isKongWang = true
				break
			}
		}
		p := models.Palace{
			Number:      pnum,
			Gua:         calendar.PalaceGua[i],
			Branch:      branch,
			Element:     calendar.PalaceElements[i],
			EarthStem:   earth[pnum],
			HeavenStem:  heaven[pnum],
			StemRelation: stemRelation(heaven[pnum], earth[pnum]),
			Door:        doors[pnum],
			Star:        stars[pnum],
			Spirit:      spirits[pnum],
			IsKongWang:  isKongWang,
			HasMaXing:    strings.Contains(branch, maXing),
			HasDayMaXing: strings.Contains(branch, dayMaXing),
		}
		plate.Palaces[i] = p

		// 构建反查索引
		if p.Door != "" {
			plate.DoorIndex[p.Door] = pnum
		}
		if p.Star != "" {
			plate.StarIndex[p.Star] = pnum
		}
		if p.Spirit != "" {
			plate.SpiritIndex[p.Spirit] = pnum
		}
		if p.HeavenStem != "" {
			plate.StemIndex[p.HeavenStem] = pnum
		}
	}

	// 日干、时干落宫
	dayStem, _ := calendar.SplitGanZhi(dayGZ)
	lookupDayStem := dayStem
	if dayStem == "甲" {
		lookupDayStem = xunShouStem(dayXunShou)
	}
	if palace, ok := plate.StemIndex[lookupDayStem]; ok {
		plate.DayStemIndex = palace
	}
	if palace, ok := plate.StemIndex[lookupHourStem]; ok {
		plate.HourStemIndex = palace
	}

	return plate, nil
}

// xunShouStem 返回旬首所隐遁的六仪天干
func xunShouStem(xunShou string) string {
	switch xunShou {
	case "甲子":
		return "戊"
	case "甲戌":
		return "己"
	case "甲申":
		return "庚"
	case "甲午":
		return "辛"
	case "甲辰":
		return "壬"
	case "甲寅":
		return "癸"
	}
	return ""
}

// stemRelation 计算天盘干与地盘干的关系
func stemRelation(heavenStem, earthStem string) string {
	if heavenStem == "" || earthStem == "" {
		return ""
	}
	heavenElement := stemElement(heavenStem)
	earthElement := stemElement(earthStem)
	if heavenElement == "" || earthElement == "" {
		return ""
	}
	if heavenElement == earthElement {
		return "比和"
	}
	// 相生
	generates := map[string]string{"木": "火", "火": "土", "土": "金", "金": "水", "水": "木"}
	if generates[heavenElement] == earthElement {
		return "天生地"
	}
	if generates[earthElement] == heavenElement {
		return "地生天"
	}
	// 相克
	overcomes := map[string]string{"木": "土", "土": "水", "水": "火", "火": "金", "金": "木"}
	if overcomes[heavenElement] == earthElement {
		return "天克地"
	}
	if overcomes[earthElement] == heavenElement {
		return "地克天"
	}
	return "比和"
}

// stemElement 返回天干五行
func stemElement(stem string) string {
	switch stem {
	case "甲", "乙":
		return "木"
	case "丙", "丁":
		return "火"
	case "戊", "己":
		return "土"
	case "庚", "辛":
		return "金"
	case "壬", "癸":
		return "水"
	}
	return ""
}
