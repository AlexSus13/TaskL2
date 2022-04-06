/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	strSlice := []string{
		"a4bc10d5e",
		"a100",
		"abcd",
		"45",
		"",
		`qwe\4\5`,
		`qwe\45`,
		`qwe\\5`,
		`qwe\110`,
	}

	for _, s := range strSlice {
		unpacked, err := unpack(s)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("unpacked(%s) = (%s)\n", s, unpacked)
	}
}

func unpack(s string) (string, error) {

	if len(s) == 0 {
		return "", nil
	}

	mNum := make(map[int]string)   //Map где хранятся все числа, ключ - номер элемента.
	mSlash := make(map[int]string) //Map где хранятся все \, ключ - номер элемента.

	for i, r := range s {

		if i == 0 && (r >= '0' && r <= '9') { //Если строка начинается с цифры
			return "", errors.New("Invalid string") //возвращаем ошибку.
		}

		if r >= '0' && r <= '9' { //Если в строке встречается цифра добавляем в map.
			if _, ok := mNum[i-1]; ok { //Если предыдущее тоже цифра то объединяем в в одно число.
				mNum[i] = (mNum[i-1] + string(r))
			} else { //Иначе просто добавляем
				mNum[i] = string(r)
			}
		}

		if string(r) == `\` { //Если в строке встречается \ добавляем в map.
			mSlash[i] = string(r)
		}
	}

	unpacked, err := unpackStr(s, mNum, mSlash)
	if err != nil {
		return "", err
	}
	return unpacked, nil
}

func unpackStr(s string, mNum, mSlash map[int]string) (string, error) {

	var unpackedS string
	if len(mSlash) == 0 { //В строке нет слэш.
		unpacked, err := unpackWithoutSlash(s, mNum)
		if err != nil {
			return "", errors.New("Error when unpacking without slash")
		}
		unpackedS = unpacked
	} else {
		unpacked, err := unpackWithSlash(s, mNum, mSlash)
		if err != nil {
			return "", errors.New("Error when unpacking with slash")
		}
		unpackedS = unpacked
	}
	return unpackedS, nil
}

func unpackWithSlash(s string, mNum, mSlash map[int]string) (string, error) {

	var unpackedS string
	for i, r := range s {
		if string(r) == `\` { //Если попался слэш.

			if _, ok := mSlash[i-1]; ok { //Если у слэша предыдущий тоже символ - тоже слэш.

				if _, ok := mNum[i+1]; ok { //проверяем есть ли после слэша число.
					j := i + 1
					var repeat string
					//Проходимся по map пока не найдем число на которое надо увеличить символ.
					for lm := 0; lm <= len(mNum); lm++ {
						if _, ok := mNum[j]; ok {
							repeat = mNum[j]
							j++
						} else {
							break
						}
					}
					intRepeat, _ := strconv.Atoi(repeat)
					unpackedS += strings.Repeat(string(r), intRepeat)
				} else { //Если поле второго слэша нет числа, просто добавляем
					unpackedS += string(r)
				}

			} else { //Если попался одиночный слэш - пропускаем.
				continue
			}
		} else { //Если символ не слэш.
			if r >= '0' && r <= '9' { //Если символ - цифра.
				if _, ok := mSlash[i-1]; ok { //Предыдущий символ \.

					if _, ok := mSlash[i-2]; ok {
						//Если число после двух слэш - пропускаем, т.к. уже учли.
						continue
					} else { //Если перед числом только 1 слэш.

						if _, ok := mNum[i+1]; ok { //Проверяем есть ли после символа число.
							j := i + 1
							var repeat string
							//Проходимся по map пока не найдем число на которое надо увеличить символ.
							for lm := 0; lm <= len(mNum); lm++ {
								if _, ok := mNum[j]; ok {
									repeat = mNum[j][1:]
									j++
								} else {
									break
								}
							}
							intRepeat, _ := strconv.Atoi(repeat)
							unpackedS += strings.Repeat(string(r), intRepeat)
						} else { //Если числа нет, просто прибавляем.
							unpackedS += string(r)
						}
					}

				} else { //Предыдущий символ не \.
					continue
				}
			} else { //Если не цифра.
				if val, ok := mNum[i+1]; ok { //проверяем есть ли после символа число
					intVal, _ := strconv.Atoi(val)
					unpackedS += strings.Repeat(string(r), intVal)
				} else { //Иначе просто добавляем
					unpackedS += string(r)
				}
			}
		}
	}
	return unpackedS, nil
}

func unpackWithoutSlash(s string, mNum map[int]string) (string, error) {

	var unpackedS string
	for i, r := range s {
		if r >= '0' && r <= '9' { //Если цифра - пропускаем
			continue
		} else {
			if _, ok := mNum[i+1]; ok { //проверяем есть ли после символа число
				j := i + 1
				var repeat string
				//Проходимся по map пока не найдем число на которое надо увеличить символ.
				for lm := 0; lm <= len(mNum); lm++ {
					if _, ok := mNum[j]; ok {
						repeat = mNum[j]
						j++
					} else {
						break
					}
				}
				intRepeat, _ := strconv.Atoi(repeat)
				unpackedS += strings.Repeat(string(r), intRepeat)
			} else { //Иначе просто добавляем
				unpackedS += string(r)
			}
		}
	}

	return unpackedS, nil
}
