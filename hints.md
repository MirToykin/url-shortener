### GOROOT
через GOROOT можно настраивать разные версии Go (Preferences->Go->GOROOT)
### idle_timeout
Сервер обычно оптимизирован таким образмом, что если от одного клиента приходит несколько 
запросов с небольшим интервалом, то поскольку нет смысла на каждый запрос от одного клиента
открывать и закрывать соединение оно держится открытым какое-то время. Но если клиент долго 
не шлет запросы, то такое соединение закрывается. idle_timeout и есть время ожидания повторного
запроса до закрытия соединения.

### Must префикс в имени функции
По соглашению Must используется в именах функций которые паникуют, т.е. не возвращают ошибку.
Так стоит делать в редких случаях, когда при ошибке не возможна дальнейшая работа приложения.

### slog
slog - это не конкретный логгер, а скорее некая обертка для логгера. Это библиотека для работы
с логгерами. Он имеет несколько дефолтных логгеров - это текстовый и json. Текстовый для просмотра
информации локально при запуске, а json для отправки на сервер, где они будут обработаны сборщиком/агрегатором
логов. Но в него можно интегрировать и другой логгер, например, logrus.

**Дополнительные парметры могут быть использованы как для одного лога, так и для всех:**
- log.Info("starting url shortener", slog.String("env", cfg.Env))
- log = log.With(slog.String("env", cfg.Env))

### Структура проекта
- **internal/lib** - каталог для различных вспомогательных функций для разных компонентов проекта
- в целом короткие непонятные имена пакетов считаются плохой практикой, но если в этом пакете разме
щен функционал, который используется очень часто, то этим можно пренебречь. Например, **internal/lib/logger/sl**, пакет sl
будет использоваться при каждом логировании ошибки (sl - slog)

### Особенности языка
- если нам нужен импорт, который нигде явно не используется, то для него можно использовать
заглушку _ , например **_ "github.com/mattn/go-sqlite3"**