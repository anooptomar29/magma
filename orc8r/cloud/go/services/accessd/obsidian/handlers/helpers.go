/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"github.com/labstack/echo"
	"net/http"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/obsidian/access"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/accessd/obsidian/models"
	"magma/orc8r/cloud/go/services/certifier"
)

func getOperatorForRead(c echo.Context) (*protos.Identity, *echo.HTTPError) {
	caller, err := access.RequestOperator(c)
	if err != nil {
		return nil, handlers.HttpError(err)
	}
	operator, httpErr := getOperator(c)
	if httpErr != nil {
		return nil, httpErr
	}
	if err := accessd.CheckReadPermission(caller, operator); err != nil {
		return nil, handlers.HttpError(err, http.StatusForbidden)
	}
	return operator, nil
}

func getOperatorForWrite(c echo.Context) (*protos.Identity, *echo.HTTPError) {
	caller, err := access.RequestOperator(c)
	if err != nil {
		return nil, handlers.HttpError(err)
	}
	operator, httpErr := getOperator(c)
	if httpErr != nil {
		return nil, httpErr
	}
	if err := accessd.CheckWritePermission(caller, operator); err != nil {
		return nil, handlers.HttpError(err, http.StatusForbidden)
	}
	return operator, nil
}

func getOperator(c echo.Context) (*protos.Identity, *echo.HTTPError) {
	operatorID, httpErr := handlers.GetOperatorId(c)
	if httpErr != nil {
		return nil, httpErr
	}
	return identity.NewOperator(operatorID), nil
}

func getNetwork(c echo.Context) (*protos.Identity, *echo.HTTPError) {
	networkID, httpErr := handlers.GetNetworkId(c)
	if httpErr != nil {
		return nil, httpErr
	}
	return identity.NewNetwork(networkID), nil
}

func getCertificateSNs(operator *protos.Identity) ([]models.CertificateSn, error) {
	certificates, err := certifier.FindCertificates(operator)
	if err != nil {
		return []models.CertificateSn{}, err
	}
	modelCertificates := make([]models.CertificateSn, len(certificates))
	for i, certificate := range certificates {
		modelCertificates[i] = models.CertificateSn(certificate)
	}
	return modelCertificates, nil
}
