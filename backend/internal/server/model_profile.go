package server

import (
	"strings"
	"time"
)

// Capability tiers used by automatic model routing and profile building.
// The tier is the quality anchor: a query's minimal requirement is matched
// against these tiers before cost is optimized within the qualified set.
const (
	ModelTierBasic    = "basic"
	ModelTierStandard = "standard"
	ModelTierFlagship = "flagship"
	ModelTierFrontier = "frontier"
)

// Profile status marks how complete a model profile is.
const (
	ProfileStatusComplete = "complete"
	ProfileStatusPartial  = "partial"
)

// ValidModelTiers is the closed set of recognized tiers.
var ValidModelTiers = []string{ModelTierBasic, ModelTierStandard, ModelTierFlagship, ModelTierFrontier}

// frontierNameHints push a model to the frontier tier: the very top of the
// flagship line. Matched case-insensitively against the model name.
var frontierNameHints = []string{"frontier", "opus", "ultra", "max", "pro-max"}

// flagshipNameHints mark a model as flagship-tier.
var flagshipNameHints = []string{"pro", "plus", "premium", "large", "sonnet", "flagship", "-4.5", "-5.", "gpt-5", "o3", "o4", "r1", "deepseek-v3"}

// basicNameHints mark a model as basic-tier (small, cheap, fast).
var basicNameHints = []string{"mini", "flash", "nano", "lite", "small", "haiku", "tiny", "-1.5"}

// reasoningNameHints add the reasoning capability when absent.
var reasoningNameHints = []string{"reason", "think", "thinking", "deep", "r1", "o1", "o3", "o4", "deepseek-reasoner"}

// visionNameHints add the vision capability when absent and input modalities
// include images.
var visionNameHints = []string{"vision", "multimodal", "omni"}

// normalizeProfileTier maps any provided tier onto the closed set, or returns
// the empty string when it is not a recognized tier.
func normalizeProfileTier(tier string) string {
	tier = strings.ToLower(strings.TrimSpace(tier))
	for _, valid := range ValidModelTiers {
		if tier == valid {
			return tier
		}
	}
	return ""
}

// deriveModelTier infers a capability tier from the model name and the
// catalog-provided capabilities. Name hints are heuristics; the catalog tier
// field and explicit admin edits always win.
func deriveModelTier(model Model) string {
	if tier := normalizeProfileTier(model.Tier); tier != "" {
		return tier
	}
	name := strings.ToLower(model.Name)
	for _, hint := range frontierNameHints {
		if strings.Contains(name, hint) {
			return ModelTierFrontier
		}
	}
	// Small-model suffixes (mini/flash/nano/lite/...) are the strongest tier
	// signal and must win over family prefixes such as "gpt-5".
	for _, hint := range basicNameHints {
		if strings.Contains(name, hint) {
			return ModelTierBasic
		}
	}
	for _, hint := range flagshipNameHints {
		if strings.Contains(name, hint) {
			return ModelTierFlagship
		}
	}
	return ModelTierStandard
}

// deriveModelCapabilities fills capability gaps from name heuristics. Existing
// catalog capabilities are preserved; only missing derived tags are appended.
func deriveModelCapabilities(model Model) []string {
	capabilities := append([]string(nil), model.Capabilities...)
	name := strings.ToLower(model.Name)
	has := func(tag string) bool {
		for _, capability := range capabilities {
			if capability == tag {
				return true
			}
		}
		return false
	}
	wantsVision := false
	for _, hint := range visionNameHints {
		if strings.Contains(name, hint) {
			wantsVision = true
			break
		}
	}
	if wantsVision && !has("vision") {
		capabilities = append(capabilities, "vision")
	}
	for _, hint := range reasoningNameHints {
		if strings.Contains(name, hint) && !has("reasoning") {
			capabilities = append(capabilities, "reasoning")
			break
		}
	}
	return capabilities
}

// normalizeProfileStatus maps any provided status onto the closed set, or
// returns the empty string when it is not recognized.
func normalizeProfileStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case ProfileStatusComplete:
		return ProfileStatusComplete
	case ProfileStatusPartial:
		return ProfileStatusPartial
	default:
		return ""
	}
}

// modelProfileComplete reports whether the essential profile fields are
// populated: pricing, context window, and at least one capability.
func modelProfileComplete(model Model) bool {
	return model.InputPriceUSDPer1M > 0 && model.ContextWindow > 0 && len(model.Capabilities) > 0
}

// applyDerivedModelProfile fills any missing profile fields with derived
// values. It is applied on every model create/update so that imported models
// without a catalog entry still get a usable, searchable profile.
func applyDerivedModelProfile(model *Model, now time.Time) {
	if normalizeProfileTier(model.Tier) == "" {
		model.Tier = deriveModelTier(*model)
	}
	model.Capabilities = deriveModelCapabilities(*model)
	if model.ProfileStatus == "" {
		if modelProfileComplete(*model) {
			model.ProfileStatus = ProfileStatusComplete
		} else {
			model.ProfileStatus = ProfileStatusPartial
		}
	}
	if model.ProfileUpdated.IsZero() {
		model.ProfileUpdated = now
	}
}

// markModelProfileEdited records an admin edit on the profile fields.
func markModelProfileEdited(model *Model, now time.Time) {
	model.ProfileUpdated = now
	if modelProfileComplete(*model) {
		model.ProfileStatus = ProfileStatusComplete
	}
}
