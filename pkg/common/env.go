package common

import (
	"context"
	"github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/client"
)

type Env struct {
	Port            string
	Stage           string
	GcpProjectId    string
	MediaBucketName string
	SmtpHost        string

	// names for retrieving secrets
	GmailUserSecretName string
	GmailPassSecretName string
	MdbUserSecretName   string
	MdbPassSecretName   string
	DBNameSecretName    string

	// actual secrets
	Secrets
}

type Secrets struct {
	GmailUser string
	GmailPass string
	MdbUser   string
	MdbPass   string
	DBName    string
}

func LoadSecrets(ctx context.Context, sm client.ISecretsManagerClient, env *Env) error {
	var err error

	env.GmailUser, err = client.RetrieveSecret(ctx, sm, env.GmailUserSecretName, env.GcpProjectId)
	if err != nil {
		return err
	}

	env.GmailPass, err = client.RetrieveSecret(ctx, sm, env.GmailPassSecretName, env.GcpProjectId)
	if err != nil {
		return err
	}

	env.DBName, err = client.RetrieveSecret(ctx, sm, env.DBNameSecretName, env.GcpProjectId)
	if err != nil {
		return err
	}

	env.MdbUser, err = client.RetrieveSecret(ctx, sm, env.MdbUserSecretName, env.GcpProjectId)
	if err != nil {
		return err
	}

	env.MdbPass, err = client.RetrieveSecret(ctx, sm, env.MdbPassSecretName, env.GcpProjectId)
	if err != nil {
		return err
	}

	return nil
}
