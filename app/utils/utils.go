package utils

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"strconv"
	"time"
)

var months = map[time.Month]string{
	time.January:   "января",
	time.February:  "февраля",
	time.March:     "марта",
	time.April:     "апреля",
	time.May:       "мая",
	time.June:      "июня",
	time.July:      "июля",
	time.August:    "августа",
	time.September: "сентября",
	time.October:   "октября",
	time.November:  "ноября",
	time.December:  "декабря",
}

var currencySuffix = map[int]string{
	0: "рублей",
	1: "рубль",
	2: "рубля",
	3: "рубля",
	4: "рубля",
	5: "рублей",
	6: "рублей",
	7: "рублей",
	8: "рублей",
	9: "рублей",
}

func GetMonthWordByOrderNumber(month time.Month) string {
	return months[month]
}

func GetCurrencySuffix(units int) string {
	return currencySuffix[units]
}

func ConvertExcelFileToBytes(f *excelize.File) ([]byte, error) {
	fileName := os.TempDir() + "/" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	err := f.SaveAs(fileName)
	if err != nil {
		log.Printf("[ERROR] Cannot create temporary file #%v", err)
		return nil, err
	}
	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			log.Printf("[WARN] failed to remove %s from FS: %s", fileName, err)
		}
	}()

	res, _ := os.ReadFile(fileName)
	return res, nil
}
