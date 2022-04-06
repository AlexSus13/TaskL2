Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
error - это интерфейс: type error interface {
		          Error() string
		       }
Структура os.PathError реализует этот интерфейс т.к. есть метод:
func (e *PathError) Error() string 

Интерфейс под капотом:
type iface struct {
	tab  *itab
	data unsafe.Pointer
}
В поле tab хранится информация о конкретном типе объекта, который был преобразован в интерфейс.
В поле data - ссылка на реальную область памяти, в которой лежат данные изначального объекта, в нашем случае os.PathError.
fmt.Printf("%#v\n", err) //(*fs.PathError)(nil)
Интерфейс может быть равен nil только если оба поля tab и data не определены.
```
