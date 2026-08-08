package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func parseUsersTableQuery(r *http.Request) (usersTableQuery, error) {
	values := r.URL.Query()
	query := usersTableQuery{
		Start:       0,
		Size:        25,
		Filters:     []usersTableFilter{},
		FilterModes: map[string]string{},
		Sorting:     []usersTableSorting{},
	}

	if raw := strings.TrimSpace(values.Get("start")); raw != "" {
		start, err := strconv.Atoi(raw)
		if err != nil {
			return query, err
		}
		if start > 0 {
			query.Start = start
		}
	}
	if raw := strings.TrimSpace(values.Get("size")); raw != "" {
		size, err := strconv.Atoi(raw)
		if err != nil {
			return query, err
		}
		if size > 0 {
			query.Size = size
		}
		if query.Size > 1000 {
			query.Size = 1000
		}
	}
	if raw := strings.TrimSpace(values.Get("filters")); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.Filters); err != nil {
			return query, err
		}
	}
	if raw := strings.TrimSpace(values.Get("filterModes")); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.FilterModes); err != nil {
			return query, err
		}
	}
	if raw := strings.TrimSpace(values.Get("sorting")); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.Sorting); err != nil {
			return query, err
		}
	}

	return query, nil
}

func filterUsersTableResponse(users []userAPI, filters []usersTableFilter, modes map[string]string) []userAPI {
	if len(filters) == 0 {
		return users
	}

	filtered := make([]userAPI, 0, len(users))
	for _, user := range users {
		matches := true
		for _, filter := range filters {
			if !matchesUsersTableFilter(user, filter, modes[filter.ID]) {
				matches = false
				break
			}
		}
		if matches {
			filtered = append(filtered, user)
		}
	}

	return filtered
}

func matchesUsersTableFilter(user userAPI, filter usersTableFilter, mode string) bool {
	if filter.Value == nil {
		return true
	}
	values := normalizeUsersTableFilterValues(filter.Value)
	if len(values) == 0 {
		return true
	}

	if filter.ID == "activeInternalSquads" {
		for _, value := range values {
			for _, squad := range user.ActiveInternalSquads {
				if strings.EqualFold(squad.UUID, value) || strings.EqualFold(squad.Name, value) {
					return true
				}
			}
		}
		return false
	}

	field := usersTableFieldValue(user, filter.ID)
	if field == nil {
		return false
	}

	if isNumericUsersTableFilterMode(mode) {
		return matchesUsersTableNumericFilter(field, values, mode)
	}

	fieldText := strings.ToLower(strings.TrimSpace(fmt.Sprint(field)))
	if fieldText == "" {
		return false
	}
	if mode == "" {
		mode = "contains"
	}

	for _, value := range values {
		needle := strings.ToLower(strings.TrimSpace(value))
		if needle == "" {
			continue
		}
		switch mode {
		case "equals", "equalsString":
			if fieldText == needle {
				return true
			}
		case "startsWith":
			if strings.HasPrefix(fieldText, needle) {
				return true
			}
		case "endsWith":
			if strings.HasSuffix(fieldText, needle) {
				return true
			}
		default:
			if strings.Contains(fieldText, needle) {
				return true
			}
		}
	}

	return false
}

func normalizeUsersTableFilterValues(value any) []string {
	switch typed := value.(type) {
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			if item == nil {
				continue
			}
			values = append(values, strings.TrimSpace(fmt.Sprint(item)))
		}
		return values
	case []string:
		return typed
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil
		}
		return []string{typed}
	case float64:
		return []string{strconv.FormatFloat(typed, 'f', -1, 64)}
	case int:
		return []string{strconv.Itoa(typed)}
	case int64:
		return []string{strconv.FormatInt(typed, 10)}
	default:
		text := strings.TrimSpace(fmt.Sprint(typed))
		if text == "" || text == "<nil>" {
			return nil
		}
		return []string{text}
	}
}

func usersTableFieldValue(user userAPI, columnID string) any {
	switch columnID {
	case "uuid":
		return user.ID
	case "id":
		return user.ID
	case "shortUuid":
		return user.ShortUUID
	case "username":
		return user.Username
	case "status":
		return user.Status
	case "trafficLimitBytes":
		return user.TrafficLimitBytes
	case "trafficLimitStrategy":
		return user.TrafficLimitStrategy
	case "expireAt":
		return user.ExpireAt
	case "telegramId":
		if user.TelegramID == nil {
			return nil
		}
		return *user.TelegramID
	case "email":
		if user.Email == nil {
			return nil
		}
		return *user.Email
	case "description":
		if user.Description == nil {
			return nil
		}
		return *user.Description
	case "tag":
		if user.Tag == nil {
			return nil
		}
		return *user.Tag
	case "hwidDeviceLimit":
		if user.HwidDeviceLimit == nil {
			return nil
		}
		return *user.HwidDeviceLimit
	case "externalSquadUuid":
		if user.ExternalSquadUUID == nil {
			return nil
		}
		return *user.ExternalSquadUUID
	case "trojanPassword":
		return user.TrojanPassword
	case "vlessUuid":
		return user.VlessUUID
	case "ssPassword":
		return user.SSPassword
	case "naivePassword":
		return user.NaivePassword
	case "shadowtlsPassword":
		return user.ShadowtlsPassword
	case "hysteria2Password":
		return user.Hysteria2Password
	case "anytlsPassword":
		return user.AnytlsPassword
	case "createdAt":
		return user.CreatedAt
	case "updatedAt":
		return user.UpdatedAt
	case "subRevokedAt":
		return user.SubRevokedAt
	case "lastTrafficResetAt":
		return user.LastTrafficResetAt
	case "nodeName", "userTraffic.lastConnectedNodeUuid":
		if user.UserTraffic.LastConnectedNodeUUID == nil {
			return nil
		}
		return *user.UserTraffic.LastConnectedNodeUUID
	case "usedTrafficBytes", "userTraffic.usedTrafficBytes":
		return user.UserTraffic.UsedTrafficBytes
	case "userTraffic.lifetimeUsedTrafficBytes":
		return user.UserTraffic.LifetimeUsedTrafficBytes
	case "userTraffic.onlineAt":
		return user.UserTraffic.OnlineAt
	case "userTraffic.firstConnectedAt":
		return user.UserTraffic.FirstConnectedAt
	default:
		return nil
	}
}

