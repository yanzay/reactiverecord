package reactiverecord

import (
  "fmt"
  // "github.com/eaigner/jet"
  "strings"
)

type ReactiveRelation struct {
  TableName       string
  Fields          string
  WhereConditions []string
  limit           uint
  Rec             *ReactiveRecord
  EmptyObject     interface{}
}

func (rel *ReactiveRelation) Where(query string) *ReactiveRelation {
  rel.WhereConditions = append(rel.WhereConditions, query)
  return rel
}

func (rel *ReactiveRelation) First(obj interface{}) {
  rel.Limit(1)
  rel.Run(obj)
}

func (rel *ReactiveRelation) Limit(num uint) {
  rel.limit = num
}

func (rel *ReactiveRelation) Run(obj interface{}) {
  query := rel.MakeQuery()
  rel.Rec.Query(query).Rows(obj)
}

func (rel *ReactiveRelation) MakeQuery() string {
  query := fmt.Sprintf("SELECT %s FROM %s", rel.Fields, rel.TableName)
  if len(rel.WhereConditions) > 0 {
    where := strings.Join(rel.WhereConditions, " AND ")
    query = fmt.Sprintf("%s WHERE %s", query, where)
  }
  if rel.limit != 0 {
    query = fmt.Sprintf("%s LIMIT %d", query, rel.limit)
  }
  return fmt.Sprintf("%s;", query)
}
