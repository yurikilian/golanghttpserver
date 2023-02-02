package exception

//import (
//	"github.com/labstack/echo/v4"
//)
//
//type CustomErrorHandler struct {
//}
//
//func (h *CustomErrorHandler) Handle(err error, ctx echo.Context) {
//
//	//TODO: handle methods not allowed, not found and other statuses
//	problem, ok := err.(*Problem)
//	if !ok {
//		ctx.Logger().Error(err)
//		ctx.JSON(500, NewInternalServerError())
//	} else {
//		ctx.JSON(problem.Code, problem)
//	}
//}
//
//func NewCustomErrorHandler() *CustomErrorHandler {
//	return &CustomErrorHandler{}
//}
