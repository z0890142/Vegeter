package controller

import (
	"Vegeter/helper/Comman"
	"Vegeter/helper/DB"
	"Vegeter/model"
	"Vegeter/websocket"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var hub *websocket.Hub

func init() {
	hub = websocket.NewHub()
	go hub.Run()
}
func Register(w http.ResponseWriter, r *http.Request) {
	Log, _ := Comman.LogInit("Register", "Vegeter", logrus.DebugLevel)
	var response model.ApiResponse

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) //io.LimitReader限制大小
	defer r.Body.Close()
	var userData model.User
	err = json.Unmarshal(body, &userData)
	if err != nil {
		Comman.ResponseWithJson(w, http.StatusOK, response, Log)
		return
	}

	isExist, uuidString := DB.CheckUserEsixt(userData)
	if isExist {
		response.ResultCode = "200"
		response.ResultMessage = uuidString
		Comman.ResponseWithJson(w, http.StatusOK, response, Log)
		return
	}

	userData.Uuid = fmt.Sprintf("%s", uuid.NewV4())

	err = DB.InsertUser(userData)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Insert User Error")
		response.ResultCode = "400"
		response.ResultMessage = err.Error()
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}
	response.ResultCode = "200"
	response.ResultMessage = userData.Uuid
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
	Log, _ := Comman.LogInit("AddFriend", "Vegeter", logrus.DebugLevel)
	var response model.ApiResponse

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) //io.LimitReader限制大小
	defer r.Body.Close()
	var friendShip model.FriendShip
	err = json.Unmarshal(body, &friendShip)
	if err != nil {
		response.ResultCode = "400"
		response.ResultMessage = err.Error()
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}

	err = DB.InsertFriendShip(friendShip)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Insert User Error")
		response.ResultCode = "400"
		response.ResultMessage = err.Error()
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)

	}
	response.ResultCode = "200"
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)
}

func AddPrice(w http.ResponseWriter, r *http.Request) {
	Log, _ := Comman.LogInit("AddPrice", "Vegeter", logrus.DebugLevel)
	var response model.ApiResponse

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) //io.LimitReader限制大小
	defer r.Body.Close()
	var price model.Price
	err = json.Unmarshal(body, &price)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Unmarshal Error")
		response.ResultCode = "400"
		response.ResultMessage = err
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}

	oldPrice, err := DB.GetCurrentPrice(price.Uuid)
	if err != nil && err.Error() != "Get CurrentPrice : sql: no rows in result set" {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Get Current Price Error")
		response.ResultCode = "400"
		response.ResultMessage = err.Error()
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}
	if oldPrice != 0 {
		err = DB.UpdatePrice(price)
		if err != nil {
			Log.WithFields(logrus.Fields{
				"Error": err,
			}).Error("Update Price Error")
			response.ResultCode = "400"
			response.ResultMessage = err.Error()
			Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
			return
		}
		response.ResultCode = "200"
		Comman.ResponseWithJson(w, http.StatusOK, response, Log)
		return
	}

	err = DB.InsertPrice(price)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Insert Price Error")
		response.ResultCode = "400"
		response.ResultMessage = err.Error()
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}

	response.ResultCode = "200"
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)
}

func GetPriceRecord(w http.ResponseWriter, r *http.Request) {
	Log, _ := Comman.LogInit("GetPriceRecord", "Vegeter", logrus.DebugLevel)
	var response model.ApiResponse
	var uuid string
	vars := mux.Vars(r)
	uuid = vars["uuid"]
	priceList, err := DB.GetPriceRecord(uuid)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Get Current Price Error")
		response.ResultCode = "400"
		response.ResultMessage = err.Error()
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}
	response.ResultCode = "200"
	response.ResultMessage = priceList
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	Log, _ := Comman.LogInit("GetUserInfo", "Vegeter", logrus.DebugLevel)
	var response model.ApiResponse

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	userInfo := DB.GetUserInfo(uuid)
	response.ResultCode = "200"
	response.ResultMessage = userInfo
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)
}

func GetAddFriend(w http.ResponseWriter, r *http.Request) {
	Log, _ := Comman.LogInit("GetAddFriend", "Vegeter", logrus.DebugLevel)
	var response model.ApiResponse

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	friendList, err := DB.GetCheckFriend(uuid)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"uuid":  uuid,
			"Error": err,
		}).Error("GetCheckFriend error")
		response.ResultCode = "400"
		response.ResultMessage = err
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}

	response.ResultCode = "200"
	response.ResultMessage = friendList
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)

}

func ConfirmFriend(w http.ResponseWriter, r *http.Request) {
	Log, _ := Comman.LogInit("ConfirmFriend", "Vegeter", logrus.DebugLevel)
	var response model.ApiResponse

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) //io.LimitReader限制大小
	defer r.Body.Close()
	var friendShip model.FriendShip
	err = json.Unmarshal(body, &friendShip)
	if err != nil {
		response.ResultCode = "400"
		response.ResultMessage = err.Error()
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}

	err = DB.UpdateFriendShip(friendShip)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Insert User Error")
		response.ResultCode = "400"
		response.ResultMessage = err.Error()
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}
	response.ResultCode = "200"
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)
}

func GetFriend(w http.ResponseWriter, r *http.Request) {
	Log, _ := Comman.LogInit("GetFriend", "Vegeter", logrus.DebugLevel)
	var response model.ApiResponse

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	friendList, err := DB.GetFriend(uuid)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"uuid":  uuid,
			"Error": err,
		}).Error("GetFriend error")
		response.ResultCode = "400"
		response.ResultMessage = err
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}

	response.ResultCode = "200"
	response.ResultMessage = friendList
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)

}

func GetALLRecord(w http.ResponseWriter, r *http.Request) {
	Log, _ := Comman.LogInit("GetALLRecord", "Vegeter", logrus.DebugLevel)
	var response model.ApiResponse

	priceList, err := DB.GetAllPriceRecord()
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Get All Price Error")
		response.ResultCode = "400"
		response.ResultMessage = err
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}
	response.ResultCode = "200"
	response.ResultMessage = priceList
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)
}

func StartWebSocket(w http.ResponseWriter, r *http.Request) {
	websocket.ServeWs(hub, w, r)
}
