package swagger

import (
	docs "go-mongo-auth/docs"
	"go-mongo-auth/internal/config"
)

func ConfigureSwagger(appConfig config.App) {
	docs.SwaggerInfo.BasePath = appConfig.BasePath
	docs.SwaggerInfo.Title = appConfig.Name
	docs.SwaggerInfo.Version = appConfig.Version
	docs.SwaggerInfo.Description = appConfig.Description
	docs.SwaggerInfo.Schemes = []string{"http"}
}
