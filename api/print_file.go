package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/tashfi04/printbin-server/config"
	"github.com/tashfi04/printbin-server/conn"
	"github.com/tashfi04/printbin-server/models"
	"github.com/tashfi04/printbin-server/repos"
	"github.com/tashfi04/printbin-server/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func updateStatus(w http.ResponseWriter, r *http.Request) {

	trackingID := r.Header.Get("Tracking-ID")
	if trackingID == "" {
		utils.Logger().Errorln("Tracking-ID missing in header")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Tracking-ID missing in header",
			Error:      "Tracking-ID missing in header",
		})
		return
	}

	if err = repos.PrintFileRepo().UpdateStatus(conn.DB(), trackingID); err != nil {
		utils.Logger().Errorln(err)
		if err == gorm.ErrRecordNotFound {
			rndr.JSON(w, http.StatusOK, utils.Response{
				StatusCode: http.StatusOK,
				Message:    "File not found",
				Error:      err.Error(),
			})
			return
		}
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Database Query failed",
			Error:      err.Error(),
		})
		return
	}

	rndr.JSON(w, http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Message:    "File status updated",
	})
	return
}

func listFiles(w http.ResponseWriter, r *http.Request) {

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

	status := r.URL.Query().Get("status")
	if _, exists := models.FileStatusTypeToValue[models.FileStatusType(status)]; !exists {
		utils.Logger().Errorln("invalid status")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid status",
			Error:      "invalid status",
		})
		return
	}

	var roomList []string
	if err = json.NewDecoder(r.Body).Decode(&roomList); err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusUnprocessableEntity, utils.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed to parse rooms from request body",
			Error:      err.Error(),
		})
		return
	}

	for _, room := range roomList {
		if _, exists := config.App().RoomListMap[strings.TrimSpace(room)]; !exists {
			utils.Logger().Errorln("Invalid room_number")
			rndr.JSON(w, http.StatusBadRequest, utils.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid room_number",
				Error:      "invalid room_number",
			})
			return
		}
	}

	fileList, err := repos.PrintFileRepo().ListFiles(conn.DB(), currentPage, pageLimit, status, searchParam, roomList)
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

func serveFile(w http.ResponseWriter, r *http.Request) {

	ctx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
	//dir := path.Join(config.LocalStorage().BasePath, config.LocalStorage().AnnotationPath)
	dir := config.App().StoragePath
	fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(dir)))
	fs.ServeHTTP(w, r)
}

func listRooms(w http.ResponseWriter, r *http.Request) {

	rndr.JSON(w, http.StatusOK, utils.Response{
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"rooms": config.App().RoomList},
	})
	return
}
