package common

import (
	"context"
	"github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/client"
	"os"
)

type Env struct {
	SMTP_PORT      string
	STAGE          string
	GCP_PROJECT_ID string
	SMTP_HOST      string

	// names for retrieving secrets
	GMAIL_SECRET_NAME    string
	GMAIL_PASS_NAME      string
	MDB_USER_SECRET_NAME string
	MDB_PASS_SECRET_NAME string
	DB_SECRET_NAME       string

	// actual secrets
	Secrets
}

type Secrets struct {
	GMAIL_USER string
	GMAIL_PASS string
	MDB_USER   string
	MDB_PASS   string
	DB_NAME    string
}

func LoadEnv() *Env {
	return &Env{
		SMTP_PORT:            os.Getenv("PORT"),
		STAGE:                os.Getenv("ENVIRONMENT"),
		GCP_PROJECT_ID:       os.Getenv("PROJECT_ID"),
		GMAIL_SECRET_NAME:    os.Getenv("GMAIL_USER_SECRET_NAME"),
		GMAIL_PASS_NAME:      os.Getenv("GMAIL_PASS_SECRET_NAME"),
		MDB_PASS_SECRET_NAME: os.Getenv("MDB_PASS_SECRET_NAME"),
		MDB_USER_SECRET_NAME: os.Getenv("MDB_USER_SECRET_NAME"),
		DB_SECRET_NAME:       os.Getenv("DB_SECRET_NAME"),
		SMTP_HOST:            os.Getenv("SMTP_HOST"),
	}
}

func LoadSecrets(ctx context.Context, sm client.ISecretsManagerClient, env *Env) error {
	var err error

	env.GMAIL_USER, err = client.RetrieveSecret(ctx, sm, env.GMAIL_SECRET_NAME, env.GCP_PROJECT_ID)
	if err != nil {
		return err
	}

	env.GMAIL_PASS, err = client.RetrieveSecret(ctx, sm, env.GMAIL_PASS_NAME, env.GCP_PROJECT_ID)
	if err != nil {
		return err
	}

	env.DB_NAME, err = client.RetrieveSecret(ctx, sm, env.DB_SECRET_NAME, env.GCP_PROJECT_ID)
	if err != nil {
		return err
	}

	env.MDB_USER, err = client.RetrieveSecret(ctx, sm, env.MDB_USER_SECRET_NAME, env.GCP_PROJECT_ID)
	if err != nil {
		return err
	}

	env.MDB_PASS, err = client.RetrieveSecret(ctx, sm, env.MDB_PASS_SECRET_NAME, env.GCP_PROJECT_ID)
	if err != nil {
		return err
	}

	return nil
}
