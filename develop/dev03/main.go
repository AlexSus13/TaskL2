package main

import (
	"flag"
	"fmt"
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flags struct {
	k *int
	n *bool
	r *bool
	u *bool
}

type ObjectSort struct {
	flags
	data []string
}

func main() {

	ObjectSort := NewObjectSort()

	err := ObjectSort.Sort()
	if err != nil {
		log.Fatal(err)
	}

}

func (ObS *ObjectSort) Sort() error {

	file, err := os.Open(os.Args[len(os.Args)-1]) //Открываем файл указанный в командной строке.
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) //NewScanner возвращает новый сканер(*Scanner) для чтения из "file". 
					  //Функция разделения по умолчанию использует ScanLines.

	//Scan переводит Scanner на следующий объект, который затем будет доступен с помощью метода Bytes или Text. 
	//Он возвращает значение false, когда сканирование останавливается, либо при достижении конца ввода, либо при ошибке. 
	for scanner.Scan() {
		ObS.data = append(ObS.data, scanner.Text())//Добавляем строки в слайс
	}
	//После того, как Scan вернет значение false, метод Err вернет любую ошибку, возникшую во время сканирования, 
	//за исключением того, что если это был io.EOF, Err вернет значение nil.
	if err := scanner.Err(); err != nil {
		return err
	}

	if *ObS.flags.n { //Сортировать по числовому значению.
		sort.Slice(ObS.data, func(i, j int) bool {
			columnElement1, err := strconv.Atoi(strings.Split(ObS.data[i], " ")[*ObS.flags.k-1])
			if err != nil {
				fmt.Printf("Integer conversion error: %v\n", err)
			}
			columnElement2, err := strconv.Atoi(strings.Split(ObS.data[j], " ")[*ObS.flags.k-1])
			if err != nil {
				fmt.Printf("Integer conversion error: %v\n", err)
			}
			return columnElement1 < columnElement2
		})
	} else {
		sort.Slice(ObS.data, func(i, j int) bool {
			columnElement1 := strings.Split(ObS.data[i], " ")[*ObS.flags.k-1]
			columnElement2 := strings.Split(ObS.data[j], " ")[*ObS.flags.k-1]
			return columnElement1 < columnElement2
		})
	}

	if *ObS.flags.r { //сортировать в обратном порядке
		for i, j := 0, len(ObS.data)-1; i < j; i, j = i+1, j-1 {
			ObS.data[i], ObS.data[j] = ObS.data[j], ObS.data[i]
		}
	}

	if *ObS.flags.u { //не выводить повторяющиеся строки
		uniqueLine := make(map[string]struct{})
		uniqueObSData := make([]string, 0)

		for _, line := range ObS.data {
			if _, ok := uniqueLine[line]; !ok {
				uniqueLine[line] = struct{}{}
				uniqueObSData = append(uniqueObSData, line)
			}
		}

		ObS.data = uniqueObSData
	}

	for _, line := range ObS.data {
		fmt.Println(line)
	}

	return nil
}

func NewObjectSort() *ObjectSort {

	//Определяем флаг с указанным именем, значением по умолчанию и "способом приминения".
	//Возвращаемое значение - это адрес переменной, в которой хранится значение флага.
	var k = flag.Int("k", 1, "specifying a column to sort")
	var n = flag.Bool("n", false, "sort by numeric value")
	var r = flag.Bool("r", false, "sort in reverse order")
	var u = flag.Bool("u", false, "do not output duplicate lines")

	flag.Parse() //Извлечения флага из командной строки.
	//Func считывает значение флага из командной строки и присваивает его содержимое
	//переменной. Нужно вызвать ее *до* использования переменнных
	//иначе они всегда будет содержать значение по умолчанию.
	//Если есть ошибки во время извлечения данных - приложение будет остановлено.

	return &ObjectSort{
		flags: flags{
			k: k,
			n: n,
			r: r,
			u: u,
		},
		data: make([]string, 0),
	}
}
