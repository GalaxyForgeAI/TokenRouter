package server

import (
	"testing"
	"time"
)

func TestDeriveModelTierByHints(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"gpt-5.6-frontier-ultra", ModelTierFrontier},
		{"claude-opus-4.5", ModelTierFrontier},
		{"gemini-2.5-pro-max", ModelTierFrontier},
		{"gpt-5.6-sol", ModelTierFlagship},
		{"claude-sonnet-4.5", ModelTierFlagship},
		{"deepseek-v3", ModelTierFlagship},
		{"gpt-5.6-nano", ModelTierBasic},
		{"gemini-flash-2.0", ModelTierBasic},
		{"claude-haiku", ModelTierBasic},
		{"qwen-7b", ModelTierStandard},
		{"bge-m3", ModelTierStandard},
	}
	for _, tc := range cases {
		got := deriveModelTier(Model{Name: tc.name})
		if got != tc.want {
			t.Errorf("deriveModelTier(%q) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestDeriveModelTierRespectsExplicitValue(t *testing.T) {
	model := Model{Name: "gpt-5.6-nano", Tier: "flagship"}
	if got := deriveModelTier(model); got != ModelTierFlagship {
		t.Fatalf("explicit tier should win, got %q", got)
	}
	if got := deriveModelTier(Model{Name: "x", Tier: "not-a-tier"}); got != ModelTierStandard {
		t.Fatalf("invalid tier should fall through to derivation, got %q", got)
	}
}

func TestDeriveModelCapabilitiesAppendsMissingTags(t *testing.T) {
	model := Model{Name: "gpt-5.6-vision-reasoning", Capabilities: []string{"chat"}}
	got := deriveModelCapabilities(model)
	if !containsString(got, "vision") {
		t.Errorf("expected vision capability to be appended, got %v", got)
	}
	if !containsString(got, "reasoning") {
		t.Errorf("expected reasoning capability to be appended, got %v", got)
	}
	if !containsString(got, "chat") {
		t.Errorf("existing capabilities must be preserved, got %v", got)
	}
}

func TestDeriveModelCapabilitiesKeepsCatalogTags(t *testing.T) {
	model := Model{Name: "qwen-7b", Capabilities: []string{"chat", "tools"}}
	got := deriveModelCapabilities(model)
	if len(got) != 2 {
		t.Fatalf("no derived tags expected for plain name, got %v", got)
	}
}

func TestApplyDerivedModelProfileFillsFields(t *testing.T) {
	now := time.Now().UTC()
	model := Model{Name: "gpt-5.6-vision-nano", ContextWindow: 128000}
	applyDerivedModelProfile(&model, now)
	if model.Tier != ModelTierBasic {
		t.Errorf("tier = %q, want basic", model.Tier)
	}
	if model.ProfileStatus != ProfileStatusPartial {
		t.Errorf("profile_status = %q, want partial (no pricing)", model.ProfileStatus)
	}
	if model.ProfileUpdated.IsZero() {
		t.Error("profile_updated must be set")
	}
	if !containsString(model.Capabilities, "vision") {
		t.Errorf("capabilities = %v, want vision appended", model.Capabilities)
	}
}

func TestApplyDerivedModelProfileComplete(t *testing.T) {
	model := Model{
		Name:             "claude-sonnet-4.5",
		ContextWindow:    200000,
		InputPriceUSDPer1M: 3,
		Capabilities:     []string{"chat"},
	}
	applyDerivedModelProfile(&model, time.Now().UTC())
	if model.ProfileStatus != ProfileStatusComplete {
		t.Errorf("profile_status = %q, want complete", model.ProfileStatus)
	}
	if model.Tier != ModelTierFlagship {
		t.Errorf("tier = %q, want flagship", model.Tier)
	}
}

func TestNormalizeProfileTier(t *testing.T) {
	if got := normalizeProfileTier("FLAGSHIP"); got != ModelTierFlagship {
		t.Errorf("normalizeProfileTier(FLAGSHIP) = %q", got)
	}
	if got := normalizeProfileTier("unknown"); got != "" {
		t.Errorf("normalizeProfileTier(unknown) = %q, want empty", got)
	}
}
