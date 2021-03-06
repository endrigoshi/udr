/*
 * Nudr_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package datarepository

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/src/udr/handler/message"
	"free5gc/src/udr/logger"

	"github.com/gin-gonic/gin"
)

// RemoveEeGroupSubscriptions - Deletes a eeSubscription for a group of UEs or any UE
func RemoveEeGroupSubscriptions(c *gin.Context) {
	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["ueGroupId"] = c.Params.ByName("ueGroupId")
	req.Params["subsId"] = c.Params.ByName("subsId")

	handlerMsg := message.NewHandlerMessage(message.EventRemoveEeGroupSubscriptions, req)
	message.SendMessage(handlerMsg)

	rsp := <-handlerMsg.ResponseChan

	HTTPResponse := rsp.HTTPResponse

	c.JSON(HTTPResponse.Status, HTTPResponse.Body)
}

// UpdateEeGroupSubscriptions - Stores an individual ee subscription of a group of UEs or any UE
func UpdateEeGroupSubscriptions(c *gin.Context) {
	var eeSubscription models.EeSubscription
	if err := c.ShouldBindJSON(&eeSubscription); err != nil {
		logger.DataRepoLog.Panic(err.Error())
	}

	req := http_wrapper.NewRequest(c.Request, eeSubscription)
	req.Params["ueGroupId"] = c.Params.ByName("ueGroupId")
	req.Params["subsId"] = c.Params.ByName("subsId")

	handlerMsg := message.NewHandlerMessage(message.EventUpdateEeGroupSubscriptions, req)
	message.SendMessage(handlerMsg)

	rsp := <-handlerMsg.ResponseChan

	HTTPResponse := rsp.HTTPResponse

	c.JSON(HTTPResponse.Status, HTTPResponse.Body)
}
