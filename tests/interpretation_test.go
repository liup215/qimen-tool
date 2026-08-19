package tests

import (
	"testing"

	"qimen-tool/internal/interpretation"
	"qimen-tool/internal/plate"
)

func TestInterpret_Career(t *testing.T) {
	tm := beijingTime(2024, 8, 19, 14, 30, 0)
	p, err := plate.BuildPlate(tm)
	if err != nil {
		t.Fatalf("BuildPlate failed: %v", err)
	}

	report, err := interpretation.Interpret(p, "career")
	if err != nil {
		t.Fatalf("Interpret failed: %v", err)
	}

	if report.Meta.QuestionType != "career" {
		t.Errorf("question type = %s, want career", report.Meta.QuestionType)
	}
	if report.Self.Symbol != "日干乙" {
		t.Errorf("self symbol = %s, want 日干乙", report.Self.Symbol)
	}
	if report.Target.Symbol != "开门" {
		t.Errorf("target symbol = %s, want 开门", report.Target.Symbol)
	}
	if report.Self.Palace != 1 || report.Target.Palace != 6 {
		t.Errorf("self palace = %d, target palace = %d, want 1 and 6", report.Self.Palace, report.Target.Palace)
	}
	if report.Relationship.Type != "事生我" {
		t.Errorf("relationship = %s, want 事生我", report.Relationship.Type)
	}
	if len(report.Resources) == 0 {
		t.Error("expected resources, got none")
	}
	if len(report.Threats) == 0 {
		t.Error("expected threats, got none")
	}
	// 该时辰可能未触发具名格局，不强制非空
	if len(report.VetoChecks) == 0 {
		t.Error("expected veto checks, got none")
	}
	if report.Summary == "" {
		t.Error("summary is empty")
	}
}

func TestInterpret_Wealth(t *testing.T) {
	tm := beijingTime(2024, 8, 19, 14, 30, 0)
	p, _ := plate.BuildPlate(tm)
	report, err := interpretation.Interpret(p, "wealth")
	if err != nil {
		t.Fatalf("Interpret failed: %v", err)
	}
	if report.Target.Symbol != "生门" {
		t.Errorf("target symbol = %s, want 生门", report.Target.Symbol)
	}
}

func TestInterpret_Health(t *testing.T) {
	tm := beijingTime(2024, 8, 19, 14, 30, 0)
	p, _ := plate.BuildPlate(tm)
	report, err := interpretation.Interpret(p, "health")
	if err != nil {
		t.Fatalf("Interpret failed: %v", err)
	}
	if report.Target.Symbol != "天芮" {
		t.Errorf("target symbol = %s, want 天芮", report.Target.Symbol)
	}
}

func TestInterpret_UnknownTopic(t *testing.T) {
	tm := beijingTime(2024, 8, 19, 14, 30, 0)
	p, _ := plate.BuildPlate(tm)
	report, err := interpretation.Interpret(p, "unknown")
	if err != nil {
		t.Fatalf("Interpret failed: %v", err)
	}
	if report.Meta.QuestionType != "general" {
		t.Errorf("unknown topic should fallback to general, got %s", report.Meta.QuestionType)
	}
}

// TestInterpret_DayStemJia 验证日干为甲时按旬首遁干定位（甲戌日甲遁于己）
func TestInterpret_DayStemJia(t *testing.T) {
	tm := beijingTime(2026, 8, 28, 15, 30, 0)
	p, err := plate.BuildPlate(tm)
	if err != nil {
		t.Fatalf("BuildPlate failed: %v", err)
	}

	report, err := interpretation.Interpret(p, "travel")
	if err != nil {
		t.Fatalf("Interpret failed: %v", err)
	}

	if report.Self.Symbol != "日干甲（遁己）" {
		t.Errorf("self symbol = %s, want 日干甲（遁己）", report.Self.Symbol)
	}
	if report.Self.Palace != 4 {
		t.Errorf("self palace = %d, want 4", report.Self.Palace)
	}
	if report.Target.Symbol != "景门" {
		t.Errorf("target symbol = %s, want 景门", report.Target.Symbol)
	}
}

// TestInterpret_HourStemJia 验证时干为甲时按旬首遁干定位（甲子时甲遁于戊）
func TestInterpret_HourStemJia(t *testing.T) {
	tm := beijingTime(2026, 8, 28, 0, 0, 0)
	p, err := plate.BuildPlate(tm)
	if err != nil {
		t.Fatalf("BuildPlate failed: %v", err)
	}

	report, err := interpretation.Interpret(p, "general")
	if err != nil {
		t.Fatalf("Interpret failed: %v", err)
	}

	if report.Target.Symbol != "时干甲（遁戊）" {
		t.Errorf("target symbol = %s, want 时干甲（遁戊）", report.Target.Symbol)
	}
	if report.Target.Palace != 7 {
		t.Errorf("target palace = %d, want 7", report.Target.Palace)
	}
}
