package user

import "strings"

func appendUserExistsFilter(query string, value, table string) string {
	if value == "" || value == "all" {
		return query
	}
	exists := "EXISTS (SELECT 1 FROM " + table + " WHERE " + table + ".user_id = users.id)"
	if value == "true" {
		return query + " AND " + exists
	}
	if value == "false" {
		return query + " AND NOT " + exists
	}
	return query
}

func appendUserSearchFilter(query string, args []any, search string) (string, []any) {
	trimmed := strings.TrimSpace(search)
	if trimmed == "" {
		return query, args
	}
	query += " AND (email LIKE ? OR name LIKE ?)"
	like := "%" + trimmed + "%"
	return query, append(args, like, like)
}
