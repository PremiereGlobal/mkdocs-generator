module github.com/PremiereGlobal/mkdocs-generator

go 1.15

replace github.com/PremiereGlobal/mkdocs-generator => ./

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday v2.0.0+incompatible

require (
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	gopkg.in/russross/blackfriday.v2 v2.1.0
	gopkg.in/yaml.v2 v2.4.0
)
