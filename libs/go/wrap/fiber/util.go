package wrap

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

func GET[O any](a *fiber.App, path string, fn func(fiber.Ctx) (*O, error)) {
	a.Get(path, func(c fiber.Ctx) error {
		out, err := fn(c)
		if err != nil {
			return err
		}
		return c.JSON(out)
	})
}

func POST[I, O any](a *fiber.App, path string, fn func(fiber.Ctx, *I) (*O, error)) {
	a.Post(path, func(c fiber.Ctx) error {
		var in I
		if err := c.Bind().Body(&in); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		out, err := fn(c, &in)
		if err != nil {
			return err
		}
		return c.JSON(out)
	})
}

func PUT[I, O any](a *fiber.App, path string, fn func(fiber.Ctx, *I) (*O, error)) {
	a.Put(path, func(c fiber.Ctx) error {
		var in I
		if err := c.Bind().Body(&in); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		out, err := fn(c, &in)
		if err != nil {
			return err
		}
		return c.JSON(out)
	})
}

func DELETE[O any](a *fiber.App, path string, fn func(fiber.Ctx) (*O, error)) {
	a.Delete(path, func(c fiber.Ctx) error {
		out, err := fn(c)
		if err != nil {
			return err
		}
		return c.JSON(out)
	})
}

/*
OneOf Example

	func registerExampleRoutes(a *fiber.App) {
		POST[CreateUserRequest, OneOf[UserResponse, ValidationError]](a, "/users", createUser)
	}
*/
type OneOf[A, B any] struct {
	a   *A
	b   *B
	isA bool
}

func First[A, B any](val A) OneOf[A, B] {
	return OneOf[A, B]{a: &val, isA: true}
}

func Second[A, B any](val B) OneOf[A, B] {
	return OneOf[A, B]{b: &val, isA: false}
}

func (o OneOf[A, B]) IsFirst() bool  { return o.isA }
func (o OneOf[A, B]) IsSecond() bool { return !o.isA }

func (o OneOf[A, B]) First() (A, bool) {
	if o.isA && o.a != nil {
		return *o.a, true
	}
	var zero A
	return zero, false
}

func (o OneOf[A, B]) Second() (B, bool) {
	if !o.isA && o.b != nil {
		return *o.b, true
	}
	var zero B
	return zero, false
}

func (o OneOf[A, B]) Match(onA func(A), onB func(B)) {
	if o.isA && o.a != nil {
		onA(*o.a)
	} else if o.b != nil {
		onB(*o.b)
	}
}

func (o OneOf[A, B]) MarshalJSON() ([]byte, error) {
	if o.isA && o.a != nil {
		return json.Marshal(o.a)
	}
	return json.Marshal(o.b)
}
