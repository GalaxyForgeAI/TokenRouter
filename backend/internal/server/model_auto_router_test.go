package server

import (
	"context"
	"testing"
)

func autoTestServer(t *testing.T, models []Model, routes []ModelRoute) *Server {
	t.Helper()
	store := NewMemoryStore()
	for _, model := range models {
		store.AddModel(model)
	}
	for _, route := range routes {
		store.AddRoute(route)
	}
	return New(store)
}

func TestResolveAutoModelPicksCheapestEligible(t *testing.T) {
	server := autoTestServer(t, []Model{
		{Name: "cheap-mini", Tier: ModelTierBasic, InputPriceUSDPer1M: 0.1, OutputPriceUSDPer1M: 0.2, Capabilities: []string{"chat"}},
		{Name: "mid-standard", Tier: ModelTierStandard, InputPriceUSDPer1M: 0.5, OutputPriceUSDPer1M: 1, Capabilities: []string{"chat"}},
		{Name: "flagship-pro", Tier: ModelTierFlagship, InputPriceUSDPer1M: 3, OutputPriceUSDPer1M: 10, Capabilities: []string{"chat", "vision"}},
	}, []ModelRoute{
		{ModelName: "cheap-mini", ProviderID: "prv_a", Status: StatusActive},
		{ModelName: "mid-standard", ProviderID: "prv_a", Status: StatusActive},
		{ModelName: "flagship-pro", ProviderID: "prv_a", Status: StatusActive},
	})

	// Short plain chat must resolve to the cheapest basic model.
	req := &ChatCompletionRequest{Model: AutoModelName, Messages: []ChatMessage{{Role: "user", Content: "hi"}}}
	selected, _, err := server.resolveAutoModel(context.Background(), req)
	if err != nil {
		t.Fatalf("resolveAutoModel: %v", err)
	}
	if selected.Name != "cheap-mini" {
		t.Fatalf("expected cheap-mini, got %s", selected.Name)
	}
}

func TestResolveAutoModelRequiresCapabilities(t *testing.T) {
	server := autoTestServer(t, []Model{
		{Name: "text-only-mini", Tier: ModelTierBasic, InputPriceUSDPer1M: 0.1, Capabilities: []string{"chat"}},
		{Name: "vision-pro", Tier: ModelTierFlagship, InputPriceUSDPer1M: 3, Capabilities: []string{"chat", "vision"}},
	}, []ModelRoute{
		{ModelName: "text-only-mini", ProviderID: "prv_a", Status: StatusActive},
		{ModelName: "vision-pro", ProviderID: "prv_a", Status: StatusActive},
	})

	req := &ChatCompletionRequest{
		Model:    AutoModelName,
		Messages: []ChatMessage{{Role: "user", Content: []any{map[string]any{"type": "text", "text": "describe"}, map[string]any{"type": "image_url", "image_url": map[string]any{"url": "data:image/png;base64,xxx"}}}}},
	}
	selected, _, err := server.resolveAutoModel(context.Background(), req)
	if err != nil {
		t.Fatalf("resolveAutoModel: %v", err)
	}
	if selected.Name != "vision-pro" {
		t.Fatalf("image request must land on a vision model, got %s", selected.Name)
	}
}

func TestResolveAutoModelRequiresTools(t *testing.T) {
	server := autoTestServer(t, []Model{
		{Name: "plain-mini", Tier: ModelTierBasic, InputPriceUSDPer1M: 0.1, Capabilities: []string{"chat"}},
		{Name: "tool-standard", Tier: ModelTierStandard, InputPriceUSDPer1M: 0.4, Capabilities: []string{"chat", "tools"}},
	}, []ModelRoute{
		{ModelName: "plain-mini", ProviderID: "prv_a", Status: StatusActive},
		{ModelName: "tool-standard", ProviderID: "prv_a", Status: StatusActive},
	})

	req := &ChatCompletionRequest{
		Model:    AutoModelName,
		Messages: []ChatMessage{{Role: "user", Content: "call a tool"}},
		Tools:    []any{map[string]any{"type": "function", "function": map[string]any{"name": "lookup"}}},
	}
	selected, _, err := server.resolveAutoModel(context.Background(), req)
	if err != nil {
		t.Fatalf("resolveAutoModel: %v", err)
	}
	if selected.Name != "tool-standard" {
		t.Fatalf("tools request must land on a tools model, got %s", selected.Name)
	}
}

func TestResolveAutoModelEscalatesTierForLongPrompt(t *testing.T) {
	server := autoTestServer(t, []Model{
		{Name: "cheap-mini", Tier: ModelTierBasic, InputPriceUSDPer1M: 0.1, Capabilities: []string{"chat"}},
		{Name: "big-standard", Tier: ModelTierStandard, InputPriceUSDPer1M: 0.5, ContextWindow: 200000, Capabilities: []string{"chat"}},
	}, []ModelRoute{
		{ModelName: "cheap-mini", ProviderID: "prv_a", Status: StatusActive},
		{ModelName: "big-standard", ProviderID: "prv_a", Status: StatusActive},
	})

	longText := ""
	for i := 0; i < 200000; i++ {
		longText += "a"
	}
	req := &ChatCompletionRequest{Model: AutoModelName, Messages: []ChatMessage{{Role: "user", Content: longText}}}
	selected, _, err := server.resolveAutoModel(context.Background(), req)
	if err != nil {
		t.Fatalf("resolveAutoModel: %v", err)
	}
	if selected.Name != "big-standard" {
		t.Fatalf("long prompt must escalate tier, got %s", selected.Name)
	}
}

func TestResolveAutoModelExcludesUnroutedAndDisabled(t *testing.T) {
	server := autoTestServer(t, []Model{
		{Name: "unrouted-mini", Tier: ModelTierBasic, InputPriceUSDPer1M: 0.01, Capabilities: []string{"chat"}},
		{Name: "disabled-mini", Tier: ModelTierBasic, InputPriceUSDPer1M: 0.02, Capabilities: []string{"chat"}, Status: StatusDisabled},
		{Name: "routed-mini", Tier: ModelTierBasic, InputPriceUSDPer1M: 0.05, Capabilities: []string{"chat"}},
	}, []ModelRoute{
		{ModelName: "routed-mini", ProviderID: "prv_a", Status: StatusActive},
	})

	req := &ChatCompletionRequest{Model: AutoModelName, Messages: []ChatMessage{{Role: "user", Content: "hi"}}}
	selected, _, err := server.resolveAutoModel(context.Background(), req)
	if err != nil {
		t.Fatalf("resolveAutoModel: %v", err)
	}
	if selected.Name != "routed-mini" {
		t.Fatalf("expected routed-mini (only eligible), got %s", selected.Name)
	}
}

func TestResolveAutoModelNoEligibleError(t *testing.T) {
	server := autoTestServer(t, []Model{
		{Name: "no-route-mini", Tier: ModelTierBasic, InputPriceUSDPer1M: 0.1, Capabilities: []string{"chat"}},
	}, nil)

	req := &ChatCompletionRequest{Model: AutoModelName, Messages: []ChatMessage{{Role: "user", Content: "hi"}}}
	if _, _, err := server.resolveAutoModel(context.Background(), req); err == nil {
		t.Fatal("expected error when no model has an active route")
	}
}
