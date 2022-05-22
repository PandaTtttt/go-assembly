module github.com/PandaTtttt/go-assembly

require (
	github.com/Shopify/sarama v1.26.4
	github.com/gin-gonic/gin v1.7.0
	github.com/go-redis/redis/v7 v7.4.0
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/stretchr/testify v1.5.1
	go.uber.org/zap v1.15.0
	gorm.io/driver/mysql v1.0.1
	gorm.io/gorm v1.20.0
)

replace (
	gorm.io/driver/mysql v1.0.1 => github.com/go-gorm/mysql v1.0.1
	gorm.io/gorm v1.20.0 => github.com/go-gorm/gorm v1.20.0
)

go 1.14
