package tests

import (
	"reflect"
	"testing"

	"qimen-tool/internal/plate"
)

// 验证地盘按戊己庚辛壬癸丁丙乙排列且覆盖 1-9 宫
func TestEarthPlate_Order(t *testing.T) {
	yang := plate.BuildEarthPlate(1, false)
	expectedYang := map[int]string{
		1: "戊", 2: "己", 3: "庚", 4: "辛", 5: "壬", 6: "癸", 7: "丁", 8: "丙", 9: "乙",
	}
	if !reflect.DeepEqual(yang, expectedYang) {
		t.Errorf("阳遁 1 局地盘错误: got %v, want %v", yang, expectedYang)
	}

	yin := plate.BuildEarthPlate(9, true)
	expectedYin := map[int]string{
		9: "戊", 8: "己", 7: "庚", 6: "辛", 5: "壬", 4: "癸", 3: "丁", 2: "丙", 1: "乙",
	}
	if !reflect.DeepEqual(yin, expectedYin) {
		t.Errorf("阴遁 9 局地盘错误: got %v, want %v", yin, expectedYin)
	}
}

// 验证八门相对顺序固定：休→生→伤→杜→景→死→惊→开
// 阳遁从 1 起，按洛书圆周顺时针：1,8,3,4,9,2,7,6（跳中五）
func TestDoors_RelativeOrder(t *testing.T) {
	doors, _ := plate.BuildDoors(1, "甲子", "子", false)
	order := []string{"休", "生", "伤", "杜", "景", "死", "惊", "开"}
	positions := []int{1, 8, 3, 4, 9, 2, 7, 6}
	for i, p := range positions {
		want := order[i]
		if doors[p] != want {
			t.Errorf("八门顺序错误: %d宫 got %s, want %s", p, doors[p], want)
		}
	}
}

// 验证九星相对顺序
func TestStars_RelativeOrder(t *testing.T) {
	stars := plate.BuildStars(1, 1, false)
	order := []string{"天蓬", "天任", "天冲", "天辅", "天英", "天芮", "天柱", "天心"}
	positions := []int{1, 8, 3, 4, 9, 2, 7, 6}
	for i, p := range positions {
		want := order[i]
		if stars[p] != want {
			t.Errorf("九星顺序错误: %d宫 got %s, want %s", p, stars[p], want)
		}
	}
}

// 验证八神阳遁顺序
func TestSpirits_YangOrder(t *testing.T) {
	spirits := plate.BuildSpirits(1, false)
	order := []string{"值符", "螣蛇", "太阴", "六合", "白虎", "玄武", "九地", "九天"}
	positions := []int{1, 8, 3, 4, 9, 2, 7, 6}
	for i, p := range positions {
		want := order[i]
		if spirits[p] != want {
			t.Errorf("阳遁八神顺序错误: %d宫 got %s, want %s", p, spirits[p], want)
		}
	}
}

// 验证八神阴遁顺序（从 1 开始逆排）
func TestSpirits_YinOrder(t *testing.T) {
	spirits := plate.BuildSpirits(1, true)
	order := []string{"值符", "螣蛇", "太阴", "六合", "白虎", "玄武", "九地", "九天"}
	positions := []int{1, 6, 7, 2, 9, 4, 3, 8}
	for i, p := range positions {
		want := order[i]
		if spirits[p] != want {
			t.Errorf("阴遁八神顺序错误: %d宫 got %s, want %s", p, spirits[p], want)
		}
	}
}

// 验证旬首在中五宫时，值符值使按坤二宫处理
func TestXunShouInCenter(t *testing.T) {
	// 阳遁 4 局：4戊 5己 6庚...，旬首甲戌己在 5 宫，寄坤 2
	// 值使门应为坤 2 原始门 "死"
	doors, zhiShiPalace := plate.BuildDoors(5, "甲戌", "戌", false)
	if zhiShiPalace != 2 {
		t.Errorf("旬首在中五宫时值使门落宫 = %d, want 2（中五寄坤二）", zhiShiPalace)
	}
	if doors[2] != "死" {
		t.Errorf("旬首在中五宫时坤二宫门 = %s, want 死", doors[2])
	}
}
