package tests

import (
	"books/database"
	"books/router"
	"books/validator"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type DatabaseTestSuite struct {
	suite.Suite
	DB     *gorm.DB
	Router *gin.Engine
}

func (suite *DatabaseTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.DB = database.CreateDb().Begin()
	suite.Router = router.SetupRouter(suite.DB)
	validator.Init(suite.DB)
}

func (suite *DatabaseTestSuite) TearDownTest() {
	if suite.DB != nil {
		suite.DB.Rollback()
	}
}
