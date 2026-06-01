package configservice

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// computeTextDiff 逐行对比 before/after，生成 git-style diff 行。
func computeTextDiff(path, before, after string) FileDiff {
	oldLines := splitLines(before)
	newLines := splitLines(after)

	action := "modify"
	if before == "" && after != "" {
		action = "create"
	} else if before != "" && after == "" {
		action = "delete"
	} else if before == after {
		action = "modify"
	}

	lines := lcsDiff(oldLines, newLines)
	return FileDiff{Path: path, Action: action, Lines: lines}
}

// computeJSONDiff 将两个 map 格式化为 JSON 后逐行对比。
func computeJSONDiff(path string, before, after map[string]any) FileDiff {
	oldContent := ""
	if before != nil {
		oldContent = formatJSON(before)
	}
	newContent := ""
	if after != nil {
		newContent = formatJSON(after)
	}
	return computeTextDiff(path, oldContent, newContent)
}

// formatJSON 将 map 格式化为缩进 JSON 字符串。
func formatJSON(data map[string]any) string {
	if data == nil {
		return ""
	}
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", data)
	}
	return string(content) + "\n"
}

// maskSensitiveValue 对单个值进行脱敏。
// 保留头尾字符，中间用 *** 遮蔽；短值也保留首尾，避免不同短密钥都显示为 ***。
func maskSensitiveValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	runes := []rune(value)
	length := len(runes)
	if length == 1 {
		return "***"
	}
	if length <= 5 {
		return string(runes[:1]) + "***" + string(runes[length-1:])
	}
	if length <= 10 {
		return string(runes[:3]) + "***" + string(runes[length-2:])
	}
	return string(runes[:8]) + "***" + string(runes[length-5:])
}

// maskTextSensitiveFields 只对指定 TOML 字段的字符串值进行脱敏。
func maskTextSensitiveFields(content string, fields ...string) string {
	result := content
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		re := regexp.MustCompile(`(?m)^(\s*` + regexp.QuoteMeta(field) + `\s*=\s*)("([^"]*)")(\s*(?:#.*)?$)`)
		result = re.ReplaceAllStringFunc(result, func(line string) string {
			match := re.FindStringSubmatch(line)
			if len(match) < 5 {
				return line
			}
			return match[1] + strconv.Quote(maskSensitiveValue(match[3])) + match[4]
		})
	}
	return result
}

// sensitiveFieldKeys 需要脱敏的配置字段名。
var sensitiveFieldKeys = []string{
	"ANTHROPIC_API_KEY",
	"ANTHROPIC_AUTH_TOKEN",
	"OPENAI_API_KEY",
}

// maskMapSensitiveKeys 对 map 中指定 key 的值进行脱敏（返回新 map，不修改原 map）。
func maskMapSensitiveKeys(data map[string]any, keys ...string) map[string]any {
	if data == nil {
		return nil
	}
	result := make(map[string]any, len(data))
	for k, v := range data {
		result[k] = v
	}
	for _, key := range keys {
		if val, ok := result[key]; ok {
			if s, ok := val.(string); ok && s != "" {
				result[key] = maskSensitiveValue(s)
			}
		}
	}
	return result
}

// maskJSONSensitiveKeys 对 JSON map 中嵌套的 env map 内的敏感字段进行脱敏。
func maskJSONSensitiveKeys(data map[string]any) map[string]any {
	if data == nil {
		return nil
	}
	result := make(map[string]any, len(data))
	for k, v := range data {
		result[k] = v
	}
	if env, ok := result["env"].(map[string]any); ok {
		result["env"] = maskMapSensitiveKeys(env, sensitiveFieldKeys...)
	}
	return result
}

// maskTextSensitiveValues 对指定字段的字符串值进行行内脱敏。
// keyValues 只使用 key 作为字段名，value 不参与匹配，避免敏感短值误伤普通文本。
func maskTextSensitiveValues(content string, keyValues map[string]string) string {
	if len(keyValues) == 0 {
		return content
	}
	fields := make([]string, 0, len(keyValues))
	for field := range keyValues {
		fields = append(fields, field)
	}
	return maskTextSensitiveFields(content, fields...)
}

// computeJSONDiffWithMask 用原始数据判定变更类型（removed/added），
// 用字段级掩码后的数据做展示，避免两个不同密钥掩码后相同时漏报变更。
func computeJSONDiffWithMask(path string, oldData, newData map[string]any, keys ...string) FileDiff {
	oldRaw := formatJSON(oldData)
	newRaw := formatJSON(newData)
	if len(keys) == 0 {
		return computeTextDiff(path, oldRaw, newRaw)
	}
	// 字段级掩码只替换敏感字段的值、不改变 map 结构，
	// 因此掩码后与原始 JSON 行数、字段顺序一致，可逐行对齐。
	// 这样既避免文本级 ReplaceAll 对短值（如 "key"）的子串误伤（如 auth_mode="apikey"），
	// 又保留用原始内容判定变更、用掩码内容展示的双轨机制。
	oldMasked := formatJSON(maskDataSensitiveKeys(oldData, keys))
	newMasked := formatJSON(maskDataSensitiveKeys(newData, keys))
	return computeTextDiffFromMasked(path, oldRaw, newRaw, oldMasked, newMasked)
}

