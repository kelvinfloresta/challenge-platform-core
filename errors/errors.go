package core_errors

var ErrWithoutPrivileges Conflict = "USER_WITHOUT_PRIVILEGES"
var ErrUserWithoutDepartment NotFound = "USER_WITHOUT_DEPARTMENT"
var ErrCompanyNotFound NotFound = "COMPANY_NOT_FOUND"
var ErrDeleteDepartmentWithUsers Conflict = "DELETE_DEPARTMENT_WITH_USERS"

type Conflict string

func (e Conflict) Error() string {
	return string(e)
}

type Forbidden string

func (e Forbidden) Error() string {
	return string(e)
}

type NotFound string

func (e NotFound) Error() string {
	return string(e)
}

type Unauthorized string

func (e Unauthorized) Error() string {
	return string(e)
}

type BadRequest string

func (e BadRequest) Error() string {
	return string(e)
}
