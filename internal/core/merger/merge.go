package merger

import "encoding/json"

func (m *Merger) Merge(configs map[string]map[string]any) map[string]any {
	result := make(map[string]any)

	types := []string{"global", "app", "env"}

	for _, t := range types {
		if config, ok := configs[t]; ok {
			result = m.deepMerge(result, config)
		}
	}

	return result
}

func (m *Merger) deepMerge(dst, src map[string]any) map[string]any {
	if dst == nil {
		dst = make(map[string]any)
	}

	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			if srcMap, ok := srcVal.(map[string]any); ok {
				if dstMap, ok := dstVal.(map[string]any); ok {
					dst[key] = m.deepMerge(dstMap, srcMap)
					continue
				}
			}
		}
		dst[key] = srcVal
	}

	return dst
}

func (m *Merger) MergeJSON(configs map[string]json.RawMessage) (json.RawMessage, error) {
	parsed := make(map[string]map[string]any)

	for k, v := range configs {
		var data map[string]any
		if err := json.Unmarshal(v, &data); err != nil {
			return nil, err
		}
		parsed[k] = data
	}

	merged := m.Merge(parsed)

	result, err := json.Marshal(merged)
	return result, err
}