// maskDataSensitiveKeys 递归深拷贝 map，并对任意层级中命中 keys 的字符串值脱敏。
// 不修改原始 map。
func maskDataSensitiveKeys(data map[string]any, keys []string) map[string]any {
	if data == nil {
		return nil
	}
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}
	result := make(map[string]any, len(data))
	for k, v := range data {
		switch val := v.(type) {
		case string:
			if keySet[k] && val != "" {
				result[k] = maskSensitiveValue(val)
			} else {
				result[k] = val
			}
		case map[string]any:
			result[k] = maskDataSensitiveKeys(val, keys)
		default:
			result[k] = v
		}
	}
	return result
}

// computeTextDiffWithMask 在原始文本上计算 diff，再按字段对展示内容进行脱敏。
func computeTextDiffWithMask(path, before, after string, fields map[string]string) FileDiff {
	return computeTextDiffWithSensitiveFields(path, before, after, mapKeys(fields)...)
}

// computeTextDiffWithSensitiveFields 在原始文本上计算 diff，再只对指定字段值进行脱敏。
func computeTextDiffWithSensitiveFields(path, before, after string, fields ...string) FileDiff {
	oldMasked := maskTextSensitiveFields(before, fields...)
	newMasked := maskTextSensitiveFields(after, fields...)
	return computeTextDiffFromMasked(path, before, after, oldMasked, newMasked)
}

func mapKeys(values map[string]string) []string {
	if len(values) == 0 {
		return nil
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

// computeTextDiffFromMasked 在原始内容上做 LCS diff 确定变更类型，
// 再将掩码后的内容填充到对应 diff 行。
func computeTextDiffFromMasked(path, oldRaw, newRaw, oldMasked, newMasked string) FileDiff {
	oldRawLines := splitLines(oldRaw)
	newRawLines := splitLines(newRaw)
	oldMaskedLines := splitLines(oldMasked)
	newMaskedLines := splitLines(newMasked)

	action := "modify"
	if oldRaw == "" && newRaw != "" {
		action = "create"
	} else if oldRaw != "" && newRaw == "" {
		action = "delete"
	}

	rawDiff := lcsDiff(oldRawLines, newRawLines)

	// 将掩码后的内容映射到 diff 行
	oldIdx, newIdx := 0, 0
	lines := make([]DiffLine, len(rawDiff))
	for i, d := range rawDiff {
		switch d.Type {
		case "context":
			lines[i] = DiffLine{Type: "context", Content: oldMaskedLines[oldIdx]}
			oldIdx++
			newIdx++
		case "removed":
			lines[i] = DiffLine{Type: "removed", Content: oldMaskedLines[oldIdx]}
			oldIdx++
		case "added":
			lines[i] = DiffLine{Type: "added", Content: newMaskedLines[newIdx]}
			newIdx++
		}
	}

	return FileDiff{Path: path, Action: action, Lines: lines}
}

// extractNestedStringValues 从 map 中提取指定 key 的字符串值，
// 支持嵌套 map（如 env 子 map）的递归查找。
func extractNestedStringValues(data map[string]any, keys []string) map[string]string {
	result := make(map[string]string, len(keys))
	if data == nil {
		return result
	}
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}
	// 顶层直接匹配
	for _, key := range keys {
		if val, ok := data[key]; ok {
			if s, ok := val.(string); ok {
				result[key] = s
			}
		}
	}
	// 递归查找嵌套 map（如 env 子 map）
	for _, v := range data {
		if sub, ok := v.(map[string]any); ok {
			for key := range keySet {
				if _, exists := result[key]; exists {
					continue
				}
				if val, ok := sub[key]; ok {
					if s, ok := val.(string); ok {
						result[key] = s
					}
				}
			}
		}
	}
	return result
}

// splitLines 将文本按换行符分割为行切片。空文本返回空切片。
func splitLines(text string) []string {
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	// 去除末尾空行（由末尾换行符产生）
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// lcsDiff 使用 LCS 算法生成 diff 行序列。
func lcsDiff(oldLines, newLines []string) []DiffLine {
	m, n := len(oldLines), len(newLines)
	if m == 0 && n == 0 {
		return nil
	}

	// 特殊情况优化
	if m == 0 {
		lines := make([]DiffLine, n)
		for i, l := range newLines {
			lines[i] = DiffLine{Type: "added", Content: l}
		}
		return lines
	}
	if n == 0 {
		lines := make([]DiffLine, m)
		for i, l := range oldLines {
			lines[i] = DiffLine{Type: "removed", Content: l}
		}
		return lines
	}

	// LCS DP 表
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if oldLines[i-1] == newLines[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	// 回溯生成 diff（逆序收集，最后反转）
	var reversed []DiffLine
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && oldLines[i-1] == newLines[j-1] {
			reversed = append(reversed, DiffLine{Type: "context", Content: oldLines[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			reversed = append(reversed, DiffLine{Type: "added", Content: newLines[j-1]})
			j--
		} else {
			reversed = append(reversed, DiffLine{Type: "removed", Content: oldLines[i-1]})
			i--
		}
	}

	// 反转得到正确顺序
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}
	return reversed
}
