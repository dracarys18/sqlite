package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
)

type PageType byte

const (
	InteriorIndex = 2
	InteriorTable = 5
	LeafIndex     = 10
	LeafTable     = 13
)

func (t PageType) ToString() string {
	var s string = ""

	switch t {
	case InteriorIndex:
		s = "InteriorIndex"
	case InteriorTable:
		s = "InteriorTable"
	case LeafIndex:
		s = "LeafIndex"
	case LeafTable:
		s = "LeafTable"
	default:
		log.Println("Invalid option")
	}

	return s
}

type DatabaseHeader struct {
	PageSize uint16
}

func ParseDatabaseHeader(file io.ReadSeeker) (DatabaseHeader, error) {
	header := make([]byte, 100)

	_, err := file.Read(header)
	if err != nil {
		log.Fatal(err)
	}

	var pageSize uint16
	if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
		fmt.Println("Failed to read integer:", err)
		return DatabaseHeader{}, errors.New("Failed to read the pagesize")
	}

	return DatabaseHeader{PageSize: pageSize}, nil
}

type PageHeader struct {
	Type     PageType
	CellNums uint16
}

func parseByte[T any](file io.ReadSeeker, len int) T {
	buf := make([]byte, len)

	if _, err := file.Read(buf); err != nil {
		log.Fatal(err)
	}

	var v T

	if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &v); err != nil {
		log.Fatal(err)
	}

	return v
}

func ParsePageHeader(file io.ReadSeeker) PageHeader {
	// Move out of database headers
	if _, err := file.Seek(100, io.SeekStart); err != nil {
		log.Fatal("Invalid file cannot find database header")
	}

	pageType := parseByte[PageType](file, 1)

	//Skip the freeblock page
	if _, err := file.Seek(2, io.SeekCurrent); err != nil {
		log.Fatal("Can't skip the freeblock byte")
	}

	cellNums := parseByte[uint16](file, 2)
	return PageHeader{Type: pageType, CellNums: cellNums}
}

type LeafTablePage struct {
	Header PageHeader
}

func ParseLeafTablePage(file io.ReadSeeker) LeafTablePage {
	header := ParsePageHeader(file)
	return LeafTablePage{Header: header}
}
