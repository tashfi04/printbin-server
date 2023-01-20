package api

import (
	"encoding/csv"
	"fmt"
	"github.com/tashfi04/printbin-server/api/middlewares"
	"github.com/tashfi04/printbin-server/config"
	"github.com/tashfi04/printbin-server/conn"
	"github.com/tashfi04/printbin-server/models"
	"github.com/tashfi04/printbin-server/repos"
	"github.com/tashfi04/printbin-server/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func uploadUser(w http.ResponseWriter, r *http.Request) {

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

	if userInfo.Role > 1 {
		utils.Logger().Infoln("Unauthorized access prohibited")
		rndr.JSON(w, http.StatusUnauthorized, utils.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized access prohibited",
			Error:      "Unauthorized access prohibited",
		})
		return
	}

	var maxFileSize int64 = 10 << 20
	if err = r.ParseMultipartForm(maxFileSize); err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusInternalServerError, utils.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to parse multipart form",
			Error:      err.Error(),
		})
		return
	}

	file, header, err := r.FormFile("users")
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

	if !strings.HasSuffix(header.Filename, ".csv") {
		utils.Logger().Errorln("invalid CSV file")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid CSV file",
			Error:      "invalid CSV file",
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

	mime := http.DetectContentType(fileHeader)
	if !strings.Contains(mime, "text/plain") && !strings.Contains(mime, "application/octet-stream") {
		utils.Logger().Errorln("invalid mime type")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid mime type",
			Error:      "invalid mime type",
		})
		return
	}

	if header.Size > maxFileSize {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("File size exceeded. Max value is %d", maxFileSize),
			Error:      "File size too large",
		})
		return
	}

	csvReader := csv.NewReader(file)

	columnHeaders, err := csvReader.Read()
	if err != nil {
		utils.Logger().Errorln(err)
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid data",
			Error:      err.Error(),
		})
		return
	}

	if strings.TrimSpace(columnHeaders[0]) != "username" || strings.TrimSpace(columnHeaders[1]) != "password" || strings.TrimSpace(columnHeaders[2]) != "team_name" || strings.TrimSpace(columnHeaders[3]) != "room_number" || strings.TrimSpace(columnHeaders[4]) != "role" {
		utils.Logger().Errorln("invalid column header")
		rndr.JSON(w, http.StatusBadRequest, utils.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid column header",
			Error:      "invalid column header",
		})
		return
	}

	var users []models.User

	validRecordColumns := 5
	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if len(record) != validRecordColumns {
			utils.Logger().Errorln("Invalid number of columns")
			rndr.JSON(w, http.StatusBadRequest, utils.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid number of columns",
				Error:      "invalid data",
			})
			return
		}

		username := record[0]
		password := record[1]
		teamName := record[2]
		roomNumber := record[3]
		stringRole := record[4]

		if strings.TrimSpace(username) == "" {
			utils.Logger().Errorln("Invalid username")
			rndr.JSON(w, http.StatusBadRequest, utils.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid username",
				Error:      "invalid username",
			})
			return
		}

		if strings.TrimSpace(password) == "" {
			utils.Logger().Errorln("Invalid password")
			rndr.JSON(w, http.StatusBadRequest, utils.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid password",
				Error:      "invalid password",
			})
			return
		}

		role, err := strconv.Atoi(stringRole)
		if err != nil {
			utils.Logger().Errorln(err)
			rndr.JSON(w, http.StatusBadRequest, utils.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid role",
				Error:      "invalid role",
			})
			return
		}
		if role < 0 || role > 2 {
			utils.Logger().Errorln("invalid role")
			rndr.JSON(w, http.StatusBadRequest, utils.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid role",
				Error:      "invalid role",
			})
			return
		}

		teamName = strings.TrimSpace(teamName)
		roomNumber = strings.TrimSpace(roomNumber)
		if role == 0 {
			if teamName == "" {
				utils.Logger().Errorln("Invalid team_name")
				rndr.JSON(w, http.StatusBadRequest, utils.Response{
					StatusCode: http.StatusBadRequest,
					Message:    "Invalid team_name",
					Error:      "invalid team_name",
				})
				return
			}

			if _, exists := config.App().RoomListMap[roomNumber]; !exists {
				utils.Logger().Errorln("Invalid room_number")
				rndr.JSON(w, http.StatusBadRequest, utils.Response{
					StatusCode: http.StatusBadRequest,
					Message:    "Invalid room_number",
					Error:      "invalid room_number",
				})
				return
			}
		}

		encryptedPassword, err := utils.Encrypt(password)
		if err != nil {
			utils.Logger().Errorln(err)
			rndr.JSON(w, http.StatusInternalServerError, utils.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to create users",
				Error:      err.Error(),
			})
			return
		}

		user := models.User{
			Username:   username,
			Password:   encryptedPassword,
			TeamName:   teamName,
			RoomNumber: roomNumber,
			Role:       uint(role),
		}
		users = append(users, user)
	}

	err = repos.AdminRepo().CreateUser(conn.DB(), users)
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
		Message:    "Users successfully created",
	})
	return
}
