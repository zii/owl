package db

import (
	"fmt"
	"server/base"
	"strings"
)

type condition struct {
	args  []interface{}
	conds []string // WHERE clause conditions
}

// Add new argument(s).
func (p *condition) With(args ...interface{}) {
	p.args = append(p.args, args...)
}

func (p *condition) Args() []interface{} {
	return p.args
}

// insert backslash before double quotes
func escape(s string) string {
	return strings.ReplaceAll(s, "\"", `\"`)
}

// quote with double quotes
func quote(s string) string {
	return "\"" + s + "\""
}

// insert backslash before double quotes, and quote with double quotes.
func literal(s string) string {
	return quote(escape(s))
}

func joinStrings(a []string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return literal(a[0])
	}
	var out = literal(a[0])
	for _, s := range a[1:] {
		out += "," + literal(s)
	}
	return out
}

func join(a interface{}) string {
	switch v := a.(type) {
	case []int:
		return base.JoinInts(v)
	case []int64:
		return base.JoinInt64s(v)
	case []string:
		return joinStrings(v)
	}
	panic(fmt.Sprintf("invalid arg type: %T", a))
}

// Add a WHERE condition with argument(s).
// e.g.
// p.And("a=?", 1).And("b!=2").And("c like ?", name+"%%")
func (p *condition) And(condition string, args ...interface{}) *condition {
	if condition != "" {
		p.conds = append(p.conds, condition)
	}
	p.args = append(p.args, args...)
	return p
}

// IN operator
// e.g.
// p.And(db.In("id", []int{1,2,3}))
func (p *condition) In(column string, a interface{}) *condition {
	s := fmt.Sprintf("%s in (%s)", column, join(a))
	return p.And(s)
}

// NOT IN operator
func (p *condition) NotIn(column string, a interface{}) *condition {
	c := fmt.Sprintf("%s not in (%s)", column, join(a))
	return p.And(c)
}

// A simple queryset to contruct SQL prepared statement.
type prepared struct {
	condition
	orderby string
	offset  int
	limit   int
}

// create a queryset object
func Prepare() *prepared {
	p := &prepared{
		limit: -1,
	}
	return p
}

// WHERE clause string
func (p *prepared) Where() string {
	if len(p.conds) == 0 {
		return ""
	}
	s := strings.Join(p.conds, " and ")
	return " where " + s
}

// Set ORDER BY keywords
func (p *prepared) Sort(keywords string) {
	p.orderby = keywords
}

// ORDER BY clause string
func (p *prepared) OrderBy() string {
	if p.orderby == "" {
		return ""
	}
	return " order by " + p.orderby
}

// Set offset and limit
func (p *prepared) Slice(offset, limit int) {
	p.offset = offset
	p.limit = limit
}

// LIMIT clause string
func (p *prepared) Limit() string {
	if p.limit < 0 {
		return ""
	}
	if p.offset == 0 {
		return fmt.Sprintf(" limit %d", p.limit)
	}
	return fmt.Sprintf(" limit %d, %d", p.offset, p.limit)
}

// clause compilation
func (p *prepared) Clause() string {
	return p.Where() + p.OrderBy() + p.Limit()
}
