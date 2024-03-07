package excel

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"reflect"
)

type Row []any

type TypedRow any

type Sheet struct {
	Name string

	writer      *excelize.StreamWriter
	rowCursor   int
	headerStyle int
}

type File struct {
	file        *excelize.File
	headerStyle int

	sheets map[string]*Sheet
}

func NewFile() *File {
	f := &File{}
	f.file = excelize.NewFile()
	f.headerStyle, _ = f.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#000000",
		},
	})
	f.sheets = make(map[string]*Sheet)
	return f
}

func (f *File) SaveAs(filename string) error {
	defer f.file.Close()
	err := f.file.DeleteSheet("Sheet1")
	if err != nil {
		return err
	}
	// flush all sheets
	for _, sheet := range f.sheets {
		err := sheet.writer.Flush()
		if err != nil {
			return err
		}
	}
	return f.file.SaveAs(filename)
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) NewSheet(name string) (*Sheet, error) {
	// create new sheet
	_, err := f.file.NewSheet(name)
	if err != nil {
		println("Error creating new sheet, err:", err.Error())
		return nil, err
	}
	writer, err := f.file.NewStreamWriter(name)
	if err != nil {
		println("Error creating new sheet writer, err:", err.Error())
		return nil, err
	}
	sheet := &Sheet{
		Name:        name,
		writer:      writer,
		headerStyle: f.headerStyle,
	}

	f.sheets[name] = sheet
	return sheet, nil
}

func (s *Sheet) writeRowWithStyle(row Row, style int) error {
	cells := make([]any, len(row))
	for i, cell := range row {
		cells[i] = excelize.Cell{Value: cell, StyleID: style}
	}
	cellName, _ := excelize.CoordinatesToCellName(1, s.rowCursor+1)
	err := s.writer.SetRow(cellName, cells)
	if err != nil {
		println("Error writing row, err:", err.Error())
		return err
	}
	s.rowCursor++
	return nil
}

func (s *Sheet) WriteRow(row Row) error {
	return s.writeRowWithStyle(row, 0)
}

func (s *Sheet) WriteRows(row []Row) error {
	for _, r := range row {
		err := s.writeRowWithStyle(r, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sheet) WriteTypedRows(row []TypedRow) error {
	if len(row) == 0 {
		return nil
	}
	headers, err := findHeaders(row[0])
	if err != nil {
		return err
	}
	err = s.WriteHeader(headers)
	if err != nil {
		return err
	}
	for _, r := range row {
		err := s.WriteRow(toRow(r))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sheet) WriteHeader(headers []string) error {
	// cast []string to []any
	anyHeaders := make([]any, len(headers))
	for i, header := range headers {
		anyHeaders[i] = header
	}
	return s.writeRowWithStyle(anyHeaders, s.headerStyle)
}

func findHeaders(row TypedRow) ([]string, error) {
	// reflect to find if the row is a struct or a map
	// if struct, use the tag "excel" to find the headers
	// if map, use the keys
	reflected := reflect.ValueOf(row)
	switch reflected.Kind() {
	case reflect.Struct:
		headers := make([]string, 0)
		for i := 0; i < reflected.NumField(); i++ {
			field := reflected.Type().Field(i)
			tag := field.Tag.Get("excel")
			if tag != "" {
				headers = append(headers, tag)
			}
		}
		return headers, nil
	case reflect.Map:
		headers := make([]string, 0)
		for _, key := range reflected.MapKeys() {
			headers = append(headers, key.String())
		}
		return headers, nil
	default:
		return nil, errors.New("cannot find headers for type: " + reflected.Kind().String())
	}
}

func toRow(row TypedRow) Row {
	reflected := reflect.ValueOf(row)
	switch reflected.Kind() {
	case reflect.Struct:
		r := make([]any, 0)
		for i := 0; i < reflected.NumField(); i++ {
			r = append(r, reflected.Field(i).Interface())
		}
		return r
	case reflect.Map:
		r := make([]any, 0)
		for _, key := range reflected.MapKeys() {
			r = append(r, reflected.MapIndex(key).Interface())
		}
		return r
	default:
		return nil
	}
}
