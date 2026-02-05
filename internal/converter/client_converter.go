package converter

import (
	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
)

func ModelReqRegisterToEntityClient(req *model.ReqRegister, client *entity.Client) {
	client.Name = req.Name
	client.Email = req.Email
	client.Password = req.Password
}

func EntityClientToModelResRegister(client *entity.Client, res *model.ResRegister) error {
	res.Name = client.Name
	email, err := client.DecryptEmail()
	if err != nil {
		return errkit.AddFuncName(err)
	}
	res.Email = email
	return nil
}

func EntityClientToModelResCurrent(client *entity.Client, res *model.ResGetClientDetail) error {
	res.ID = client.ID
	res.Name = client.Name
	email, err := client.DecryptEmail()
	if err != nil {
		return errkit.AddFuncName(err)
	}
	res.Email = email
	apiKey, err := client.DecryptAPIKey()
	if err != nil {
		return errkit.AddFuncName(err)
	}
	res.APIKey = apiKey
	res.CreatedAt = client.CreatedAt
	res.UpdatedAt = client.UpdatedAt
	return nil
}

func EntityClientToModelClientAuth(client *entity.Client, clientAuth *model.ClientAuth) error {
	clientAuth.ID = client.ID
	clientAuth.Name = client.Name
	clientAuth.Email = client.Email
	email, err := client.DecryptEmail()
	if err != nil {
		return errkit.AddFuncName(err)
	}
	clientAuth.EmailDecrypted = email
	clientAuth.APIKey = client.APIKey
	apiKey, err := client.DecryptAPIKey()
	if err != nil {
		return errkit.AddFuncName(err)
	}
	clientAuth.APIKeyDecrypted = apiKey
	clientAuth.CreatedAt = client.CreatedAt
	clientAuth.UpdatedAt = client.UpdatedAt
	return nil
}

func ModelClientAuthToModelResGetClientDetail(clientAuth *model.ClientAuth, res *model.ResGetClientDetail) {
	res.ID = clientAuth.ID
	res.Name = clientAuth.Name
	res.Email = clientAuth.EmailDecrypted
	res.APIKey = clientAuth.APIKeyDecrypted
	res.CreatedAt = clientAuth.CreatedAt
	res.UpdatedAt = clientAuth.UpdatedAt
}
