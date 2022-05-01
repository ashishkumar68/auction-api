package utils

import (
	"github.com/stretchr/testify/assert"
	"strings"
)

func (suite *UtilsTestSuite) TestRenameFile() {
	fileName := "test_file_name.jpeg"
	fileNameInfo := strings.Split(fileName, ".")
	beforeEditExtension := fileNameInfo[len(fileNameInfo)-1]

	newFileName, err := GetRenamedFileName(fileName)
	fileNameInfo = strings.Split(newFileName, ".")
	afterEditExtension := fileNameInfo[len(fileNameInfo)-1]

	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), newFileName, fileName)
	assert.Equal(suite.T(), beforeEditExtension, afterEditExtension)
}
