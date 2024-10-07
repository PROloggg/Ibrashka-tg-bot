package files

import (
	"app/parser"
	"app/storage"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Storage struct {
	basePath string
}

const (
	DefaultPerm   = 0774
	FileExtension = ".csv"
)

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {
		if err != nil {
			log.Print(err)
		}
	}()

	leadRecord := &parser.LeadRecord{}

	err = leadRecord.Do(page.Text)
	if err != nil {
		return err
	}

	err = leadRecord.Validate()
	if err != nil {
		return err
	}

	filePath := s.basePath

	if err := os.MkdirAll(filePath, DefaultPerm); err != nil {
		return err
	}

	filePath = filepath.Join(filePath, leadRecord.MkInfo+" "+fileName(page)+FileExtension)

	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, DefaultPerm)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header only if the file is newly created
	if fileInfo, err := os.Stat(filePath); err == nil && fileInfo.Size() == 0 {
		header := header()
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	// Write the page data
	record := record(leadRecord, page)

	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}

func header() []string {
	return []string{
		"МК",
		"Телефон",
		"ФИО",
		"Внесено",
		"Формат обучения",
		"Менеджер",
		"",
		"",
		"",
		"",
		"Сообщение",
		"Пользователь",
		"Дата добавления",
	}
}

func record(leadRecord *parser.LeadRecord, page *storage.Page) []string {
	managerStr := leadRecord.Manager
	if len(managerStr) == 0 {
		managerStr = page.UserName
	}

	return []string{
		leadRecord.MkInfo,
		leadRecord.Phone,
		leadRecord.ClientFio,
		fmt.Sprintf("%.2f", leadRecord.Price),
		leadRecord.Format,
		managerStr,
		"",
		"",
		"",
		"",
		page.Text + " " + page.PictureUrl,
		page.UserName,
		page.CreatedAt.Format("02-01-2006 15:04:05"),
	}
}

func fileName(p *storage.Page) string {
	return p.CreatedAt.Format("01-2006")
}
