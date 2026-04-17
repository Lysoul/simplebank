package shared

import (
	"assignments/simplebank/adapters/monitoring"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"

	// "github.com/iancoleman/strcase"
	"go.uber.org/zap"
)

var (
	ErrNotFound          = ConstError("not_found")
	ErrPermissionDenied  = ConstError("permission_denied")
	ErrInvalidParameters = ConstError("invalid_parameters")
	ErrUnexpected        = ConstError("unexpected")
)

type ConstError string

func (err ConstError) Error() string {
	return string(err)
}
func (err ConstError) Is(target error) bool {
	ts := target.Error()
	es := string(err)
	return ts == es || strings.HasPrefix(ts, es+": ")
}
func (err ConstError) Wrap(inner error) error {
	return wrapError{msg: string(err), err: inner}
}

type wrapError struct {
	err error
	msg string
}

func (err wrapError) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v", err.msg, err.err)
	}
	return err.msg
}
func (err wrapError) Unwrap() error {
	return err.err
}
func (err wrapError) Is(target error) bool {
	return ConstError(err.msg).Is(target)
}

func ErrorToHTTP(ctx *gin.Context, err error) {
	validationErrs := &validator.ValidationErrors{}
	switch {
	case err == nil:
		return
	case errors.As(err, validationErrs):
		// e := err.(*validator.ValidationErrors)
		// todo mae better details
		ctx.AbortWithStatusJSON(400, gin.H{
			"error":   "validation_failed",
			"details": MapValidationErrors(validationErrs),
		})
	case errors.Is(err, ErrNotFound):
		ctx.AbortWithStatus(404)
	case errors.Is(err, ErrPermissionDenied):
		ctx.AbortWithStatusJSON(403, gin.H{
			"error": "permission_denied",
		})
	case errors.Is(err, ErrUnexpected):
		monitoring.Logger().Ctx(ctx).Error("unexpected error", zap.Error(err))
		ctx.AbortWithStatusJSON(500, gin.H{
			"error": "unexpected_error",
		})
	default:
		monitoring.Logger().Ctx(ctx).Error("no http mapping for error",
			zap.String("type", fmt.Sprintf("%T", err)),
			zap.Error(err))

		ctx.AbortWithStatus(500)
	}
	// todo: only handle errors we know and call shared.handleError(err) for the rest
}

// todo investigate using json schema for validation
// https://github.com/santhosh-tekuri/jsonschema
func MapValidationErrors(errs *validator.ValidationErrors) []string {
	details := make([]string, len(*errs))
	for errNum, err := range *errs {
		path := strings.Split(err.Namespace(), ".")
		field := make([]string, len(path)-1)
		for i, p := range path[1:] {
			field[i] += strcase.ToLowerCamel(p)
		}
		details[errNum] = strings.Join(field, ".") + " " + translateTag(err)
	}

	return details
}

func translateTag(err validator.FieldError) string {
	switch err.ActualTag() {
	case "gtfield":
		return "must greater than " + strcase.ToLowerCamel(err.Param())
	case "gtefield":
		return "must be greater than or equal " + strcase.ToLowerCamel(err.Param())
	case "lte":
		return "must be maximum " + err.Param()
	case "gt":
		return "must be greater than " + err.Param()
	case "gte":
		return "must be greater than or equal " + err.Param()
	case "ltefield":
		return "must be less then or equal " + strcase.ToLowerCamel(err.Param())
	case "required_with":
		return "required with " + strcase.ToLowerCamel(err.Param())
	case "oneof":
		return "must be one of " + err.Param()
	case "max":
		return "cannot exceed " + err.Param() + " characters"
	case "filetype":
		return "must be a valid " + err.Param() + " file"
	case "alphanum":
		return "must contain only alphanumeric characters"
	case "url":
		return "must be a valid URL"
	default:
		return err.ActualTag()
	}
}
