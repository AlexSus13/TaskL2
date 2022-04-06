/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
package main

import (
	"github.com/PuerkitoBio/goquery"
	errorss "github.com/pkg/errors"

	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"strings"
	"fmt"
	"bytes"
)



type flags struct {
	r   *bool   //включение рекурсивной загрузки
	l   *int    //глубина рекурсии
	O   *string //изменить имя сохраненного файла
}

type ObjectWget struct {
        flags
	URL string
        pathToSaveFiles string
	uniqueLinks map[string]struct{}
}
func main() {

	ObjectWget := NewObjectWget()

	//Создаем директорию для сохранения файлов.
	err := os.MkdirAll(ObjectWget.pathToSaveFiles, 0777)
	if err != nil {
		log.Fatal(err)
	}

	err = ObjectWget.Wget()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wget OK")
}

func (ow *ObjectWget) Wget() error {

	if len(os.Args) < 2 {
		return errors.New("Specify the url address")
	}

	if *ow.flags.l < 1 {
		return errors.New("Recursion depth cannot be negative")
	}

	err := GetHTML(ow.URL, ow)
	if err != nil {
		return errorss.Wrap(err, "func GetHTML")
	}

	return nil
}

func GetHTML(link  string, ow *ObjectWget) error {


	//Если ссылка есть в мап, значит страница скачена.
	if _, ok := ow.uniqueLinks[link]; ok {
		return nil
	} else {
		ow.uniqueLinks[link] = struct{}{}
	}

        //Client с Timeout.
        client := http.Client{
                Timeout: 5 * time.Second,
        }

        //Получаем ответ в байтах по url.
        r, err := client.Get(link)
        if err != nil {
                return errorss.Wrap(err, "func client.Get(f.Url)")
        }
        defer r.Body.Close()

        //Проверяем статус код ответа.
        if r.StatusCode != 200 {
                return errors.New("Received non 200 response code")
        }

	FileContents, err := io.ReadAll(r.Body)
	if err != nil {
		return errorss.Wrap(err, "func io.ReadAll(r.Body)")
	}

	//Чтобы прочитать тело запроса дважды, получаем из
	//FileContents []byte - объект типа io.Reader.
	repBody := bytes.NewReader(FileContents)

        //NewDocumentFromReader возвращает документ из io.Reader.
        //Он возвращает ошибку в качестве второго значения, если
        //данные не могут быть проанализированы как html.
        //doc структура типа *Document, представляет HTML-документ.
        doc, err := goquery.NewDocumentFromReader(repBody)
        if err != nil {
                return errorss.Wrap(err, "func goquery.NewDocumentFromReader")
        }

        doc.Find("a").Each(func(index int, selectObj *goquery.Selection) {

                link, _ := selectObj.Attr("href")

		if _, ok := ow.uniqueLinks[link]; !ok {
			ow.uniqueLinks[link] = struct{}{}
		}
        })

	fmt.Println("MAP", ow.uniqueLinks)
	//Сохраняем файл.
	err = SaveFile(FileContents, link, ow)

	//Если выключена рекурсивная загрузка или глубина рекурсии равна 1,
	//то скачаются ссылки только первой страницы страницы.
	if *ow.flags.r || *ow.flags.l == 1 {

		fmt.Println("recurs")
		*ow.flags.l--

		for link, _ := range ow.uniqueLinks {
			err := GetHTML(link, ow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SaveFile(FileContents []byte, link string, ow *ObjectWget) error {

	//"Парсинг URL"
	link = strings.TrimPrefix(link, "http://")
	link = strings.TrimPrefix(link, "http://")
	link = strings.TrimSuffix(link, "/")

	FN := fmt.Sprintf("%s.html", link )
	//Создаем файл по указанному пути.
        file, err := os.Create(ow.pathToSaveFiles + FN)
        if err != nil {
                return err
        }
        defer file.Close()

        file.Write(FileContents)

        return nil
}

func NewObjectWget() *ObjectWget {

	//Определяем флаг с указанным именем, значением по умолчанию и "способом приминения".
	//Возвращаемое значение - это адрес переменной, в которой хранится значение флага.
	var r = flag.Bool("r", false, "activation recursive loading")
	var l = flag.Int("l", 1, "recursion depth")
	var O = flag.String("n", "index.html", "change name saved file")
	//Получаем url.
	URL := os.Args[len(os.Args)-1]

	flag.Parse() //Извлечения флага из командной строки.
	//Func считывает значение флага из командной строки и присваивает его содержимое
	//переменной. Нужно вызвать ее *до* использования переменнных
	//иначе они всегда будет содержать значение по умолчанию.
	//Если есть ошибки во время извлечения данных - приложение будет остановлено.

        return &ObjectWget{
                flags: flags{
                        r: r,
                        l: l,
                        O: O,
                },
		URL: URL,
		uniqueLinks: make(map[string]struct{}),
		pathToSaveFiles: "newWget/",
        }
}
