package wrap

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GET[O any](e *echo.Echo, path string, fn func(echo.Context) (*O, error)) {
	e.GET(path, func(c echo.Context) error {
		out, err := fn(c)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, out)
	})
}

func POST[I, O any](e *echo.Echo, path string, fn func(echo.Context, *I) (*O, error)) {
	e.POST(path, func(c echo.Context) error {
		var in I
		if err := c.Bind(&in); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		out, err := fn(c, &in)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, out)
	})
}

func PUT[I, O any](e *echo.Echo, path string, fn func(echo.Context, *I) (*O, error)) {
	e.PUT(path, func(c echo.Context) error {
		var in I
		if err := c.Bind(&in); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		out, err := fn(c, &in)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, out)
	})
}

func DELETE[O any](e *echo.Echo, path string, fn func(echo.Context) (*O, error)) {
	e.DELETE(path, func(c echo.Context) error {
		out, err := fn(c)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, out)
	})
}

type Group struct {
	g *echo.Group
}

func NewGroup(e *echo.Echo, prefix string) *Group {
	return &Group{g: e.Group(prefix)}
}

func (grp *Group) G() *echo.Group {
	return grp.g
}

func GroupGET[O any](g *echo.Group, path string, fn func(echo.Context) (*O, error)) {
	g.GET(path, func(c echo.Context) error {
		out, err := fn(c)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, out)
	})
}

func GroupPOST[I, O any](g *echo.Group, path string, fn func(echo.Context, *I) (*O, error)) {
	g.POST(path, func(c echo.Context) error {
		var in I
		if err := c.Bind(&in); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		out, err := fn(c, &in)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, out)
	})
}

func GroupPUT[I, O any](g *echo.Group, path string, fn func(echo.Context, *I) (*O, error)) {
	g.PUT(path, func(c echo.Context) error {
		var in I
		if err := c.Bind(&in); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		out, err := fn(c, &in)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, out)
	})
}

func GroupDELETE[O any](g *echo.Group, path string, fn func(echo.Context) (*O, error)) {
	g.DELETE(path, func(c echo.Context) error {
		out, err := fn(c)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, out)
	})
}

/*
OneOf Example

	func registerExampleRoutes(e *echo.Echo) {
		POST[CreateUserRequest, OneOf[UserResponse, ValidationError]](e, "/users", createUser)
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
