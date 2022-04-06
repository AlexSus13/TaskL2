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
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flags struct {
	A *int
	B *int
	C *int
	c *bool
	i *bool
	v *bool
	F *bool
	n *bool
}

type ObjectGrep struct {
	flags
	data []string
}

func main() {

	ObjectGrep := NewObjectGrep()

	err := ObjectGrep.Grep()
	if err != nil {
		log.Fatal(err)
	}

}

func (ObG *ObjectGrep) Grep() error {

	pattern := flag.Arg(0)
	if len(pattern) == 0 {
		return errors.New("Add a pattern for the search")
	}

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
		ObG.data = append(ObG.data, scanner.Text()) //Добавляем строки в слайс
	}
	//После того, как Scan вернет значение false, метод Err вернет любую ошибку, возникшую во время сканирования,
	//за исключением того, что если это был io.EOF, Err вернет значение nil.
	if err := scanner.Err(); err != nil {
		return err
	}

	if *ObG.C > 0 {
		*ObG.A = *ObG.C
		*ObG.B = *ObG.C
	}

	counter, last := 0, 0
	for n, line := range ObG.data {
		//ok = true если pattern есть в строке или pattern есть с игнорированием регитра(-i).
		ok := strings.Contains(line, pattern) || (*ObG.i && strings.Contains(strings.ToLower(line), strings.ToLower(pattern)))

		if *ObG.F { //точное совпадение со строкой
			if pattern != line {
				continue
			}
		}
		//Если шаблон есть в строке и нет исключающего флага
		//или нет шаблона в строке и есть исключающий флаг
		if (ok && !*ObG.v) || (!ok && *ObG.v) {
			counter++

			for *ObG.B > 0 { //Установлен флаг - показать n строк до строки с шаблоном
				if n-*ObG.B < 0 { //
					*ObG.B--
					continue
				}

				if *ObG.n { //Флаг печати номера строки.
					fmt.Printf("%d-%s\n", n-*ObG.B+1, ObG.data[n-*ObG.B])
				} else {
					fmt.Println(ObG.data[n-*ObG.B])
				}
				*ObG.B--
			}

			if *ObG.n { //Флаг печати номера строки.
				fmt.Printf("%d:%s\n", n+1, line)
			} else {
				fmt.Println(line)
			}

			last = n
		}
	}

	for i := 0; i < *ObG.A; i++ {
		if last+1+i >= len(ObG.data) {
			continue
		}

		if *ObG.n { //Флаг печати номера строки.
			fmt.Printf("%d-%s\n", last+1+i+1, ObG.data[last+1+i])
		} else {
			fmt.Println(ObG.data[last+1+i])
		}
	}

	if *ObG.c { //Фдаг печати колличества строк с шаблоном.
		fmt.Println(counter)
	}

	return nil
}
func NewObjectGrep() *ObjectGrep {

	//Определяем флаг с указанным именем, значением по умолчанию и "способом приминения".
	//Возвращаемое значение - это адрес переменной, в которой хранится значение флага.
	var A = flag.Int("A", 0, "print +N lines after match")
	var B = flag.Int("B", 0, "print +N lines before match")
	var C = flag.Int("C", 0, "(A+B) print ±N lines around the match")
	var c = flag.Bool("c", false, "number of lines")
	var i = flag.Bool("i", false, "ignore case")
	var v = flag.Bool("v", false, "instead of match, exclude")
	var F = flag.Bool("F", false, "exact string match, not a pattern")
	var n = flag.Bool("n", false, "print line number")

	flag.Parse() //Извлечения флага из командной строки.
	//Func считывает значение флага из командной строки и присваивает его содержимое
	//переменной. Нужно вызвать ее *до* использования переменнных
	//иначе они всегда будет содержать значение по умолчанию.
	//Если есть ошибки во время извлечения данных - приложение будет остановлено.

	return &ObjectGrep{
		flags: flags{
			A: A,
			B: B,
			C: C,
			c: c,
			i: i,
			v: v,
			F: F,
			n: n,
		},
		data: make([]string, 0),
	}
}
