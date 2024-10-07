package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	validateWrongPrice     = "validateWrongPrice"
	validateWrongFormat    = "validateWrongFormat"
	validateWrongPhone     = "validateWrongPhone"
	validateWrongMkInfo    = "validateWrongMkInfo"
	validateWrongClientFio = "validateWrongClientFio"
)

var ValidationMessages = map[string]string{
	validateWrongPrice:     "Не удалось определить Цену",
	validateWrongFormat:    "Не удалось определить Формат обучения",
	validateWrongPhone:     "Не удалось определить Телефон",
	validateWrongMkInfo:    "Не удалось определить Информацию по МК",
	validateWrongClientFio: "Не удалось определить Фио клиента",
}

type LeadRecord struct {
	Price     float32
	Format    string
	Phone     string
	MkInfo    string
	Manager   string
	ClientFio string
}

func isTextFormatted(text string) bool {
	count := 0

	for _, char := range text {
		if char == ':' {
			count++
		}
	}

	return count > 2
}

func (lr *LeadRecord) Do(text string) error {
	if isTextFormatted(text) == true {
		return lr.formatParsing(text)
	}
	lowerText := strings.ToLower(text)

	return lr.intuitiveParsing(lowerText)
}

func (lr *LeadRecord) formatParsing(text string) error {
	// Создаем карту для хранения ключей и значений
	parsedData := make(map[string]string)

	// Создаем новый сканер для построчного чтения
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		// Разделяем строку на ключ и значение
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.ToLower(strings.TrimSpace(parts[0])) // Ключ - текст до ':'
			value := strings.TrimSpace(parts[1])                // Значение - текст после ':'
			parsedData[key] = value                             // Добавляем в карту
		}
	}

	// Обработка ошибок
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении текста:", err)
		return err
	}

	price, err := strconv.ParseFloat(parsedData["внесено"], 32)
	if err != nil {
		return err
	}
	lr.Price = float32(price)
	lr.Format = parsedData["формат обучения"]
	lr.Phone = parsedData["телефон"]
	lr.MkInfo = parsedData["мк"]
	lr.Manager = parsedData["менеджер"]
	lr.ClientFio = parsedData["фио"]

	lr.cleanData()

	return nil
}

func (lr *LeadRecord) intuitiveParsing(text string) error {
	lr.MkInfo = extractMkInfo(text)
	lr.Phone = extractPhone(text)
	lr.ClientFio = getFioBeforePhone(text, lr.Phone)
	lr.Format = extractFormat(text)

	//убираем цифры телефона
	replacedText := strings.ReplaceAll(text, lr.Phone, "")
	lr.Price = extractPrice(replacedText)

	//вырезаем все что удалось получить
	replacedText = strings.ReplaceAll(text, lr.Format, "")
	replacedText = strings.ReplaceAll(text, lr.MkInfo, "")

	lr.Manager = extractManager(replacedText)
	lr.cleanData()

	return nil
}

func (lr *LeadRecord) Validate() error {
	if lr.Price <= 0 {
		return errors.New(validateWrongPrice)
	}
	if lr.Format == "" {
		return errors.New(validateWrongFormat)
	}
	if lr.Phone == "" {
		return errors.New(validateWrongPhone)
	}
	if lr.MkInfo == "" {
		return errors.New(validateWrongMkInfo)
	}
	if lr.ClientFio == "" {
		return errors.New(validateWrongClientFio)
	}

	return nil
}

func (lr *LeadRecord) cleanData() {
	// Перебираем все поля структуры
	val := reflect.ValueOf(lr).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.String {
			// Удаляем переносы строк и специальные символы
			cleanedValue := strings.ReplaceAll(field.String(), "\n", "")
			field.SetString(cleanedValue) // Устанавливаем очищенное значение
		}
	}
}
func getFioBeforePhone(text, phone string) string {
	var fio string
	// Находим индекс номера телефона в тексте
	phoneIndex := strings.Index(text, phone)
	if phoneIndex != -1 {
		// Получаем текст перед номером телефона
		beforePhone := text[:phoneIndex]
		// Разбиваем на слова
		words := strings.Fields(beforePhone)
		// Проверяем, есть ли как минимум два слова
		if len(words) >= 2 {
			fio = capitalizeWords(words[len(words)-2]) + " " + capitalizeWords(words[len(words)-1])
		} else if len(words) == 1 {
			fio = capitalizeWords(words[0])
		}
	}

	return fio
}

func extractMkInfo(text string) string {
	words := strings.Fields(text)
	if len(words) < 3 {
		return text
	}

	return capitalizeWords(words[1]) + " " + capitalizeWords(words[2])
}

func extractManager(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		if strings.ToLower(word) == "менеджер" {
			if i+1 < len(words) {
				return capitalizeWords(words[i+1])
			}

			return capitalizeWords(word)
		}
	}

	return ""
}

func capitalizeWords(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		if len(word) > 0 {
			// Преобразуем первую букву в заглавную, обрабатывая кириллицу
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

func extractPrice(text string) float32 {
	// Регулярное выражение для поиска всех чисел, включая числа с буквой "к"
	re := regexp.MustCompile(`\d+к?|\d+`)
	matches := re.FindAllString(text, -1)

	if len(matches) > 0 {

		priceStr := matches[0]
		var price float32
		// Проверяем, есть ли буква "к" в конце
		if len(priceStr) > 1 && rune(priceStr[len(priceStr)-1]) == 186 {
			// Убираем "к" и умножаем на 1000
			numericPart := strings.ReplaceAll(priceStr, "к", "")
			numericValue, err := strconv.ParseFloat(numericPart, 32)
			if err != nil {
				log.Println("Error parsing number:", err)
				return 0
			}

			price = float32(numericValue * 1000)
		} else {
			// Если "к" нет, просто парсим число
			numericValue, err := strconv.ParseFloat(priceStr, 32)
			if err != nil {
				log.Println("Error parsing number:", err)
				return 0
			}
			price = float32(numericValue)
		}

		return price
	}

	return 0
}

func extractPhone(text string) string {
	return find(
		`(?:\+375|(?:\+7|7|8|\b)?)\s*(\d{3})\s*(\d{3})[-\s]?(\d{2})[-\s]?(\d{2})`,
		text,
	)
}

func extractFormat(text string) string {
	result := find(
		`(?i)с отработкой`,
		text,
	)

	if len(result) > 0 {
		return "C отработкой"
	}

	return "Без отработки"
}

func find(regStr string, text string) string {
	re := regexp.MustCompile(regStr)

	matches := re.FindStringSubmatch(text)

	if len(matches) > 0 {
		return matches[0]
	}

	return ""
}
