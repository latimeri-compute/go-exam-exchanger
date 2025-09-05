package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/constraints"
)

// структура для упрощения синтаксиса
type JSONEnveloper map[string]any

type Currency float64

// обрабатывает полученный json
func UnpackJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	// установка максимального размер данных
	// 1 мегабайт
	maxBytes := 1_048_576
	// разбор тела реквеста
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(destination)
	if err != nil {
		// обработка ошибок
		var (
			// ошибка синтаксиса, невозможно разобрать объект
			syntaxError *json.SyntaxError
			// неправильный тип JSON
			unmarshalTypeError *json.UnmarshalTypeError
			// неверный аргумент, неподдерживающийся методом Decode()
			invalidUnmarshalError *json.InvalidUnmarshalError
			// тело запроса слишком большое
			maxBytesError *http.MaxBytesError
		)

		switch {
		// тело слишком большое
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("превышен размер запроса. Тело не может быть больше %d байтов", maxBytesError.Limit)

		// ошибка синтаксиса
		case errors.As(err, &syntaxError):
			return fmt.Errorf("поле %q содержит неправильно составленный JSON", syntaxError.Offset)

		// неправильный тип JSON
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("поле %q содержит недопустимый примитив JSON", unmarshalTypeError.Field)
			}
			return fmt.Errorf("недопустимый примитив JSON на символе %d", unmarshalTypeError.Offset)

		// пустое тело
		case errors.Is(err, io.EOF):
			return fmt.Errorf("тело запроса не должено быть пустым")

		// внезапный конец строки
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("тело запроса содержит неправильно составленный JSON")

		// неизвестное поле
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			unknownField := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("тело содержит неизвестное поле %s", unknownField)

		// в данном случае лучше всего запаниковать
		// тк ошибка не со стороны клиента
		// и в принципе не понятно а что а куда
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// Возвращает ошибку io.EOF если тело содержало только один объект JSON
	// При возврате любого другого значения -- было передано более одного значения
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("тело должно содержать только один объект JSON")
	}

	return nil
}

func Abs[T constraints.Signed](num T) T {
	if num < 0 {
		num = -num
	}
	return num
}

func (c *Currency) MarshalJSON() ([]byte, error) {
	f := float64(*c)
	res := fmt.Appendf(nil, "%.2f", f)
	return res, nil
}
