package department_controller

import (
	"conformity-core/context"
	"conformity-core/enums/user_company_enum"
	"conformity-core/usecases/department_case"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func PaginateUsers(c *fiber.Ctx) error {
	actualPage, err := strconv.Atoi(c.Params("actualPage"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(c.Params("pageSize"))
	if err != nil {
		return err
	}

	data := department_case.PaginateUsersInput{
		ActualPage:   actualPage,
		PageSize:     pageSize,
		Name:         c.Query("name"),
		DepartmentID: c.Query("departmentId"),
		Status:       c.Query("status"),
		Role:         user_company_enum.Role(c.Query("role")),
		Email:        c.Query("email"),
	}

	users, err := department_case.Singleton.PaginateUsers(context.New(c.Context()), data)

	if err != nil {
		return err
	}

	return c.JSON(users)
}
