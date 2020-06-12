package handler

import goke "github.com/remorses/goke/src"

func isAdminSecretValid(admins []goke.AdminConfig, secret string) bool {
	for _, config := range admins {
		if config.Secret == secret {
			return true
		}
	}
	return false
}
