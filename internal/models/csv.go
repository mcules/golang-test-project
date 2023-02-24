package models

import (
	"bytes"
	"encoding/csv"
	"errors"
	"log"
	"os"
)

type CSVWriter [][]string

func (record *CSVWriter) Init(header []string) {
	record.Add(header)
}

func (record *CSVWriter) Add(line []string) {
	*record = append(*record, line)
}

func (record *CSVWriter) Get() ([]byte, error) {
	if *record == nil || len(*record) == 0 {
		return nil, errors.New("records cannot be nil or empty")
	}

	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	if err := csvWriter.WriteAll(*record); err != nil {
		return nil, err
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (record *CSVWriter) SaveFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)

	content, err := record.Get()
	if err != nil {
		return err
	}

	_, err = f.WriteString(string(content))
	if err != nil {
		return err
	}

	return nil
}
