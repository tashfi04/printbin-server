package api

import (
	"github.com/tashfi04/printbin-server/api/middlewares"
	"github.com/tashfi04/printbin-server/config"
	"github.com/tashfi04/printbin-server/conn"
	"github.com/tashfi04/printbin-server/data"
	"github.com/tashfi04/printbin-server/repos"
	"github.com/tashfi04/printbin-server/utils"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func listUserFiles(w http.ResponseWriter, r *http.Request) {

	userInfo, err := middlewares.GetUserInfo(r)
	if err != nil {
		utils.Logger().Errorln("Missing user context in header: ", err)
		rndr.JSON(w, http.StatusUnprocessableEntity, utils.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Missing user context in header",
			Error:      "Missing user context in header",
		})
		return
	}

	stringCurrentPage := r.URL.Query().Get("page")
	currentPage, err := strconv.Atoi(stringCurrentPage)
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid value for page query param",
			Error:      err.Error(),
		})
		return
	}

	stringPageLimit := r.URL.Query().Get("limit")
	pageLimit, err := strconv.Atoi(stringPageLimit)
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid value for limit query param",
			Error:      err.Error(),
		})
		return
	}

	if currentPage < 1 {
		currentPage = 1
	}
	if pageLimit < 5 {
		pageLimit = 5
	}

	searchParam := r.URL.Query().Get("search")
	searchParam = strings.TrimSpace(searchParam)
	if searchParam != "" {
		searchParam = "%" + searchParam + "%"
	}

	fileList, err := repos.FileRepo().ListUserFiles(conn.DB(), currentPage, pageLimit, searchParam, userInfo.UserID)
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database Query failed",
			Error:      err.Error(),
		})
		return
	}

	rndr.JSON(w, http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Data:       fileList,
	})
	return
}

func submitFile(w http.ResponseWriter, r *http.Request) {

	userInfo, err := middlewares.GetUserInfo(r)
	if err != nil {
		utils.Logger().Errorln("Missing user context in header: ", err)
		rndr.JSON(w, http.StatusUnprocessableEntity, utils.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Missing user context in header",
			Error:      "Missing user context in header",
		})
		return
	}

	user, err := data.User().GetUserByID(conn.DB(), userInfo.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rndr.JSON(w, http.StatusNotFound, utils.Response{
				StatusCode: http.StatusNotFound,
				Message:    "User not found",
				Error:      err.Error(),
			})
			return
		}
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database query failed",
			Error:      err.Error(),
		})
		return
	}

	stringFilePageCount := r.URL.Query().Get("pages")
	filePageCount, err := strconv.Atoi(stringFilePageCount)
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid value for pages query param",
			Error:      err.Error(),
		})
		return
	}
	if filePageCount < 1 {
		utils.Logger().Errorln("invalid value for pages query param")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid value for pages query param",
			Error:      "invalid value for pages query param",
		})
		return
	}

	userPrintPageCount := int(user.PrintPageCount) + filePageCount

	if userPrintPageCount > config.App().UserPrintLimit {
		utils.Logger().Errorln("print limit reached for user")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Print limit reached for user",
			Error:      "print limit reached for user",
		})
		return
	}

	trackingID := r.Header.Get("Tracking-ID")
	if trackingID == "" {
		utils.Logger().Errorln("Tracking-ID missing in form data")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Tracking-ID missing in form data",
		})
		return
	}

	err = r.ParseMultipartForm(8 << 20)
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to parse multipart form",
			Error:      err.Error(),
		})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to read file",
			Error:      err.Error(),
		})
		return
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".docx") {
		utils.Logger().Errorln("invalid file type")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid file type",
			Error:      "invalid file type",
		})
		return
	}

	fileHeader := make([]byte, 512)
	if _, err = file.Read(fileHeader); err != nil {
		if err != io.EOF {
			utils.Logger().Errorln(err)
			rndr.JSON(w, http.StatusBadRequest, utils.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Failed to read file header",
				Error:      err.Error(),
			})
			return
		}
	}

	if _, err = file.Seek(0, 0); err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to set seek file in start position",
			Error:      err.Error(),
		})
		return
	}

	//mime := http.DetectContentType(fileHeader)
	//if !strings.Contains(mime, "video/webm") { // && !strings.Contains(mime, "application/octet-stream") {
	//	utils.Logger().Errorln("invalid mime type")
	//	rndr.JSON(w, http.StatusBadRequest, utils.Response{
	//		StatusCode: http.StatusBadRequest,
	//		Message:    "Invalid mime type",
	//		Error:      "invalid mime type",
	//	})
	//	return
	//}

	directory := config.App().StoragePath
	if err = utils.CreateDirectory(directory); err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save file",
			Error:      err.Error(),
		})
	}

	dest := utils.BuildPath(trackingID + ".docx")
	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save file",
			Error:      err.Error(),
		})
	}
	defer func() {
		f.Close()
	}()

	_, err = io.Copy(f, file)
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save file",
			Error:      err.Error(),
		})
		return
	}

	err = repos.FileRepo().SubmitFile(conn.DB(), userInfo.UserID, uint(userPrintPageCount), trackingID)
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database Query failed",
			Error:      err.Error(),
		})
		return
	}

	rndr.JSON(w, http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Message:    "File submitted successfully",
		Data:       map[string]interface{}{"available_print_page_count": config.App().UserPrintLimit - userPrintPageCount},
	})
	return
}
