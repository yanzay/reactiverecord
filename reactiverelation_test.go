package reactiverecord

import (
  // "fmt"
  "testing"
)

// type TestObject struct {
//   Some    int
//   Another string
// }

func TestMakeQuery(t *testing.T) {
  obj := TestObject{}
  rec, _ := Connect("sqlite", "test")
  rel := new(ReactiveRelation)
  rel.TableName = rec.tableName(obj)
  rel.Fields = rec.attributeNamesString(obj)
  rel.Where("1 = 1")
  query := rel.MakeQuery()
  good := "SELECT Some, Another FROM test_object WHERE 1 = 1;"
  if query != good {
    t.Errorf("Incorrect query, expected %s, got %s", good, query)
  }
}

func TestFirstRelation(t *testing.T) {
  st, _ := Connect("sqlite", "test.db")
  to := TestObject{Some: 1, Another: "helo"}
  st.CreateTable(to)
  st.TruncateTable(to)
  st.Create(to)
  var result []TestObject
  st.Where(TestObject{}, "Another = 'helo'").First(&result)
  if result[0].Some != 1 || result[0].Another != "helo" {
    t.Errorf("Unexpected first element %v", result)
  }

}
