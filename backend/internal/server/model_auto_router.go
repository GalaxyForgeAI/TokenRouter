package server

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// AutoModelName is the reserved model name that asks the gateway to pick the
// best model for a request. The model is resolved from the request features
// (multimodality, tools, structured output, prompt size) against the model
// profile (tier, capabilities, pricing) of every published model, then routed
// exactly like an explicitly named model — so quotas, audit, and usage
// attribution all use the resolved model.
const AutoModelName = "auto"

// AutoModelError is returned when no published model satisfies the request's
// quality constraints. Callers surface it as a 400 with a diagnostic code so
// the rejection is visible in the request log.
type AutoModelError struct {
	Reason string
}

func (e *AutoModelError) Error() string {
	if e.Reason == "" {
		return "no eligible model for auto routing"
	}
	return e.Reason
}

// autoTierRank orders the capability tiers; higher is more capable.
func autoTierRank(tier string) int {
	switch tier {
	case ModelTierFrontier:
		return 3
	case ModelTierFlagship:
		return 2
	case ModelTierStandard:
		return 1
	case ModelTierBasic:
		return 0
	default:
		return 0
	}
}

type autoRequestFeatures struct {
	hasImage        bool
	hasTools        bool
	wantsStructured bool
	estTokens       int64
}

// extractAutoRequestFeatures extracts routing-relevant features from a chat
// request. Content may be a plain string or a multimodal part array; both are
// handled.
func extractAutoRequestFeatures(req *ChatCompletionRequest) autoRequestFeatures {
	var features autoRequestFeatures
	features.hasTools = false
	if tools, ok := req.Tools.([]any); ok {
		features.hasTools = len(tools) > 0
	}
	if req.ResponseFormat != nil {
		if formatted, ok := req.ResponseFormat.(map[string]any); ok {
			if kind, ok := formatted["type"].(string); ok && kind == "json_object" || kind == "json_schema" {
				features.wantsStructured = true
			}
		}
	}
	for _, message := range req.Messages {
		switch content := message.Content.(type) {
		case string:
			features.estTokens += int64(len(content)/4 + 1)
		case []any:
			for _, part := range content {
				block, ok := part.(map[string]any)
				if !ok {
					continue
				}
				switch block["type"] {
				case "text":
					if text, ok := block["text"].(string); ok {
						features.estTokens += int64(len(text)/4 + 1)
					}
				case "image_url", "input_image":
					features.hasImage = true
				}
			}
		}
	}
	return features
}

// autoMinTierRank derives the minimum capability tier a request may use. The
// tier is the quality anchor: long or complex prompts must land on models with
// enough headroom, while short queries may use the cheapest tier.
func autoMinTierRank(features autoRequestFeatures) int {
	switch {
	case features.estTokens > 80000:
		return autoTierRank(ModelTierFlagship)
	case features.estTokens > 16000:
		return autoTierRank(ModelTierStandard)
	case features.hasImage:
		// Multimodal inference is heavier; prefer a capable model.
		return autoTierRank(ModelTierStandard)
	default:
		return autoTierRank(ModelTierBasic)
	}
}

// autoRequiredCapabilities lists capability tags a request demands. Models
// lacking any required tag are excluded before cost optimization.
func autoRequiredCapabilities(features autoRequestFeatures) []string {
	var required []string
	if features.hasImage {
		required = append(required, "vision")
	}
	if features.hasTools {
		required = append(required, "tools")
	}
	if features.wantsStructured {
		required = append(required, "structured_outputs")
	}
	return required
}

// resolveAutoModel picks the cheapest published model that satisfies the
// request's quality constraints. Candidates are models with at least one
// active route (otherwise the model could not be served), whose tier meets the
// minimum derived from the request features, and whose capabilities cover the
// request's demands. Ties on input price are broken by output price.
func (s *Server) resolveAutoModel(ctx context.Context, req *ChatCompletionRequest) (Model, error) {
	features := extractAutoRequestFeatures(req)
	minTier := autoMinTierRank(features)
	required := autoRequiredCapabilities(features)

	routedModels := map[string]bool{}
	for _, route := range s.store.ListRoutes() {
		if route.Status == StatusActive {
			routedModels[route.ModelName] = true
		}
	}

	var candidates []Model
	for _, model := range s.store.ListModels() {
		if model.Status != StatusActive || !routedModels[model.Name] {
			continue
		}
		if autoTierRank(model.Tier) < minTier {
			continue
		}
		if !modelHasAllCapabilities(model.Capabilities, required) {
			continue
		}
		candidates = append(candidates, model)
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		left, right := candidates[i], candidates[j]
		leftPrice, rightPrice := autoSortablePrice(left), autoSortablePrice(right)
		if leftPrice != rightPrice {
			return leftPrice < rightPrice
		}
		return autoSortablePriceOut(left) < autoSortablePriceOut(right)
	})
	if len(candidates) == 0 {
		return Model{}, &AutoModelError{Reason: fmt.Sprintf(
			"auto routing found no eligible model: min_tier=%d required_capabilities=%v",
			minTier, required,
		)}
	}
	return candidates[0], nil
}

// modelHasAllCapabilities reports whether the model covers every required tag.
func modelHasAllCapabilities(actual []string, required []string) bool {
	for _, tag := range required {
		found := false
		for _, capability := range actual {
			if strings.EqualFold(capability, tag) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// autoSortablePrice treats missing or zero pricing as infinitely expensive so
// un-priced models sort last instead of being free.
func autoSortablePrice(model Model) float64 {
	if model.InputPriceUSDPer1M <= 0 {
		return 1e18
	}
	return model.InputPriceUSDPer1M
}

func autoSortablePriceOut(model Model) float64 {
	if model.OutputPriceUSDPer1M <= 0 {
		return 1e18
	}
	return model.OutputPriceUSDPer1M
}
