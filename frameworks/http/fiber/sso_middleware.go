package app

import (
	"conformity-core/config"
	core_errors "conformity-core/errors"
	"conformity-core/usecases/auth_case"
	"conformity-core/usecases/company_case"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"strings"

	"fmt"

	"github.com/crewjam/saml/samlsp"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func SingleSignOnSAML(app *fiber.App) func(c *fiber.Ctx) error {
	keyPair := parseCerticate()

	rootURL, err := url.Parse(config.APIHost)
	if err != nil {
		panic(err)
	}

	return func(c *fiber.Ctx) error {
		workspace := c.Params("workspace")

		companyFound, err := company_case.Singleton.GetOneByFilter(company_case.GetOneByFilterInput{
			Workspace: workspace,
		})
		if err != nil {
			return err
		}

		if companyFound == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if companyFound.IdentityProviderMetadata == "" {
			return c.Redirect(config.AppURL+"/error/MISSING_IDP", http.StatusSeeOther)
		}

		idpMetadata, err := samlsp.ParseMetadata([]byte(companyFound.IdentityProviderMetadata))

		if err != nil {
			return err
		}

		path := c.Path()
		if strings.HasSuffix(path, "/check") {
			return c.SendStatus(fiber.StatusOK)
		}

		samlSP, _ := samlsp.New(samlsp.Options{
			EntityID:    config.AppURL,
			URL:         *rootURL,
			Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
			Certificate: keyPair.Leaf,
			IDPMetadata: idpMetadata,
		})

		samlSP.ServiceProvider.AcsURL.Path = "/saml/" + workspace + "/acs"

		if strings.HasSuffix(path, "/auth") {
			return adaptor.HTTPHandler(samlSP.RequireAccount(http.HandlerFunc(samlCallback)))(c)
		}

		if strings.HasSuffix(path, "/acs") {
			return adaptor.HTTPHandlerFunc(samlSP.ServeACS)(c)
		}

		if strings.HasSuffix(path, "/metadata") {
			return adaptor.HTTPHandlerFunc(samlSP.ServeMetadata)(c)
		}

		return c.SendStatus(fiber.StatusNotFound)
	}
}

func samlCallback(w http.ResponseWriter, r *http.Request) {
	s := samlsp.SessionFromContext(r.Context())
	if s == nil {
		return
	}

	sa, ok := s.(samlsp.SessionWithAttributes)
	if !ok {
		return
	}

	workspace := strings.TrimLeft(strings.TrimRight(r.RequestURI, "/auth"), "/saml")
	attr := sa.GetAttributes()

	input := auth_case.LoginSSOInput{
		OID:            attr.Get("objectGUID"),
		Email:          attr.Get("mail"),
		Name:           attr.Get("displayName"),
		DepartmentName: attr.Get("department"),
		Workspace:      workspace,
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
		Path:   "/",
	})

	token, err := auth_case.Singleton.LoginSSO(input)

	if err != nil {
		handleSamlError(w, r, err, input)
		return
	}

	if token == "" {
		http.Redirect(w, r, config.AppURL+"/forbidden", http.StatusSeeOther)
		return
	}

	newUrl := fmt.Sprint(config.AppURL, "/token/", token)
	http.Redirect(w, r, newUrl, http.StatusSeeOther)
}

func parseCerticate() tls.Certificate {
	keyPair, err := tls.LoadX509KeyPair("conformity-core.cert", "conformity-core.key")
	if err != nil {
		panic(err)
	}

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err)
	}

	return keyPair
}

func handleSamlError(w http.ResponseWriter, r *http.Request, err error, input auth_case.LoginSSOInput) {
	switch err.(type) {
	case core_errors.Unauthorized:
		http.Redirect(w, r, config.AppURL+"/error/"+err.Error(), http.StatusSeeOther)

	case core_errors.Forbidden:
		http.Redirect(w, r, config.AppURL+"/forbidden", http.StatusSeeOther)

	case core_errors.NotFound:
		http.Redirect(w, r, config.AppURL+"/error/"+err.Error(), http.StatusSeeOther)

	case core_errors.BadRequest:
		http.Redirect(w, r, config.AppURL+"/error/"+err.Error(), http.StatusSeeOther)

	case core_errors.Conflict:
		http.Redirect(w, r, config.AppURL+"/error/"+err.Error(), http.StatusSeeOther)

	default:
		hub := sentry.CurrentHub().Clone()
		scope := hub.Scope()
		scope.SetUser(sentry.User{
			Email:    input.Email,
			ID:       input.OID,
			Username: input.Name,
		})
		scope.SetRequest(r)
		hub.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
