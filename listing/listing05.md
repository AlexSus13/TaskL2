Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error - это интерфейс: type error interface {
                          Error() string
                       }
Структура customError реализует этот интерфейс т.к. есть метод:
func (e *customError) Error() string 

Интерфейс под капотом:
type iface struct {
        tab  *itab
        data unsafe.Pointer
}
В поле tab хранится информация о конкретном типе объекта, который был преобразован в интерфейс.
В поле data - ссылка на реальную область памяти, в которой лежат данные изначального объекта, в нашем случае customError.
fmt.Printf("%#v\n", err) //(*main.customError)(nil)
Интерфейс может быть равен nil только если оба поля tab и data не определены.
```
