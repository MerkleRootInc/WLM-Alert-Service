package common

import "os"

var (
	PORT                         = os.Getenv("PORT")
	ENVIRONMENT                  = os.Getenv("ENVIRONMENT")
	ALCHEMY_ENDPOINT_SECRET_NAME = os.Getenv("ALCHEMY_ENDPOINT_SECRET_NAME")
	GCP_PROJECT_ID               = os.Getenv("GCP_PROJECT_ID")
	MEDIA_BUCKET_NAME            = os.Getenv("MEDIA_BUCKET_NAME")
	MONGO_DB_USER                = os.Getenv("MONGO_DB_USER")
	MONGO_DB_PASS                = os.Getenv("MONGO_DB_PASS")
	DB_NAME                      = os.Getenv("DB_NAME")
)
