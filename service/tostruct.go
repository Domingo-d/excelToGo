package service

import (
	"bufio"
	"excelToGo/common"
	"fmt"
	"github.com/fatih/color"
	"github.com/xuri/excelize/v2"
	"os"
	"strings"
)

type (
	ElementInfo struct {
		Name string
		tp   string
	}

	ExcelToGoStruct struct {
		name    string
		tp      string
		link    string
		element []*ExcelToGoStruct
	}
)

func ExcelToGo(path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	defer f.Close()

	sheetInfo := make(map[string]*ExcelToGoStruct)
	for _, sheet := range f.GetSheetList() {
		sheetInfo[sheet] = nil
		dimension, err := f.GetSheetDimension(sheet)
		if err != nil {
			continue
		}

		colCount, _, err := excelize.CellNameToCoordinates(strings.Split(dimension, ":")[1])
		if err != nil {
			continue
		}

		var firstNode *ExcelToGoStruct
	outLoop:
		for colIdx := 1; colIdx <= colCount; colIdx++ {
			colNode := new(ExcelToGoStruct)
			for rowIdx := 1; rowIdx <= 6; rowIdx++ {
				point, err := excelize.CoordinatesToCellName(colIdx, rowIdx)
				if err != nil {
					continue
				}

				cell, err := f.GetCellValue(sheet, point)
				if err != nil {
					continue
				}

				if rowIdx == 1 && colIdx == 1 {
					colNode.tp = cell
					colNode.name = sheet
					firstNode = colNode
					sheetInfo[sheet] = firstNode
					continue outLoop
				}

				switch rowIdx {
				case 2:
					switch cell {
					case "int":
						colNode.tp = "uint32"
					case "float":
						colNode.tp = "float32"
					case "float[]":
						colNode.tp = "[]float32"
					case "int[]":
						colNode.tp = "[]uint32"
					case "str[]":
						colNode.tp = "[]string"
					case "bool":
						colNode.tp = "bool"
					case "string", "str":
						colNode.tp = "string"
					case "ref", "":
					default:
						color.Red("字段类型错误:%s", cell)
					}
				case 3:
					if cell != "" {
						colNode.link = cell
					}
				case 5:
					colNode.name = cell
				case 6:
					if !strings.Contains(cell, "s") {
						continue outLoop
					}
				}
			}

			firstNode.element = append(firstNode.element, colNode)
		}
	}

	saveFile(sheetInfo)

	return nil
}

func saveFile(sheetInfo map[string]*ExcelToGoStruct) {
	path := common.OutputPath
	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}

	file, err := os.Create(path)
	if err != nil {
		color.Red("create file error:%s", err.Error())
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	writer.WriteString("type (\n")
	for _, data := range sheetInfo {
		if len(data.element) <= 0 {
			continue
		}

		writer.WriteString(fmt.Sprintf("\t%s struct {\n", strings.ToUpper(data.name[:1])+data.name[1:]))
		for _, element := range data.element {
			tp := element.tp
			name := element.name
			if name == "ref" || name == "" || name == data.name {
				continue
			}

			name = strings.ToUpper(name[:1]) + name[1:]

			tag := element.name
			if "" != element.link {
				tag = element.link
			}

			if obj, ok := sheetInfo[tag]; ok {
				tp = strings.ToUpper(tag[:1]) + tag[1:]
				var prefix string
				for _, str := range strings.Split(obj.tp, ",") {
					switch str {
					case "arr":
						prefix += "[]"
					case "obj":
						prefix += "map[uint32]"
					}
				}

				tp = prefix + tp
			}

			writer.WriteString(fmt.Sprintf("\t\t%s %s `json:\"%s\"` \n", name,
				tp, tag))
		}
		writer.WriteString("\t}\n")
	}
	writer.WriteString(")\n")

	writer.Flush()
}
