package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flags struct {
	f *int
	d *string
	s *bool
}

type ObjectCut struct {
	flags
	data []string
}

func main() {

	ObjectCut := NewObjectCut()

	err := ObjectCut.Cut()
	if err != nil {
		log.Fatal(err)
	}

}

func (ObC *ObjectCut) Cut() error {

	if *ObC.f == 0 {
		return errors.New(`Using the "-f" flag, specify the column`)
	}

	//NewScanner возвращает новый сканер(*Scanner) для чтения из "os.Stdin".
	scanner := bufio.NewScanner(os.Stdin)
	//Scan переводит Scanner на следующий объект, который затем будет доступен с помощью метода Bytes или Text.
	//Он возвращает значение false, когда сканирование останавливается, либо при достижении конца ввода, либо при ошибке.
	for scanner.Scan() {

		ObC.data = append(ObC.data, scanner.Text()) //Заносим данные в слайс

		for _, line := range ObC.data {
			if strings.Contains(line, *ObC.d) { //Если разделитель присутствует в строке

				dline := strings.Split(line, *ObC.d) //Получим массив элементов строки разделенных разделителем

				if *ObC.f <= len(dline) { //Если указанный номер калонки есть в слайсе
					fmt.Println(dline[*ObC.f-1])
				} else {
					fmt.Println()
				}
			} else { //Если разделител нет в строке
				if *ObC.s { //и установлен флаг вывода строки без разделителей
					fmt.Println(line)
					continue
				}
				break
			}
		}

	}
	//После того, как Scan вернет значение false, метод Err вернет любую ошибку, возникшую во время сканирования,
	//за исключением того, что если это был io.EOF, Err вернет значение nil.
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func NewObjectCut() *ObjectCut {

	//Определяем флаг с указанным именем, значением по умолчанию и "способом приминения".
	//Возвращаемое значение - это адрес переменной, в которой хранится значение флага.
	var f = flag.Int("f", 0, "select fields (columns)")
	var d = flag.String("d", "\t", "use a different delimiter")
	var s = flag.Bool("s", false, "only strings with delimiter")

	flag.Parse() //Извлечения флага из командной строки.
	//Func считывает значение флага из командной строки и присваивает его содержимое
	//переменной. Нужно вызвать ее *до* использования переменнных
	//иначе они всегда будет содержать значение по умолчанию.
	//Если есть ошибки во время извлечения данных - приложение будет остановлено.

	return &ObjectCut{
		flags: flags{
			f: f,
			d: d,
			s: s,
		},
		data: make([]string, 0),
	}
}
