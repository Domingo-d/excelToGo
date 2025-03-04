package json

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type (
	JSONNode struct {
		Type     string
		Value    any
		Children []*JSONNode
		Tag      string
	}
)

var (
	numberRegexp = regexp.MustCompile(`^\d+$`)
)

func parseJSON(data []byte) (*JSONNode, error) {
	var raw any
	if err := jsoniter.Unmarshal(data, &raw); nil != err {
		return nil, err
	}

	return parseNode(raw)
}

func GetKindType(v any) string {
	cType := "string"
	tp := reflect.TypeOf(v).Kind()
	switch tp {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		cType = "uint32"
	case reflect.String:
		cType = "string"
	case reflect.Bool:
		cType = "bool"
	case reflect.Array, reflect.Slice:
		cType = "array"
	case reflect.Struct:
		cType = "struct"
	case reflect.Map:
		cType = "map"
	}

	return cType
}

func parseNode(data any) (*JSONNode, error) {
	switch v := data.(type) {
	case map[string]any:
		return parseObject(v)
	case []any:
		return parseArray(v)
	default:
		return &JSONNode{Type: GetKindType(v), Value: v}, nil
	}
}

func parseObject(data map[string]any) (*JSONNode, error) {
	var fields []*JSONNode
	for key, value := range data {

		child, err := parseNode(value)
		if err != nil {
			return nil, err
		}

		fields = append(fields, child)
		if numberRegexp.MatchString(key) {
			return &JSONNode{
				//Type:     fmt.Sprintf("map[uint32]*%s", child.Type),
				Type:     "struct",
				Value:    data,
				Children: fields,
			}, nil
		}

		child.Tag = fmt.Sprintf("json:\"%s\"", key)
	}
	return &JSONNode{
		Type:     "",
		Value:    data,
		Children: fields,
	}, nil
}

func parseArray(data []any) (*JSONNode, error) {
	var children []*JSONNode
	if len(data) <= 0 {
		return nil, nil
	}

	child, err := parseNode(data[0])
	if err != nil {
		return nil, err
	}
	children = append(children, child)

	return &JSONNode{
		//Type:     fmt.Sprintf("[]*%s", child.Type),
		Type:     "array",
		Value:    data,
		Children: children,
	}, nil
}

func generateStruct(node *JSONNode) (string, error) {
	switch node.Type {
	case "object":
		return generateObject(node)
	case "array":
		return generateArray(node)
	default:
		return fmt.Sprintf("%s %s `%s`", strings.ToTitle(node.Type), node.Type, node.Tag), nil
	}
}

func generateObject(node *JSONNode) (string, error) {
	var fields []string
	for _, child := range node.Children {
		field, err := generateStruct(child)
		if nil != err {
			return "", err
		}

		fields = append(fields, fmt.Sprintf("%s %s `%s`", strings.ToTitle(child.Tag), field, child.Tag))
	}

	return fmt.Sprintf("%s struct {\n%s\n}", node.Tag, strings.Join(fields, "\n")), nil
}

func generateArray(node *JSONNode) (string, error) {
	if len(node.Children) <= 0 {
		return "", nil
	}

	elementType, err := generateStruct(node.Children[0])
	if nil != err {
		return "", err
	}

	return fmt.Sprintf("%s []*%s `%s`", node.Tag, elementType, node.Tag), nil
}

func jsonToGo(excelName string) error {
	file, err := os.Open(excelName)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
		panic(err)
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	structInfo, err := parseJSON(content)
	if err != nil {
		return err
	}

	goFileName := strings.TrimSuffix(structInfo.Type, ".json")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("type(\n\t %s struct {\n", goFileName))
	structCode, err := generateStruct(structInfo)
	if nil != err {
		return err
	}
	sb.WriteString(structCode)
	sb.WriteString("\n}\n)")

	goFile, err := os.OpenFile(goFileName+".go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Error creating Go file: %s", err)
	}

	defer goFile.Close()

	goFile.WriteString(sb.String())

	return nil
}
