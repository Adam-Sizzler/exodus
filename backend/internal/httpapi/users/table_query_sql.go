package users

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// usersFilterColumnMap is a 1:1 port of Remnawave's USERS_FILTER_COLUMN_MAP
// (users.repository.ts). Any filter/sort id NOT in this map is silently
// ignored — exactly Remnawave's `if (!(filter.id in USERS_FILTER_COLUMN_MAP))
// continue;` — never treated as "exclude everything".
//
// naivePassword/shadowtlsPassword/hysteria2Password/anytlsPassword/
// trafficLimitStrategy/updatedAt have no Remnawave equivalent (sing-box-only
// protocols and an Exodus-only field). They're added here and go through the
// same generic fallback switch Remnawave applies to any of its own
// non-special-cased columns, so they behave the same way a Remnawave column
// with no special case would.
var usersFilterColumnMap = map[string]string{
	"id":                  "u.id",
	"uuid":                "u.uuid",
	"createdAt":           "u.created_at",
	"created_at":          "u.created_at",
	"expireAt":            "u.expire_at",
	"expire_at":           "u.expire_at",
	"lastTrafficResetAt":  "u.last_traffic_reset_at",
	"last_traffic_reset_at": "u.last_traffic_reset_at",
	"subRevokedAt":        "u.sub_revoked_at",
	"sub_revoked_at":      "u.sub_revoked_at",
	"telegramId":          "u.telegram_id",
	"telegram_id":         "u.telegram_id",
	"vlessUuid":           "u.vless_uuid",
	"vless_uuid":          "u.vless_uuid",
	"trojanPassword":      "u.trojan_password",
	"trojan_password":     "u.trojan_password",
	"externalSquadUuid":   "u.external_squad_uuid",
	"external_squad_uuid": "u.external_squad_uuid",
	"username":            "u.username",
	"status":              "u.status",
	"shortUuid":           "u.short_uuid",
	"short_uuid":          "u.short_uuid",
	"description":         "u.description",
	"email":               "u.email",
	"tag":                 "u.tag",
	"hwidDeviceLimit":     "u.hwid_device_limit",
	"hwid_device_limit":   "u.hwid_device_limit",
	"trafficLimitBytes":   "u.traffic_limit_bytes",
	"traffic_limit_bytes": "u.traffic_limit_bytes",
	"usedTrafficBytes":    "ut.used_traffic_bytes",
	"used_traffic_bytes":  "ut.used_traffic_bytes",

	"userTraffic.usedTrafficBytes":         "ut.used_traffic_bytes",
	"userTraffic.onlineAt":                 "ut.online_at",
	"userTraffic.firstConnectedAt":         "ut.first_connected_at",
	"userTraffic.lifetimeUsedTrafficBytes": "ut.lifetime_used_traffic_bytes",
	"userTraffic.lastConnectedNodeUuid":    "ut.last_connected_node_uuid",
	"onlineAt":                             "ut.online_at",
	"online_at":                            "ut.online_at",
	"firstConnectedAt":                     "ut.first_connected_at",
	"first_connected_at":                   "ut.first_connected_at",
	"lifetimeUsedTrafficBytes":             "ut.lifetime_used_traffic_bytes",
	"lifetime_used_traffic_bytes":          "ut.lifetime_used_traffic_bytes",
	"lastConnectedNodeUuid":                "ut.last_connected_node_uuid",
	"last_connected_node_uuid":             "ut.last_connected_node_uuid",

	// Exodus-only protocols & fields
	"trafficLimitStrategy":   "u.traffic_limit_strategy",
	"traffic_limit_strategy": "u.traffic_limit_strategy",
	"updatedAt":              "u.updated_at",
	"updated_at":             "u.updated_at",
	"ssPassword":             "u.ss_password",
	"ss_password":            "u.ss_password",
	"naivePassword":          "u.naive_password",
	"naive_password":         "u.naive_password",
	"shadowtlsPassword":      "u.shadowtls_password",
	"shadowtls_password":     "u.shadowtls_password",
	"hysteria2Password":      "u.hysteria2_password",
	"hysteria2_password":     "u.hysteria2_password",
	"anytlsPassword":         "u.anytls_password",
	"anytls_password":        "u.anytls_password",
}

