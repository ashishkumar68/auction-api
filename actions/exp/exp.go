package exp

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func SampleUploadFiles(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, utils.GetFileSystemFilePath(fmt.Sprintf("sample-uploads-%d", time.Now().UnixNano())))

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
