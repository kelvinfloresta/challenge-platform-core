package auth_case

import (
	"conformity-core/usecases/company_case"
	"conformity-core/utils"
)

func (a AuthCase) getWorkspaceByLogin(login string) (string, error) {
	domain := utils.GetEmailDomain(login)
	if domain == "" {
		return "", nil
	}

	company, err := a.companyCase.GetOneByFilter(company_case.GetOneByFilterInput{
		Domain: domain,
	})

	if err != nil {
		return "", err
	}

	if company == nil {
		return "", nil
	}

	if company.SSOEnabled {
		return company.Workspace, nil
	}

	return "", nil
}