// usersNumericFilterIDs is a 1:1 port of Remnawave's NUMERIC_FILTER_IDS.
var usersNumericFilterIDs = map[string]struct{}{
	"hwidDeviceLimit":   {},
	"id":                {},
	"trafficLimitBytes": {},
}

// usersExactDateFilterIDs is a 1:1 port of the id list Remnawave special-cases
// to an exact-equality date match. Note subRevokedAt is deliberately excluded
// here too — Remnawave doesn't special-case it either, so it falls through to
// the generic switch below, same as upstream.
var usersExactDateFilterIDs = map[string]struct{}{
	"createdAt":            {},
	"expireAt":             {},
	"lastTrafficResetAt":   {},
	"userTraffic.onlineAt": {},
}

// buildUsersTableQuery is a 1:1 port of Remnawave's getAllUsers + applyUsersFilters:
// filtering, sorting and pagination pushed down to Postgres instead of loading
// every user row and filtering/sorting it in Go.
func buildUsersTableQuery(filters []usersTableFilter, filterModes map[string]string, sorting []usersTableSorting) (whereSQL string, orderSQL string, args []any, err error) {
	clauses := make([]string, 0, len(filters))

	for _, filter := range filters {
		_, known := usersFilterColumnMap[filter.ID]
		if !known && filter.ID != "activeInternalSquads" && filter.ID != "nodeName" {
			continue // unknown id: ignored, matching Remnawave's `continue`
		}
		if filter.Value == nil {
			continue
		}
		if arr, isArr := filter.Value.([]any); isArr && len(arr) == 0 {
			continue
		}

		mode := filterModes[filter.ID]
		if mode == "" {
			mode = "contains"
		}

		switch {
		case filter.ID == "activeInternalSquads":
			// Remnawave: single squad uuid, exact match via subquery — not a
			// multi-value OR, and not a name match (the picker sends one uuid).
			v, ok := singleStringValue(filter.Value)
			if !ok {
				continue
			}
			args = append(args, v)
			clauses = append(clauses, fmt.Sprintf(
				`u.id IN (SELECT user_id FROM internal_squad_members WHERE internal_squad_uuid = $%d::uuid)`,
				len(args),
			))

		case filter.ID == "nodeName":
			// Remnawave: exact match against the node uuid (the picker sends one uuid).
			v, ok := singleStringValue(filter.Value)
			if !ok {
				continue
			}
			args = append(args, v)
			clauses = append(clauses, fmt.Sprintf(`ut.last_connected_node_uuid = $%d::uuid`, len(args)))

		case filter.ID == "externalSquadUuid":
			v, ok := singleStringValue(filter.Value)
			if !ok {
				continue
			}
			args = append(args, v)
			clauses = append(clauses, fmt.Sprintf(`u.external_squad_uuid = $%d::uuid`, len(args)))

		case filter.ID == "vlessUuid":
			v, ok := singleStringValue(filter.Value)
			if !ok {
				continue
			}
			args = append(args, "%"+v+"%")
			clauses = append(clauses, fmt.Sprintf(`u.vless_uuid::text ILIKE $%d`, len(args)))

		case filter.ID == "id":
			v, ok := singleStringValue(filter.Value)
			if !ok {
				continue
			}
			if _, parseErr := strconv.ParseInt(strings.TrimSpace(v), 10, 64); parseErr != nil {
				continue // not numeric-looking: filter dropped, matching Remnawave's swallowed catch
			}
			args = append(args, "%"+v+"%")
			clauses = append(clauses, fmt.Sprintf(`CAST(u.id AS TEXT) LIKE $%d`, len(args)))

		case filter.ID == "telegramId":
			v, ok := singleStringValue(filter.Value)
			if !ok {
				continue
			}
			if _, parseErr := strconv.ParseInt(strings.TrimSpace(v), 10, 64); parseErr != nil {
				clauses = append(clauses, `u.telegram_id IS NULL`)
				continue
			}
			args = append(args, "%"+v+"%")
			clauses = append(clauses, fmt.Sprintf(`CAST(u.telegram_id AS TEXT) LIKE $%d`, len(args)))

		case isExactDateFilter(filter.ID):
			v, ok := singleStringValue(filter.Value)
			if !ok {
				continue
			}
			ts, parseErr := parseFlexibleDate(v)
			if parseErr != nil {
				continue
			}
			args = append(args, ts)
			clauses = append(clauses, fmt.Sprintf(`%s = $%d`, usersFilterColumnMap[filter.ID], len(args)))

		default:
			col := usersFilterColumnMap[filter.ID]
			_, numeric := usersNumericFilterIDs[filter.ID]
			clause, clauseArgs, ok, clauseErr := buildGenericUsersFilterClause(col, filter.Value, mode, numeric, len(args))
			if clauseErr != nil {
				return "", "", nil, clauseErr
			}
			if !ok {
				continue
			}
			clauses = append(clauses, clause)
			args = append(args, clauseArgs...)
		}
	}

	if len(clauses) > 0 {
		whereSQL = "WHERE " + strings.Join(clauses, " AND ")
	}

	orderParts := make([]string, 0, len(sorting))
	for _, sort := range sorting {
		dir := "ASC"
		if sort.Desc {
			dir = "DESC"
		}
		if sort.ID == "usedTrafficPercentage" {
			// Remnawave: usedTrafficBytes / NULLIF(trafficLimitBytes, 0), computed at sort time.
			orderParts = append(orderParts, fmt.Sprintf(
				"CAST(COALESCE(ut.used_traffic_bytes, 0) AS NUMERIC) / NULLIF(u.traffic_limit_bytes, 0) %s NULLS LAST", dir,
			))
			continue
		}
		col, ok := usersFilterColumnMap[sort.ID]
		if !ok {
			// Unknown sort ID: gracefully skip to avoid breaking the users table.
			continue
		}
		orderParts = append(orderParts, fmt.Sprintf("%s %s NULLS LAST", col, dir))
	}
	if len(orderParts) == 0 {
		orderParts = append(orderParts, "u.id DESC")
	}
	orderSQL = "ORDER BY " + strings.Join(orderParts, ", ")

	return whereSQL, orderSQL, args, nil
}