func isNumericUsersTableFilterMode(mode string) bool {
	switch mode {
	case "equals", "greaterThan", "greaterThanOrEqualTo", "lessThan", "lessThanOrEqualTo", "between", "betweenInclusive":
		return true
	default:
		return false
	}
}

func matchesUsersTableNumericFilter(field any, values []string, mode string) bool {
	fieldValue, ok := usersTableFloatValue(field)
	if !ok {
		return false
	}
	parse := func(index int) (float64, bool) {
		if index >= len(values) || strings.TrimSpace(values[index]) == "" {
			return 0, false
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(values[index]), 64)
		return value, err == nil
	}

	switch mode {
	case "greaterThan":
		value, ok := parse(0)
		return ok && fieldValue > value
	case "greaterThanOrEqualTo":
		value, ok := parse(0)
		return ok && fieldValue >= value
	case "lessThan":
		value, ok := parse(0)
		return ok && fieldValue < value
	case "lessThanOrEqualTo":
		value, ok := parse(0)
		return ok && fieldValue <= value
	case "between", "betweenInclusive":
		min, minOK := parse(0)
		max, maxOK := parse(1)
		if minOK && maxOK {
			return fieldValue >= min && fieldValue <= max
		}
		if minOK {
			return fieldValue >= min
		}
		if maxOK {
			return fieldValue <= max
		}
		return true
	default:
		value, ok := parse(0)
		return ok && fieldValue == value
	}
}

func usersTableFloatValue(value any) (float64, bool) {
	switch typed := value.(type) {
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case float64:
		return typed, true
	case *int:
		if typed == nil {
			return 0, false
		}
		return float64(*typed), true
	case *int64:
		if typed == nil {
			return 0, false
		}
		return float64(*typed), true
	default:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(value)), 64)
		return parsed, err == nil
	}
}

func sortUsersTableResponse(users []userAPI, sorting []usersTableSorting) {
	if len(sorting) == 0 {
		return
	}

	sort.SliceStable(users, func(i, j int) bool {
		for _, sortRule := range sorting {
			comparison := compareUsersTableValues(usersTableFieldValue(users[i], sortRule.ID), usersTableFieldValue(users[j], sortRule.ID))
			if comparison == 0 {
				continue
			}
			if sortRule.Desc {
				return comparison > 0
			}
			return comparison < 0
		}
		return false
	})
}

func compareUsersTableValues(left any, right any) int {
	leftTime, leftIsTime := usersTableTimeValue(left)
	rightTime, rightIsTime := usersTableTimeValue(right)
	if leftIsTime || rightIsTime {
		if !leftIsTime && rightIsTime {
			return -1
		}
		if leftIsTime && !rightIsTime {
			return 1
		}
		if leftTime.Before(rightTime) {
			return -1
		}
		if leftTime.After(rightTime) {
			return 1
		}
		return 0
	}

	leftFloat, leftIsFloat := usersTableFloatValue(left)
	rightFloat, rightIsFloat := usersTableFloatValue(right)
	if leftIsFloat || rightIsFloat {
		if !leftIsFloat && rightIsFloat {
			return -1
		}
		if leftIsFloat && !rightIsFloat {
			return 1
		}
		if leftFloat < rightFloat {
			return -1
		}
		if leftFloat > rightFloat {
			return 1
		}
		return 0
	}

	leftText := strings.ToLower(strings.TrimSpace(fmt.Sprint(left)))
	rightText := strings.ToLower(strings.TrimSpace(fmt.Sprint(right)))
	if leftText < rightText {
		return -1
	}
	if leftText > rightText {
		return 1
	}
	return 0
}

func usersTableTimeValue(value any) (time.Time, bool) {
	switch typed := value.(type) {
	case time.Time:
		return typed, true
	case *time.Time:
		if typed == nil {
			return time.Time{}, false
		}
		return *typed, true
	default:
		return time.Time{}, false
	}
}

func paginateUsersTableResponse(users []userAPI, start int, size int) []userAPI {
	if start < 0 {
		start = 0
	}
	if size <= 0 {
		size = 25
	}
	if start >= len(users) {
		return []userAPI{}
	}
	end := start + size
	if end > len(users) {
		end = len(users)
	}
	return users[start:end]
}
