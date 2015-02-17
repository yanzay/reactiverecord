package reactiverecord

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/eaigner/jet"
)

type ReactiveRecord struct {
	Adapter SQLAdapter
}

type SQLAdapter interface {
	Query(string) jet.Runnable
}

func Connect(database string, connection string) (rec *ReactiveRecord, err error) {
	rec = new(ReactiveRecord)
	switch database {
	case "sqlite":
		rec.Adapter, err = NewSqliteAdapter(connection)
	}
	return rec, err
}

func (rec *ReactiveRecord) Query(query string) jet.Runnable {
	fmt.Printf("[RR] Query: %s\n", query)
	return rec.Adapter.Query(query)
}

func (rec *ReactiveRecord) CreateTable(obj interface{}) error {
	query := rec.CreateTableSQL(obj)
	err := rec.Query(query).Run()
	return err
}

func (rec *ReactiveRecord) TruncateTable(obj interface{}) {
	query := rec.TruncateTableSQL(obj)
	rec.Query(query).Run()
}

func (rec *ReactiveRecord) Create(obj interface{}) error {
	query := rec.InsertSQL(obj)
	err := rec.Query(query).Run()
	if err != nil {
		fmt.Printf("Fail to save %v, error: %v", obj, err)
	}
	return err
}

func (rec *ReactiveRecord) All(obj interface{}, arr interface{}) error {
	query := rec.SelectAllSQL(obj)
	err := rec.Query(query).Rows(arr)
	return err
}

func (rec *ReactiveRecord) Where(obj interface{}, condition string) *ReactiveRelation {
	rel := ReactiveRelation{Rec: rec}
	rel.TableName = rec.tableName(obj)
	rel.Fields = rec.attributeNamesString(obj)
	rel.EmptyObject = obj
	rel.Where(condition)
	return &rel
}

func (rec *ReactiveRecord) First(obj interface{}, result interface{}) {
	rel := ReactiveRelation{Rec: rec}
	rel.TableName = rec.tableName(obj)
	rel.Fields = rec.attributeNamesString(obj)
	rel.First(result)
}

func (ar ReactiveRecord) CreateTableSQL(obj interface{}) string {
	attrs := ar.attributes(obj)
	tableName := ar.tableName(obj)
	fieldStrings := make([]string, 0)
	for fieldName, fieldType := range attrs {
		var field string
		switch fieldType.Kind() {
		case reflect.Int:
			field = "INT"
		case reflect.String:
			field = "TEXT"
		case reflect.Float32:
			field = "REAL"
		}
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s %s", fieldName, field))
	}
	fieldsString := strings.Join(fieldStrings, ", ")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, fieldsString)
	return query
}

func (rec ReactiveRecord) TruncateTableSQL(obj interface{}) string {
	tableName := rec.tableName(obj)
	return fmt.Sprintf("DELETE FROM %s;", tableName)
}

func (ar ReactiveRecord) InsertSQL(obj interface{}) string {
	attrs := ar.attributeNames(obj)
	tableName := ar.tableName(obj)
	values := make([]string, 0)
	for _, attrName := range attrs {
		field := reflect.ValueOf(obj).FieldByName(attrName)
		value := getStringValueOfField(field)
		values = append(values, value)
	}
	attrString := strings.Join(attrs, ", ")
	valueString := strings.Join(values, ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, attrString, valueString)
	return query
}

func (ar ReactiveRecord) SelectAllSQL(obj interface{}) string {
	tableName := ar.tableName(obj)
	attrs := ar.attributeNames(obj)
	attrString := strings.Join(attrs, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s;", attrString, tableName)
	return query
}

func (ar ReactiveRecord) attributes(obj interface{}) map[string]reflect.Type {
	val := reflect.TypeOf(obj)
	fieldsNumber := val.NumField()
	attrs := make(map[string]reflect.Type)
	for i := 0; i < fieldsNumber; i++ {
		field := val.Field(i)
		attrs[field.Name] = field.Type
	}
	return attrs
}

func (ar ReactiveRecord) attributeNames(obj interface{}) []string {
	val := reflect.TypeOf(obj)
	fieldsNumber := val.NumField()
	attrNames := make([]string, 0)
	for i := 0; i < fieldsNumber; i++ {
		attrNames = append(attrNames, val.Field(i).Name)
	}
	return attrNames
}

func (rec ReactiveRecord) attributeNamesString(obj interface{}) string {
	return strings.Join(rec.attributeNames(obj), ", ")
}

func (ar ReactiveRecord) tableName(obj interface{}) string {
	name := reflect.TypeOf(obj).String()
	name = strings.Split(name, ".")[1]
	re := regexp.MustCompile("([a-z]+)([A-Z][a-z])")
	name = re.ReplaceAllString(name, "${1}_${2}")
	return strings.ToLower(name)
}

func getStringValueOfField(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Int:
		return strconv.Itoa(int(val.Int()))
	case reflect.String:
		return fmt.Sprintf("'%s'", val.String())
	case reflect.Float32:
		return fmt.Sprintf("%v", val.Float())
	}
	return ""
}