// buildGenericUsersFilterClause is a 1:1 port of the generic switch(mode) branch
// at the end of Remnawave's applyUsersFilters:
//
//	const value = NUMERIC_FILTER_IDS.has(filter.id) ? Number(filter.value) : filter.value;
//	switch (mode) { ... }
//
// `value` is computed once and reused across cases — including startsWith/
// endsWith/contains, which is why a numeric-id field goes through the same
// Number() normalization there too, not just for equals/comparisons. The
// 'equals' + array branch is the one exception: Remnawave binds the RAW
// filter.value array via `'in'` there, not the Number-converted `value`.
//
// A malformed numeric filter value returns a non-nil error instead of
// silently dropping the filter, matching Remnawave's behavior for this case:
// `Number('garbage')` is NaN in JS, gets bound to the query as-is, and
// Postgres rejects it at execution time — failing the whole request. Go's
// strconv.ParseFloat fails at parse time instead of bind time, but the
// caller (buildUsersTableQuery) propagates the error the same way, so the
// end result is the same request-wide failure, not a partially-applied filter.
func buildGenericUsersFilterClause(col string, rawValue any, mode string, numeric bool, argOffset int) (string, []any, bool, error) {
	// case 'equals' with an array: Remnawave uses filter.value as-is, unconverted.
	if mode == "equals" {
		if arr, ok := rawValue.([]any); ok && len(arr) > 0 {
			vals := make([]any, len(arr))
			for i, item := range arr {
				vals[i] = fmt.Sprint(item)
			}
			placeholders := make([]string, len(vals))
			for i := range vals {
				placeholders[i] = fmt.Sprintf("$%d", argOffset+i+1)
			}
			return fmt.Sprintf("%s IN (%s)", col, strings.Join(placeholders, ", ")), vals, true, nil
		}
	}

	// case 'between': doesn't use the shared `value` below at all — it reads
	// filter.value as a [from, to] tuple directly with its own per-bound cast.
	// Handled first so a whole-array Number() parse failure (JS: silently NaN,
	// unused for this mode anyway) can't short-circuit it in Go.
	if mode == "between" {
		castFn := func(v any) (any, error) {
			if !numeric {
				return fmt.Sprint(v), nil
			}
			n, err := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(v)), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid numeric filter value %q for %s: %w", fmt.Sprint(v), col, err)
			}
			return n, nil
		}
		arr, ok := rawValue.([]any)
		if !ok || len(arr) != 2 {
			return "", nil, false, nil
		}
		from, to := arr[0], arr[1]
		parts := make([]string, 0, 2)
		args := make([]any, 0, 2)
		if from != nil && strings.TrimSpace(fmt.Sprint(from)) != "" {
			v, err := castFn(from)
			if err != nil {
				return "", nil, false, err
			}
			args = append(args, v)
			parts = append(parts, fmt.Sprintf("%s >= $%d", col, argOffset+len(args)))
		}
		if to != nil && strings.TrimSpace(fmt.Sprint(to)) != "" {
			v, err := castFn(to)
			if err != nil {
				return "", nil, false, err
			}
			args = append(args, v)
			parts = append(parts, fmt.Sprintf("%s <= $%d", col, argOffset+len(args)))
		}
		if len(parts) == 0 {
			return "", nil, false, nil
		}
		return strings.Join(parts, " AND "), args, true, nil
	}

	// const value = NUMERIC_FILTER_IDS.has(filter.id) ? Number(filter.value) : filter.value;
	value := rawValue
	if numeric {
		n, err := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(rawValue)), 64)
		if err != nil {
			return "", nil, false, fmt.Errorf("invalid numeric filter value %q for %s: %w", fmt.Sprint(rawValue), col, err)
		}
		value = n
	}

	switch mode {
	case "equals":
		return fmt.Sprintf("%s = $%d", col, argOffset+1), []any{value}, true, nil

	case "startsWith":
		return fmt.Sprintf("%s LIKE $%d", col, argOffset+1), []any{fmt.Sprint(value) + "%"}, true, nil

	case "endsWith":
		return fmt.Sprintf("%s LIKE $%d", col, argOffset+1), []any{"%" + fmt.Sprint(value)}, true, nil

	case "greaterThan", "greaterThanOrEqualTo", "lessThan", "lessThanOrEqualTo":
		op := map[string]string{
			"greaterThan": ">", "greaterThanOrEqualTo": ">=",
			"lessThan": "<", "lessThanOrEqualTo": "<=",
		}[mode]
		return fmt.Sprintf("%s %s $%d", col, op, argOffset+1), []any{value}, true, nil

	default: // "contains" and any unrecognized mode
		return fmt.Sprintf("%s ILIKE $%d", col, argOffset+1), []any{"%" + fmt.Sprint(value) + "%"}, true, nil
	}
}

func isExactDateFilter(filterID string) bool {
	_, ok := usersExactDateFilterIDs[filterID]
	return ok
}

// singleStringValue extracts a single string from a filter value that's
// expected to carry exactly one id/uuid (select-style pickers), matching how
// Remnawave's special-cased branches use `filter.value as string` directly.
func singleStringValue(v any) (string, bool) {
	switch t := v.(type) {
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return "", false
		}
		return s, true
	case []any:
		if len(t) != 1 || t[0] == nil {
			return "", false
		}
		s := strings.TrimSpace(fmt.Sprint(t[0]))
		if s == "" {
			return "", false
		}
		return s, true
	case nil:
		return "", false
	default:
		s := strings.TrimSpace(fmt.Sprint(t))
		if s == "" {
			return "", false
		}
		return s, true
	}
}

// parseFlexibleDate approximates JS's `new Date(string)` for the ISO-ish
// timestamps a date filter would actually send.
func parseFlexibleDate(v string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if ts, err := time.Parse(layout, v); err == nil {
			return ts, nil
		}
	}
	return time.Time{}, fmt.Errorf("unparseable date: %s", v)
}
