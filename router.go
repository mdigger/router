package router

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	// PathDelimeter defines the path separator.
	PathDelimeter = "/"
	// NamedParamFlag is used to define a named parameter in the path.
	NamedParamFlag = byte(':')
	// CatchAllParamFlag is used to define a dynamic named parameter in route.
	CatchAllParamFlag = byte('*')
	// Splitter normalizes the path and returns it in the form of parts. If you
	// want, you can replace this function on his own.
	Splitter = func(url string) []string {
		return strings.SplitAfter(strings.TrimPrefix(url, PathDelimeter),
			PathDelimeter)
	}
)

// Paths describes the structure for quick selection handler for request path.
// Supports both a static route and path parameters.
type Paths struct {
	// хранилище статических путей, без параметров;
	// в качестве ключа используется полный путь
	static map[string]interface{}
	// хранит информацию о путях с параметрами;
	// в качестве ключа используется общее количество элементов пути
	fields map[uint16]records
	// максимальное количество частей пути во всех определениях
	maxParts uint16
	// позиция, в которой встречается самый ранний динамический параметр
	catchAll uint16
}

// Add adds a new handler for the specified path. In the description of the way
// to use named parameters (starts with': 'character) and the final a named
// parameter (starts with '*'), which indicates that the URL can longer. In the
// latter case all the rest of the path will be included in this setting. A
// starred parameter, if specified, must be the the last parameter of the path.
//
// Returns an error if the handler is not defined (nil), if the number of
// elements of a URL path greater than 32768 or option with an asterisk is not
// used in the last path element.
//
// ATTENTION! When adding a path is not verified by its uniqueness from the
// point of view named parameters. Therefore, it is possible to add two
// different handler for the same path. For example:
//	/:user/:id/:name
//	/:user/:name/:id
// In this case, the error will not happen, just volunteer to be the first
// added the handler, and the other will never be called.
//
// On the other hand, absolutely correctly fulfilled the following situation:
// 	/:user/:name
// 	/:user/test
func (r *Paths) Add(url string, handler interface{}) error {
	if handler == nil {
		return errors.New("nil handler")
	}
	parts := Splitter(url) // нормализуем путь и разбиваем его на части
	// проверяем, что количество получившихся частей не превышает максимально
	// поддерживаемое количество
	level := uint16(len(parts)) // всего элементов пути
	if level > (1<<15 - 1) {
		return fmt.Errorf("path parts overflow: %d", len(parts))
	}
	// считаем количество параметров в определении пути
	var params uint16
	for i, value := range parts {
		if value == "" {
			continue // пропускаем пустые пути
		}
		switch value[0] {
		case NamedParamFlag:
			params++ // увеличиваем счетчик параметров
		case CatchAllParamFlag:
			// такой параметр должен быть самым последним в определении путей
			if uint16(i) != level-1 {
				return errors.New("catch-all parameter must be last")
			}
			params |= 1 << 15 // взводим флаг динамического параметра
			// запоминаем позицию самого раннего встреченного динамического
			// параметра во всех добавленных путях
			if r.catchAll == 0 || r.catchAll > level {
				r.catchAll = level
			}
		}
	}
	// если в пути нет параметров, то добавляем в статические обработчики
	if params == 0 {
		if r.static == nil {
			r.static = make(map[string]interface{})
		}
		r.static[strings.Join(parts, "")] = handler
		return nil
	}
	// запоминаем максимальное количество элементов пути во всех определениях
	if r.maxParts < level {
		r.maxParts = level
	}
	// инициализируем динамические пути, если не сделали этого раньше
	if r.fields == nil {
		r.fields = make(map[uint16]records)
	}
	// добавляем в массив обработчиков с таким же количеством параметров
	r.fields[level] = append(r.fields[level], &record{params, parts, handler})
	sort.Stable(r.fields[level]) // сортируем по количеству параметров
	return nil
}

// Lookup returns the handler and the list of named parameters with their
// values. If a suitable handler is found, it returns nil.
func (r *Paths) Lookup(url string) (interface{}, Params) {
	parts := Splitter(url) // нормализуем путь и разбиваем его на части
	// сначала ищем среди статических путей; если статические пути не
	// определены, то пропускаем проверку
	if r.static != nil {
		if handler, ok := r.static[strings.Join(parts, "")]; ok {
			return handler, nil
		}
	}
	// если пути с параметрами не определены, то на этом заканчиваем проверку
	if r.fields == nil {
		return nil, nil
	}
	length := uint16(len(parts)) // вычисляем количество элементов пути
	// наши определения могут быть короче, если используются catchAll параметры,
	// поэтому вычисляем с какой длины начинать
	var total uint16
	// если длина запроса больше максимальной длины определений, то нужно
	// замахиваться на меньшее...
	if length > r.maxParts {
		// если нет динамических параметров, то ничего и не подойдет,
		// потому что наш запрос явно длиннее
		if r.catchAll == 0 {
			return nil, nil
		}
		total = r.maxParts // начнем с максимального определения пути
	} else {
		total = length // наш запрос короче самого длинного определения
	}
	// запрашиваем список обработчиков для такого же количества элементов пути
	for l := total; l > 0; l-- {
		records := r.fields[l] // получаем определения путей для данной длины
		// если обработчики для такой длины пути не зарегистрированы, то...
		if len(records) == 0 {
			// проверяем, что на этом уровне динамические пути еще встречаются
			if l < r.catchAll {
				break // больше нет динамических параметров дальше
			}
			// переходим к более короткому пути
			continue
		}
	nextRecord:
		// обработчики есть — перебираем все записи с ними
		for _, record := range records {
			// если наш путь длиннее обработчика, а он не содержит catchAll
			// параметра, то он точно нам не подойдет
			if l < length && record.params>>15 != 1 {
				continue
			}
			// здесь мы будем собирать значения параметров к данному запросу
			// если ранее они были не пустые от другого обработчика, то
			// сбрасываем их
			var params Params
		params:
			// перебираем все части пути, заданные в обработчике
			for i, part := range record.parts {
				switch part[0] {
				case byte(':'): // это одиночный параметр
					params = append(params, Param{
						// убираем ':' в начале и возможный '/' в конце
						Key: strings.TrimSuffix(part[1:], PathDelimeter),
						// элемент пути без возможного '/' в конце
						Value: strings.TrimSuffix(parts[i], PathDelimeter),
					})
					continue // переходим к следующему элементу пути
				case byte('*'): // это параметр, который заберет все
					params = append(params, Param{
						Key: part[1:], // исключаем '*' из имени
						// добавляем весь оставшийся путь
						Value: strings.Join(parts[i:], ""),
					})
					break params // больше ловить нечего — нашли
				}
				// статическая часть пути не совпадает с запрашиваемой
				if part != parts[i] {
					// переходим к следующему обработчику
					continue nextRecord
				}
			}
			// возвращаем найденный обработчик и заполненные параметры
			return record.handler, params
		}
	}
	// сюда мы попадаем, если так ничего подходящего и не нашли
	return nil, nil
}

// Path returns a list of path elements associated with this processor.
// If the handler is associated with multiple paths, return the first.
func (r *Paths) Path(handler interface{}) []string {
	// перебираем статические пути
	for url, h := range r.static {
		if h == handler {
			return Splitter(url) // нашли нужный адрес - возвращаем элементы пути
		}
	}
	// перебираем все пути с параметрами
	for _, records := range r.fields {
		for _, record := range records {
			// сравниваем адреса методов
			if handler == record.handler {
				return record.parts // возвращаем элементы пути
			}
		}
	}
	return nil // данный обработчик не зарегистрирован
}
