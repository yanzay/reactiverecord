package reactiverecord

import (
  "fmt"
  "testing"
)

type TestObject struct {
  Some    int
  Another string
}

func TestConnect(t *testing.T) {
  rec, err := Connect("sqlite", "test.db")
  if err != nil {
    t.Errorf("Error connecting to database %v\n", err)
  }
  fmt.Printf("ReactiveRecord: %v\n", rec)
}

func TestCreate(t *testing.T) {
  st, err := Connect("sqlite", "test.db")
  to := TestObject{Some: 1, Another: "helo"}
  st.CreateTable(to)
  err = st.Create(to)
  if err != nil {
    t.Errorf("Got an error %v", err)
  }
}

func TestAll(t *testing.T) {
  st, _ := Connect("sqlite", "test.db")
  to := TestObject{Some: 1, Another: "helo"}
  st.CreateTable(to)
  st.Create(to)
  var rows []TestObject
  st.All(TestObject{}, &rows)
  first := rows[0]
  if first.Some != 1 || first.Another != "helo" {
    t.Errorf("Expected {1 helo}, got %v", first)
  }
}

func TestWhere(t *testing.T) {
  st, _ := Connect("sqlite", "test.db")
  to := TestObject{Some: 1, Another: "helo"}
  st.CreateTable(to)
  st.TruncateTable(to)
  st.Create(to)
  relation := st.Where(TestObject{}, "some = 1 AND another = 'helo'")
  if relation == nil {
    t.Error("Relation should not be nil.")
  }
  var result []TestObject
  relation.Run(&result)
  fmt.Println(result)
  if len(result) != 1 {
    t.Errorf("Unexpected result %v", result)
  }
  if result[0].Some != 1 || result[0].Another != "helo" {
    t.Errorf("Error Where, expected {1 helo}, got %v", result[0])
  }
}

func TestWhereChain(t *testing.T) {
  st, _ := Connect("sqlite", "test.db")
  to := TestObject{Some: 1, Another: "helo"}
  var result []TestObject
  st.CreateTable(to)
  st.TruncateTable(to)
  st.Create(to)
  st.Where(TestObject{}, "some = 1").Where("another = 'helo'").Run(&result)
  if len(result) != 1 {
    t.Errorf("Unexpected result %v", result)
  }
  if result[0].Some != 1 || result[0].Another != "helo" {
    t.Errorf("Error Where, expected {1 helo}, got %v", result[0])
  }
}

func TestFirst(t *testing.T) {
  st, _ := Connect("sqlite", "test.db")
  to := TestObject{Some: 1, Another: "helo"}
  st.CreateTable(to)
  st.TruncateTable(to)
  st.Create(to)
  var result []TestObject
  st.First(TestObject{}, &result)
  if result[0].Some != 1 || result[0].Another != "helo" {
    t.Errorf("Unexpected first element %v", result)
  }
}
