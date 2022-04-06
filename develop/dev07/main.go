package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

func main() {

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})

		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}

	start := time.Now()
	a := <-or(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
		sig(4*time.Second),
		sig(5*time.Second),
	)

	fmt.Printf("fone after %v\n", time.Since(start))
}

func or(channels ...<-chan interface{}) <-chan interface{} {

	out := make(chan interface{}) //Создаем выходной канал
	defer close(out)

	wg := &sync.WaitGroup{}

	for _, ch := range channels { //перебираем каналы
		wg.Add(1)
		go func(ch <-chan interface{}) {
			defer wg.Done()
			for {
				v, ok := <-ch //Проверка канала на закрытие
				if ok {       //ok = true, следоватеьно открыт
					out <- v
				} else { //закыт
					break //выходим из цикла for
				}
			}
			/*
				for v := range ch {//range получает значения из канала пока он не закрыт
					out <- v
				}
			*/
		}(ch)
	}

	wg.Wait() //Ждем заершения горутин

	return out
}
