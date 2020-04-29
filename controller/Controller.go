package controller

import (
	"Vegeter/helper/Comman"
	"Vegeter/helper/DB"
	"Vegeter/model"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

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
	uuid := uuid.NewV4()
	userData.Uuid = fmt.Sprintf("%s", uuid)

	err = DB.InsertUser(userData)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Insert User Error")
		response.ResultCode = "400"
		response.ResultMessage = ""
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)

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
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}

	err = DB.InsertFirendShip(friendShip)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Insert User Error")
		response.ResultCode = "400"
		response.ResultMessage = ""
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
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Get Current Price Error")
		response.ResultCode = "400"
		response.ResultMessage = err
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
		return
	}
	if oldPrice.Uuid != "" {
		err = DB.UpdatePrice(price)
		if err != nil {
			Log.WithFields(logrus.Fields{
				"Error": err,
			}).Error("Update Price Error")
			response.ResultCode = "400"
			response.ResultMessage = err
			Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)
			return
		}
	}

	err = DB.InsertPrice(price)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Insert Price Error")
		response.ResultCode = "400"
		response.ResultMessage = err
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
		response.ResultMessage = ""
		Comman.ResponseWithJson(w, http.StatusBadRequest, response, Log)

	}
	response.ResultCode = "200"
	response.ResultMessage = priceList
	Comman.ResponseWithJson(w, http.StatusOK, response, Log)
}
