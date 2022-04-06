Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
Сначала мы получим все значения из каналов a и b, а потом будем получать нули.
Используя цикл for v := range c мы итерируемся по значениям из канала который не закрывают.
В канал с больше никто не пишет, т.к. в  merge есть select и каналы a и b закрыты,
==> при попытке получить данные из канала, которых нет, мы получим значение по умолчанию.
Чтобы этого избежать нужно проверять закрытие канала v, opened:= <-c. 
```
