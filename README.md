# Простой маршрутизатор с поддержкой именованных параметров

[![GoDoc](https://godoc.org/github.com/mdigger/router?status.svg)](https://godoc.org/github.com/mdigger/router)
[![Build Status](https://travis-ci.org/mdigger/router.svg)](https://travis-ci.org/mdigger/router)
[![Coverage Status](https://coveralls.io/repos/github/mdigger/router/badge.svg?branch=master)](https://coveralls.io/github/mdigger/router?branch=master)

router содержит достаточно простой универсальный "маршрутизатор",
который является заготовкой для замены `http.ServeMux` с поддержкой именованных
параметров в пути.

Текущая реализация позволяет ассоциировать с именованными путями любый
объекты, не только обработчики HTTP-запросов. Для именованных параметров
используется маркер `:`, а для динамических — `*`. Хотя, это все
настраиваемое. Естественно, статические пути без параметров тоже
поддерживаются.

Примеры задаваемых путей:

	/user/:name
	/user/test
	/files/*filename
	/repos/:owner/:repo/pulls

### Как использовать?

Библиотека достаточно абстрактна и не привязана впрямую к `http`. Так что можно 
использовать не только по ее прямому предназначению:

```go
package main

import (
	"fmt"
	"github.com/mdigger/router"
)

func main() {
	var paths router.Paths
	paths.Add("/users", "usersList")
	paths.Add("/users/:name", "userName")
	paths.Add("/users/me", "userMe")

	name, params := paths.Lookup("/users/mdigger")
	fmt.Printf("%v: %v\n", name, params.Get("name"))
	// Output: userName: mdigger
}
```

### Отмазка от каких-либо гарантий

Раньше это было составной частью библиотеки <github.com/mdigger/rest>, где,
по большей части, вся эта функциональность была просто скрыта и не доступна
для самостоятельного использования. Но она потребовалась мне для некоторых
внутренних проектов и я решил вынести ее в отдельную библиотеку. Я не
гарантирую, что библиотека не будет время от времени меняться под мои
собственные нужды, поэтому, если вы хотите ее использовать в своих проектах,
то лучшим способом является забрать ее себе целиком и дальше делать с ней
все, что хотите.
