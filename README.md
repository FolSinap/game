### 1ая часть

Творческое задание - простая текстовая игра.
Играть в неё можно вводя команды в консоли.
Вводим команду - получаем ответ.
Такие игры делали когда компьютеры были большими, а интернета не было вовсе.
В следующих домашних заданиях мы будем увеличивать функционал игры и переводить её в интернет, делать многопользовательской.

Игровой мир обычно состоит из комнат, где может происходить какое-то действие.
Так же у нас есть игрок.
Как у игрока, так и у команты есть состояние.
initGame делает нового игрока и задаёт ему начальное состояние.
В данной версии можно обойтись глобальными переменными для игрока.

Команда парсится как
`$команда $параметр1 $параметр2 $параметр3`

В тестах представлены последовательности команд и получаемый ответ.
Задача - пройти все тесты и сделать правильно.
Под правильным понимается универсально, чтобы можно было без проблем что-то добавить или убрать.
Т.е. бесконечный набор захардкоженных if'ов не подойдёт для всего мира не подойдёт.
Услвоия есть толкьо внутри конкретной комнаты.
Надо думать в сторону объектов, вызова функций, структур, которые описывают состояние комнаты и игрока, функций которые описывают какой-то интерактив в комнате. Не забывайте что вы можете создать мапу из функций. Или можно реализовать триггер (действие, выполняемое при каком-то событии). Или у структуры поле может иметь тип "функция".

Тестовых кейсов много. Прочитайте их внимательно, там есть результаты работы всего что вам надо.
Задание сложное, не на 5 минут
Не стесняйтесь. Не забывайте, что у нас ест ьканал в телеграме - лучше спрашивать туда.
Однако преждем чем спрашивать - попробуйте что-то сделать и четко сформулируйте, что у вас не получается.

### 2ая часть

дальнейшее развитие game-0
добавляем взаимодействие нескольких игроков
раньше был 1 игрок и игра была по сути синхронна
теперь добавляется второй ( и третий, четвёртый ) и они оба могут "вводить" команды
значит у нас появляется асинхронное взаимодействие
надо использовать каналы и горутины для обработки команд игроков в каждой комнате
Смотрите примеры с мультипрексорами для организации ввода

Большая часть команд так же относится к одному игроку, но появляется несколько новых команд:
"сказать" - транслируется всем другим игрокам в этой комнате
"сказать_игроку" - отправляет персональное сообщение другому игроку, остальные его не слышат
Если выполнить команду "осмотреться", что можно увидеть другого игрока
Это значит, что у комнаты появляется список игроков в ней, которым надо передать сообщение

Чтобы как-то отличать игроков друг от друга, теперь при инициализации надо задать имя ( функция NewPlayer )
Естественно, количество игроков не должно быть захардкожено и может меняться

У игрока появляется метод HandleInput, который принимает от него входящее сообщение.
А так же метод HandleOutput, в который мы пишем то что игро должен видеть у себя на экране.
И метод GetOutput, который возвращает канал, принимающий сообщения

### 3ая часть

Продолжаем развивать игру из ДЗ-2 ит ДЗ-3
В этот раз выкладываем её на хероку в виде бота для телеграма

Можете начать с выкалдывания бота из примеров ( 4/06_bot ) и как убедитесь что всё получится - уже выкладывать остальное

Необходимо реализовать новые возможности:
* Если игрок ещё не записан в игре - он добавляется на кухню
* Если игрок не активен в течении 15 минут (не было сообщений-команд), то он удаляется из игры, если у него был инвентарь - он возвращается туда где был.
* У админа бота есть команда сброса игры - выкинуть всех игроков, вернуть предменты и замок в первоначальное состояние. Другим игрокам эта команда недоступна
* Бот расположен не в корне сервера, как в примере (см WebhookURL, ListenForWebhook), а в подпапке /bot , это надо чтобы мы далее могли в корне ещё что-то запускать другое в следующих ДЗ

В телеграме есть несколько вариантво ввода команд: настраиваемые кнопки на клавиатуре, инлайн-кнопки ( прямо в чате ), просто команды в чат. Реализовать можно любым способом.

### 4ая часть

Продолжаем оработки по игре

Надо добавить сохранение в базу следующей информации об игре:
* От кого пришла команда
* Во сколько (время)
* Какой результат ( json где в качестве ключа игрок, в качестве значения - полученное сообщение)

После этого сделать форму логина ( юзернейм любой, пароль из конфига )
При успешной авторизации создаётся сессия ( использовать встроенный кеш внутри бинарника. как в примере с сессией ), в сессию кладётся введённый юзернейм и его юзерагент. при перезагрзке бинарника ( обновления кода ) сессии очищаются. Сессия живёт 15 минут.

При наличии корректной сессии мы выводим этому пользователю
"Welcome, %username% from %useragent%
После этого выводятся 100 последних записей из базы

База для сообщений:
```
CREATE TABLE "commands" (
  "id" serial NOT NULL,
  "from" text NOT NULL,
  "command" text NOT NULL,
  "result" text NOT NULL,
  "time" integer NOT NULL
);
```
