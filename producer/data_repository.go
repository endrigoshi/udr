package producer

import (
	"fmt"
	udr_context "free5gc/src/udr/context"
	"free5gc/src/udr/logger"

	// "context"
	"encoding/json"
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/src/udr/handler/message"
	"free5gc/src/udr/util"
	"net/http"
	"reflect"
	"strconv"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	// "strconv"
	// "strings"
)

var CurrentResourceUri string

func HandleCreateAccessAndMobilityData(respChan chan message.HandlerResponseMessage, ueId string, body models.AccessAndMobilityData) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleDeleteAccessAndMobilityData(respChan chan message.HandlerResponseMessage, ueId string) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleQueryAccessAndMobilityData(respChan chan message.HandlerResponseMessage, ueId string) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleQueryAmData(request *http_wrapper.Request) *http_wrapper.Response {
	logger.DataRepoLog.Infof("HandleQueryAmData")

	collName := "subscriptionData.provisionedData.amData"
	ueId := request.Params["ueId"]
	servingPlmnId := request.Params["servingPlmnId"]
	response, problemDetails := QueryAmDataProcedure(collName, ueId, servingPlmnId)

	if problemDetails == nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(http.StatusNotFound, nil, problemDetails)
	}
}

func QueryAmDataProcedure(collName string, ueId string, servingPlmnId string) (response *map[string]interface{}, problemDetails *models.ProblemDetails) {
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	accessAndMobilitySubscriptionData := RestfulAPIGetOne(collName, filter)
	if accessAndMobilitySubscriptionData != nil {
		return &accessAndMobilitySubscriptionData, nil
	} else {
		problemDetails = &models.ProblemDetails{}
		problemDetails.Cause = "USER_NOT_FOUND"
		return nil, problemDetails
	}
}

func HandleAmfContext3gpp(respChan chan message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.contextData.amf3gppAccess"
	filter := bson.M{"ueId": ueId}

	origValue := RestfulAPIGetOne(collName, filter)

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})

		newValue := RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleCreateAmfContext3gpp(respChan chan message.HandlerResponseMessage, ueId string, body models.Amf3GppAccessRegistration) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.contextData.amf3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryAmfContext3gpp(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.amf3gppAccess"
	filter := bson.M{"ueId": ueId}

	amf3GppAccessRegistration := RestfulAPIGetOne(collName, filter)

	if amf3GppAccessRegistration != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, amf3GppAccessRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleAmfContextNon3gpp(respChan chan message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.contextData.amfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	origValue := RestfulAPIGetOne(collName, filter)

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})

		newValue := RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleCreateAmfContextNon3gpp(respChan chan message.HandlerResponseMessage, ueId string, body models.AmfNon3GppAccessRegistration) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.contextData.amfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryAmfContextNon3gpp(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.amfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	amfNon3GppAccessRegistration := RestfulAPIGetOne(collName, filter)

	if amfNon3GppAccessRegistration != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, amfNon3GppAccessRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleModifyAuthentication(respChan chan message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	filter := bson.M{"ueId": ueId}

	origValue := RestfulAPIGetOne(collName, filter)

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})

		newValue := RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleQueryAuthSubsData(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	filter := bson.M{"ueId": ueId}

	authenticationSubscription := RestfulAPIGetOne(collName, filter)

	if authenticationSubscription != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, authenticationSubscription)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateAuthenticationSoR(respChan chan message.HandlerResponseMessage, ueId string, body models.SorData) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.ueUpdateConfirmationData.sorData"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryAuthSoR(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.ueUpdateConfirmationData.sorData"
	filter := bson.M{"ueId": ueId}

	sorData := RestfulAPIGetOne(collName, filter)

	if sorData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sorData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateAuthenticationStatus(respChan chan message.HandlerResponseMessage, ueId string, body models.AuthEvent) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.authenticationData.authenticationStatus"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryAuthenticationStatus(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.authenticationData.authenticationStatus"
	filter := bson.M{"ueId": ueId}

	authEvent := RestfulAPIGetOne(collName, filter)

	if authEvent != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, authEvent)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleApplicationDataInfluenceDataGet(respChan chan message.HandlerResponseMessage) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataInfluenceIdDelete(respChan chan message.HandlerResponseMessage, influenceId string) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataInfluenceIdPatch(respChan chan message.HandlerResponseMessage, influenceId string, body models.TrafficInfluDataPatch) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataInfluenceIdPut(respChan chan message.HandlerResponseMessage, influenceId string, body models.TrafficInfluData) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifyGet(respChan chan message.HandlerResponseMessage) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifyPost(respChan chan message.HandlerResponseMessage, body models.TrafficInfluSub) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdDelete(respChan chan message.HandlerResponseMessage, subscriptionId string) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdGet(respChan chan message.HandlerResponseMessage, subscriptionId string) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdPut(respChan chan message.HandlerResponseMessage, subscriptionId string, body models.TrafficInfluSub) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataPfdsAppIdDelete(respChan chan message.HandlerResponseMessage, pfdsAppId string) {
	collName := "applicationData.pfds"
	filter := bson.M{"applicationId": pfdsAppId}

	RestfulAPIDeleteOne(collName, filter)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleApplicationDataPfdsAppIdGet(respChan chan message.HandlerResponseMessage, pfdsAppId string) {
	collName := "applicationData.pfds"
	filter := bson.M{"applicationId": pfdsAppId}

	getData := RestfulAPIGetOne(collName, filter)

	if getData != nil {
		delete(getData, "_id")
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, getData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleApplicationDataPfdsAppIdPut(respChan chan message.HandlerResponseMessage, pfdsAppId string, body models.PfdDataForApp) {
	putData := toBsonM(body)

	collName := "applicationData.pfds"
	filter := bson.M{"applicationId": pfdsAppId}

	isExisted := RestfulAPIPutOne(collName, filter, putData)

	if isExisted {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, putData)

		//PreHandlePolicyDataChangeNotification("", pfdsAppId, body)
	} else {
		message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)
	}
}

func HandleApplicationDataPfdsGet(respChan chan message.HandlerResponseMessage, pfdsAppIdArray []string) {
	collName := "applicationData.pfds"
	filter := bson.M{}

	var pfdsArray []map[string]interface{}
	if len(pfdsAppIdArray) == 0 {
		pfdsArray = RestfulAPIGetMany(collName, filter)
		for i := 0; i < len(pfdsArray); i++ {
			delete(pfdsArray[i], "_id")
		}
	} else {
		for _, e := range pfdsAppIdArray {
			filter["applicationId"] = e
			getData := RestfulAPIGetOne(collName, filter)
			if getData != nil {
				delete(getData, "_id")
				pfdsArray = append(pfdsArray, getData)
			}
		}
	}

	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, pfdsArray)
}

func HandleExposureDataSubsToNotifyPost(respChan chan message.HandlerResponseMessage, body models.ExposureDataSubscription) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleExposureDataSubsToNotifySubIdDelete(respChan chan message.HandlerResponseMessage, subId string) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleExposureDataSubsToNotifySubIdPut(respChan chan message.HandlerResponseMessage, subId string, body models.ExposureDataSubscription) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandlePolicyDataBdtDataBdtReferenceIdDelete(respChan chan message.HandlerResponseMessage, bdtReferenceId string) {
	collName := "policyData.bdtData"
	filter := bson.M{"bdtReferenceId": bdtReferenceId}

	RestfulAPIDeleteOne(collName, filter)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandlePolicyDataBdtDataBdtReferenceIdGet(respChan chan message.HandlerResponseMessage, bdtReferenceId string) {
	collName := "policyData.bdtData"
	filter := bson.M{"bdtReferenceId": bdtReferenceId}

	bdtData := RestfulAPIGetOne(collName, filter)

	if bdtData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, bdtData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataBdtDataBdtReferenceIdPut(respChan chan message.HandlerResponseMessage, bdtReferenceId string, body models.BdtData) {
	putData := toBsonM(body)
	putData["bdtReferenceId"] = bdtReferenceId

	collName := "policyData.bdtData"
	filter := bson.M{"bdtReferenceId": bdtReferenceId}

	isExisted := RestfulAPIPutOne(collName, filter, putData)

	if isExisted {
		message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)

		PreHandlePolicyDataChangeNotification("", bdtReferenceId, body)
	} else {
		message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)

		// // TODO: need to check UPDATE_NOT_ALLOWED case
		// problemDetails := models.ProblemDetails{
		// 	Cause: "UPDATE_NOT_ALLOWED",
		// }
		// message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandlePolicyDataBdtDataGet(respChan chan message.HandlerResponseMessage) {
	collName := "policyData.bdtData"
	filter := bson.M{}

	bdtDataArray := RestfulAPIGetMany(collName, filter)

	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, bdtDataArray)
}

func HandlePolicyDataPlmnsPlmnIdUePolicySetGet(respChan chan message.HandlerResponseMessage, plmnId string) {
	collName := "policyData.plmns.uePolicySet"
	filter := bson.M{"plmnId": plmnId}

	uePolicySet := RestfulAPIGetOne(collName, filter)

	if uePolicySet != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, uePolicySet)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataSponsorConnectivityDataSponsorIdGet(respChan chan message.HandlerResponseMessage, sponsorId string) {
	collName := "policyData.sponsorConnectivityData"
	filter := bson.M{"sponsorId": sponsorId}

	sponsorConnectivityData := RestfulAPIGetOne(collName, filter)

	if sponsorConnectivityData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sponsorConnectivityData)
	} else {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	}
}

func HandlePolicyDataSubsToNotifyPost(respChan chan message.HandlerResponseMessage, body models.PolicyDataSubscription) {
	udrSelf := udr_context.UDR_Self()

	newSubscriptionID := strconv.Itoa(udrSelf.PolicyDataSubscriptionIDGenerator)
	udrSelf.PolicyDataSubscriptions[newSubscriptionID] = &body
	udrSelf.PolicyDataSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/policy-data/subs-to-notify{subsId} */
	locationHeader := fmt.Sprintf("%s/policy-data/subs-to-notify/%s", udrSelf.GetIPv4GroupUri(udr_context.NUDR_DR), newSubscriptionID)
	headers := http.Header{
		"Location": {locationHeader},
	}

	message.SendHttpResponseMessage(respChan, headers, http.StatusCreated, body)
}

func HandlePolicyDataSubsToNotifySubsIdDelete(respChan chan message.HandlerResponseMessage, subsId string) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.PolicyDataSubscriptions[subsId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}
	delete(udrSelf.PolicyDataSubscriptions, subsId)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandlePolicyDataSubsToNotifySubsIdPut(respChan chan message.HandlerResponseMessage, subsId string, body models.PolicyDataSubscription) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.PolicyDataSubscriptions[subsId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	udrSelf.PolicyDataSubscriptions[subsId] = &body

	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, body)
}

func HandlePolicyDataUesUeIdAmDataGet(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "policyData.ues.amData"
	filter := bson.M{"ueId": ueId}

	amPolicyData := RestfulAPIGetOne(collName, filter)

	if amPolicyData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, amPolicyData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataUesUeIdOperatorSpecificDataGet(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "policyData.ues.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	operatorSpecificDataContainerMapCover := RestfulAPIGetOne(collName, filter)

	if operatorSpecificDataContainerMapCover != nil {
		operatorSpecificDataContainerMap := operatorSpecificDataContainerMapCover["operatorSpecificDataContainerMap"]
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, operatorSpecificDataContainerMap)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataUesUeIdOperatorSpecificDataPatch(respChan chan message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "policyData.ues.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatchExtend(collName, filter, patchJSON, "operatorSpecificDataContainerMap")

	if success {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandlePolicyDataUesUeIdOperatorSpecificDataPut(respChan chan message.HandlerResponseMessage, ueId string, body map[string]models.OperatorSpecificDataContainer) {
	// json.NewDecoder(c.Request.Body).Decode(&operatorSpecificDataContainerMap)

	collName := "policyData.ues.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	putData := map[string]interface{}{"operatorSpecificDataContainerMap": body}
	putData["ueId"] = ueId

	RestfulAPIPutOne(collName, filter, putData)

	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandlePolicyDataUesUeIdSmDataGet(respChan chan message.HandlerResponseMessage, ueId string, snssai models.Snssai, dnn string) {
	collName := "policyData.ues.smData"
	filter := bson.M{"ueId": ueId}

	if !reflect.DeepEqual(snssai, models.Snssai{}) {
		filter["smPolicySnssaiData."+util.SnssaiModelsToHex(snssai)] = bson.M{"$exists": true}
	}
	if !reflect.DeepEqual(snssai, models.Snssai{}) && dnn != "" {
		filter["smPolicySnssaiData."+util.SnssaiModelsToHex(snssai)+".smPolicyDnnData."+dnn] = bson.M{"$exists": true}
	}

	smPolicyData := RestfulAPIGetOne(collName, filter)
	if smPolicyData != nil {
		var smPolicyDataResp models.SmPolicyData
		_ = json.Unmarshal(util.MapToByte(smPolicyData), &smPolicyDataResp)
		{
			collName := "policyData.ues.smData.usageMonData"
			filter := bson.M{"ueId": ueId}
			usageMonDataMapArray := RestfulAPIGetMany(collName, filter)

			if !reflect.DeepEqual(usageMonDataMapArray, []map[string]interface{}{}) {
				var usageMonDataArray []models.UsageMonData
				_ = json.Unmarshal(util.MapArrayToByte(usageMonDataMapArray), &usageMonDataArray)
				smPolicyDataResp.UmData = make(map[string]models.UsageMonData)
				for _, element := range usageMonDataArray {
					smPolicyDataResp.UmData[element.LimitId] = element
				}
			}
		}

		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smPolicyDataResp)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataUesUeIdSmDataPatch(respChan chan message.HandlerResponseMessage, ueId string, body map[string]models.UsageMonData) {
	collName := "policyData.ues.smData.usageMonData"
	filter := bson.M{"ueId": ueId}

	successAll := true
	for k, usageMonData := range body {
		limitId := k
		filterTmp := bson.M{"ueId": ueId, "limitId": limitId}
		success := RestfulAPIMergePatch(collName, filterTmp, toBsonM(usageMonData))
		if !success {
			successAll = false
		} else {
			var usageMonData models.UsageMonData
			usageMonDataBsonM := RestfulAPIGetOne(collName, filterTmp)
			_ = json.Unmarshal(util.MapToByte(usageMonDataBsonM), &usageMonData)
			PreHandlePolicyDataChangeNotification(ueId, limitId, usageMonData)
		}
	}

	if successAll {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})

		smPolicyDataBsonM := RestfulAPIGetOne(collName, filter)
		var smPolicyData models.SmPolicyData
		_ = json.Unmarshal(util.MapToByte(smPolicyDataBsonM), &smPolicyData)
		{
			collName := "policyData.ues.smData.usageMonData"
			filter := bson.M{"ueId": ueId}
			usageMonDataMapArray := RestfulAPIGetMany(collName, filter)

			if !reflect.DeepEqual(usageMonDataMapArray, []map[string]interface{}{}) {
				var usageMonDataArray []models.UsageMonData
				_ = json.Unmarshal(util.MapArrayToByte(usageMonDataMapArray), &usageMonDataArray)
				smPolicyData.UmData = make(map[string]models.UsageMonData)
				for _, element := range usageMonDataArray {
					smPolicyData.UmData[element.LimitId] = element
				}
			}
		}
		PreHandlePolicyDataChangeNotification(ueId, "", smPolicyData)
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandlePolicyDataUesUeIdSmDataUsageMonIdDelete(respChan chan message.HandlerResponseMessage, ueId string, usageMonId string) {
	collName := "policyData.ues.smData.usageMonData"
	filter := bson.M{"ueId": ueId, "usageMonId": usageMonId}

	RestfulAPIDeleteOne(collName, filter)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandlePolicyDataUesUeIdSmDataUsageMonIdGet(respChan chan message.HandlerResponseMessage, ueId string, usageMonId string) {
	collName := "policyData.ues.smData.usageMonData"
	filter := bson.M{"ueId": ueId, "usageMonId": usageMonId}

	usageMonData := RestfulAPIGetOne(collName, filter)

	if usageMonData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, usageMonData)
	} else {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	}
}

func HandlePolicyDataUesUeIdSmDataUsageMonIdPut(respChan chan message.HandlerResponseMessage, ueId string, usageMonId string, body models.UsageMonData) {
	putData := toBsonM(body)
	putData["ueId"] = ueId
	putData["usageMonId"] = usageMonId

	collName := "policyData.ues.smData.usageMonData"
	filter := bson.M{"ueId": ueId, "usageMonId": usageMonId}

	RestfulAPIPutOne(collName, filter, putData)

	message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)
}

func HandlePolicyDataUesUeIdUePolicySetGet(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "policyData.ues.uePolicySet"
	filter := bson.M{"ueId": ueId}

	uePolicySet := RestfulAPIGetOne(collName, filter)

	if uePolicySet != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, uePolicySet)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataUesUeIdUePolicySetPatch(respChan chan message.HandlerResponseMessage, ueId string, body models.UePolicySet) {
	patchData := toBsonM(body)
	patchData["ueId"] = ueId

	collName := "policyData.ues.uePolicySet"
	filter := bson.M{"ueId": ueId}

	success := RestfulAPIMergePatch(collName, filter, patchData)

	if success {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})

		var uePolicySet models.UePolicySet
		uePolicySetBsonM := RestfulAPIGetOne(collName, filter)
		_ = json.Unmarshal(util.MapToByte(uePolicySetBsonM), &uePolicySet)
		PreHandlePolicyDataChangeNotification(ueId, "", uePolicySet)
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}

}

func HandlePolicyDataUesUeIdUePolicySetPut(respChan chan message.HandlerResponseMessage, ueId string, body models.UePolicySet) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "policyData.ues.uePolicySet"
	filter := bson.M{"ueId": ueId}

	isExisted := RestfulAPIPutOne(collName, filter, putData)

	if !isExisted {
		message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)
	} else {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	}
}

func HandleCreateAMFSubscriptions(respChan chan message.HandlerResponseMessage, ueId string, subsId string, body []models.AmfSubscriptionInfo) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}
	udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId].AmfSubscriptionInfos = body

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleRemoveAmfSubscriptionsInfo(respChan chan message.HandlerResponseMessage, ueId string, subsId string) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	if udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "AMFSUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId].AmfSubscriptionInfos = nil

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleModifyAmfSubscriptionInfo(respChan chan message.HandlerResponseMessage, ueId string, subsId string, patchItem []models.PatchItem) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	if udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "AMFSUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	patchJSON, _ := json.Marshal(patchItem)
	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "MODIFY_NOT_ALLOWED"
		problemDetails.Detail = "PatchItem attributes are invalid"
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
		return
	}
	original, _ := json.Marshal((udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId]).AmfSubscriptionInfos)
	modified, err := patch.Apply(original)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "MODIFY_NOT_ALLOWED"
		problemDetails.Detail = "Occur error when applying PatchItem"
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
		return
	}
	var modifiedData []models.AmfSubscriptionInfo
	_ = json.Unmarshal(modified, &modifiedData)

	udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId].AmfSubscriptionInfos = modifiedData

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleGetAmfSubscriptionInfo(respChan chan message.HandlerResponseMessage, ueId string, subsId string) {
	udrSelf := udr_context.UDR_Self()

	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	if udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId].AmfSubscriptionInfos == nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "AMFSUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId].AmfSubscriptionInfos)
}

func HandleQueryEEData(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.eeProfileData"
	filter := bson.M{"ueId": ueId}

	eeProfileData := RestfulAPIGetOne(collName, filter)

	if eeProfileData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, eeProfileData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleRemoveEeGroupSubscriptions(respChan chan message.HandlerResponseMessage, ueGroupId string, subsId string) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.UEGroupCollection[ueGroupId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UEGroupCollection[ueGroupId].EeSubscriptions[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}
	delete(udrSelf.UEGroupCollection[ueGroupId].EeSubscriptions, subsId)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleUpdateEeGroupSubscriptions(respChan chan message.HandlerResponseMessage, ueGroupId string, subsId string, body models.EeSubscription) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.UEGroupCollection[ueGroupId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UEGroupCollection[ueGroupId].EeSubscriptions[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}
	udrSelf.UEGroupCollection[ueGroupId].EeSubscriptions[subsId] = &body

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleCreateEeGroupSubscriptions(respChan chan message.HandlerResponseMessage, ueGroupId string, body models.EeSubscription) {
	udrSelf := udr_context.UDR_Self()

	_, ok := udrSelf.UEGroupCollection[ueGroupId]
	if !ok {
		udrSelf.UEGroupCollection[ueGroupId] = new(udr_context.UEGroupSubsData)
	}
	if udrSelf.UEGroupCollection[ueGroupId].EeSubscriptions == nil {
		udrSelf.UEGroupCollection[ueGroupId].EeSubscriptions = make(map[string]*models.EeSubscription)
	}

	newSubscriptionID := strconv.Itoa(udrSelf.EeSubscriptionIDGenerator)
	udrSelf.UEGroupCollection[ueGroupId].EeSubscriptions[newSubscriptionID] = &body
	udrSelf.EeSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/nudr-dr/v1/subscription-data/group-data/{ueGroupId}/ee-subscriptions */
	locationHeader := fmt.Sprintf("%s/nudr-dr/v1/subscription-data/group-data/%s/ee-subscriptions/%s", udrSelf.GetIPv4GroupUri(udr_context.NUDR_DR), ueGroupId, newSubscriptionID)
	headers := http.Header{
		"Location": {locationHeader},
	}

	message.SendHttpResponseMessage(respChan, headers, http.StatusCreated, body)
}

func HandleQueryEeGroupSubscriptions(respChan chan message.HandlerResponseMessage, ueGroupId string) {
	udrSelf := udr_context.UDR_Self()

	_, ok := udrSelf.UEGroupCollection[ueGroupId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	var eeSubscriptionSlice []models.EeSubscription

	for _, v := range udrSelf.UEGroupCollection[ueGroupId].EeSubscriptions {
		eeSubscriptionSlice = append(eeSubscriptionSlice, *v)
	}

	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, eeSubscriptionSlice)
}

func HandleRemoveeeSubscriptions(respChan chan message.HandlerResponseMessage, ueId string, subsId string) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}
	delete(udrSelf.UESubsCollection[ueId].EeSubscriptionCollection, subsId)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleUpdateEesubscriptions(respChan chan message.HandlerResponseMessage, ueId string, subsId string, body models.EeSubscription) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}
	udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[subsId].EeSubscriptions = &body

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleCreateEeSubscriptions(respChan chan message.HandlerResponseMessage, ueId string, body models.EeSubscription) {
	udrSelf := udr_context.UDR_Self()

	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		udrSelf.UESubsCollection[ueId] = new(udr_context.UESubsData)
	}
	if udrSelf.UESubsCollection[ueId].EeSubscriptionCollection == nil {
		udrSelf.UESubsCollection[ueId].EeSubscriptionCollection = make(map[string]*udr_context.EeSubscriptionCollection)
	}

	newSubscriptionID := strconv.Itoa(udrSelf.EeSubscriptionIDGenerator)
	udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[newSubscriptionID] = new(udr_context.EeSubscriptionCollection)
	udrSelf.UESubsCollection[ueId].EeSubscriptionCollection[newSubscriptionID].EeSubscriptions = &body
	udrSelf.EeSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/subscription-data/{ueId}/context-data/ee-subscriptions/{subsId} */
	locationHeader := fmt.Sprintf("%s/subscription-data/%s/context-data/ee-subscriptions/%s", udrSelf.GetIPv4GroupUri(udr_context.NUDR_DR), ueId, newSubscriptionID)
	headers := http.Header{
		"Location": {locationHeader},
	}

	message.SendHttpResponseMessage(respChan, headers, http.StatusCreated, body)
}

func HandleQueryeesubscriptions(respChan chan message.HandlerResponseMessage, ueId string) {
	udrSelf := udr_context.UDR_Self()

	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	var eeSubscriptionSlice []models.EeSubscription

	for _, v := range udrSelf.UESubsCollection[ueId].EeSubscriptionCollection {
		eeSubscriptionSlice = append(eeSubscriptionSlice, *v.EeSubscriptions)
	}

	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, eeSubscriptionSlice)
}

func HandlePatchOperSpecData(respChan chan message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	origValue := RestfulAPIGetOne(collName, filter)

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})

		newValue := RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleQueryOperSpecData(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	operatorSpecificDataContainer := RestfulAPIGetOne(collName, filter)

	// The key of the map is operator specific data element name and the value is the operator specific data of the UE.

	if operatorSpecificDataContainer != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, operatorSpecificDataContainer)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetppData(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.ppData"
	filter := bson.M{"ueId": ueId}

	ppData := RestfulAPIGetOne(collName, filter)

	if ppData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, ppData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateSessionManagementData(respChan chan message.HandlerResponseMessage, ueId string, pduSessionId int32, body models.PduSessionManagementData) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleDeleteSessionManagementData(respChan chan message.HandlerResponseMessage, ueId string, pduSessionId int32) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleQuerySessionManagementData(respChan chan message.HandlerResponseMessage, ueId string, pduSessionId int32) {
	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleQueryProvisionedData(respChan chan message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	var provisionedDataSets models.ProvisionedDataSets

	{
		collName := "subscriptionData.provisionedData.amData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		accessAndMobilitySubscriptionData := RestfulAPIGetOne(collName, filter)
		if accessAndMobilitySubscriptionData != nil {
			var tmp models.AccessAndMobilitySubscriptionData
			err := mapstructure.Decode(accessAndMobilitySubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.AmData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smfSelectionSubscriptionData := RestfulAPIGetOne(collName, filter)
		if smfSelectionSubscriptionData != nil {
			var tmp models.SmfSelectionSubscriptionData
			err := mapstructure.Decode(smfSelectionSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmfSelData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smsData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smsSubscriptionData := RestfulAPIGetOne(collName, filter)
		if smsSubscriptionData != nil {
			var tmp models.SmsSubscriptionData
			err := mapstructure.Decode(smsSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmsSubsData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		sessionManagementSubscriptionDatas := RestfulAPIGetMany(collName, filter)
		if sessionManagementSubscriptionDatas != nil {
			var tmp []models.SessionManagementSubscriptionData
			err := mapstructure.Decode(sessionManagementSubscriptionDatas, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmData = tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.traceData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		traceData := RestfulAPIGetOne(collName, filter)
		if traceData != nil {
			var tmp models.TraceData
			err := mapstructure.Decode(traceData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.TraceData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smsMngData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smsManagementSubscriptionData := RestfulAPIGetOne(collName, filter)
		if smsManagementSubscriptionData != nil {
			var tmp models.SmsManagementSubscriptionData
			err := mapstructure.Decode(smsManagementSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmsMngData = &tmp
		}
	}

	if !reflect.DeepEqual(provisionedDataSets, models.ProvisionedDataSets{}) {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, provisionedDataSets)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleModifyPpData(respChan chan message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.ppData"
	filter := bson.M{"ueId": ueId}

	origValue := RestfulAPIGetOne(collName, filter)

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})

		newValue := RestfulAPIGetOne(collName, filter)
		PreHandleOnDataChangeNotify(ueId, CurrentResourceUri, patchItem, origValue, newValue)
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleGetIdentityData(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.identityData"
	filter := bson.M{"ueId": ueId}

	identityData := RestfulAPIGetOne(collName, filter)

	if identityData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, identityData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetOdbData(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.operatorDeterminedBarringData"
	filter := bson.M{"ueId": ueId}

	operatorDeterminedBarringData := RestfulAPIGetOne(collName, filter)

	if operatorDeterminedBarringData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, operatorDeterminedBarringData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetSharedData(respChan chan message.HandlerResponseMessage, sharedDataIds []string) {

	collName := "subscriptionData.sharedData"
	var sharedDataArray []map[string]interface{}
	for _, sharedDataId := range sharedDataIds {
		filter := bson.M{"sharedDataId": sharedDataId}
		sharedData := RestfulAPIGetOne(collName, filter)
		if sharedData != nil {
			sharedDataArray = append(sharedDataArray, sharedData)
		}
	}

	if sharedDataArray != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sharedDataArray)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleRemovesdmSubscriptions(respChan chan message.HandlerResponseMessage, ueId string, subsId string) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UESubsCollection[ueId].SdmSubscriptions[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}
	delete(udrSelf.UESubsCollection[ueId].SdmSubscriptions, subsId)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleUpdatesdmsubscriptions(respChan chan message.HandlerResponseMessage, ueId string, subsId string, body models.SdmSubscription) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	_, ok = udrSelf.UESubsCollection[ueId].SdmSubscriptions[subsId]

	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}
	body.SubscriptionId = subsId
	udrSelf.UESubsCollection[ueId].SdmSubscriptions[subsId] = &body

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleCreateSdmSubscriptions(respChan chan message.HandlerResponseMessage, ueId string, body models.SdmSubscription) {
	udrSelf := udr_context.UDR_Self()

	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		udrSelf.UESubsCollection[ueId] = new(udr_context.UESubsData)
	}
	if udrSelf.UESubsCollection[ueId].SdmSubscriptions == nil {
		udrSelf.UESubsCollection[ueId].SdmSubscriptions = make(map[string]*models.SdmSubscription)
	}

	newSubscriptionID := strconv.Itoa(udrSelf.SdmSubscriptionIDGenerator)
	body.SubscriptionId = newSubscriptionID
	udrSelf.UESubsCollection[ueId].SdmSubscriptions[newSubscriptionID] = &body
	udrSelf.SdmSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/subscription-data/{ueId}/context-data/sdm-subscriptions/{subsId}' */
	locationHeader := fmt.Sprintf("%s/subscription-data/%s/context-data/sdm-subscriptions/%s", udrSelf.GetIPv4GroupUri(udr_context.NUDR_DR), ueId, newSubscriptionID)
	headers := http.Header{
		"Location": {locationHeader},
	}

	message.SendHttpResponseMessage(respChan, headers, http.StatusCreated, body)
}

func HandleQuerysdmsubscriptions(respChan chan message.HandlerResponseMessage, ueId string) {
	udrSelf := udr_context.UDR_Self()

	_, ok := udrSelf.UESubsCollection[ueId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	var sdmSubscriptionSlice []models.SdmSubscription

	for _, v := range udrSelf.UESubsCollection[ueId].SdmSubscriptions {
		sdmSubscriptionSlice = append(sdmSubscriptionSlice, *v)
	}

	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sdmSubscriptionSlice)
}

func HandleQuerySmData(respChan chan message.HandlerResponseMessage, ueId string, servingPlmnId string, singleNssai models.Snssai, dnn string) {
	collName := "subscriptionData.provisionedData.smData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	if !reflect.DeepEqual(singleNssai, models.Snssai{}) {
		if singleNssai.Sd == "" {
			filter["singleNssai.sst"] = singleNssai.Sst
		} else {
			filter["singleNssai.sst"] = singleNssai.Sst
			filter["singleNssai.sd"] = singleNssai.Sd
		}
	}

	if dnn != "" {
		filter["dnnConfigurations."+dnn] = bson.M{"$exists": true}
	}

	sessionManagementSubscriptionDatas := RestfulAPIGetMany(collName, filter)

	message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sessionManagementSubscriptionDatas)
}

func HandleCreateSmfContextNon3gpp(respChan chan message.HandlerResponseMessage, ueId string, pduSessionId int32, body models.SmfRegistration) {
	pduSessionIdInt := pduSessionId

	putData := toBsonM(body)
	putData["ueId"] = ueId
	putData["pduSessionId"] = int32(pduSessionIdInt)

	collName := "subscriptionData.contextData.smfRegistrations"
	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	isExisted := RestfulAPIPutOne(collName, filter, putData)

	if !isExisted {
		message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)
	} else {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, putData)
	}
}

func HandleDeleteSmfContext(respChan chan message.HandlerResponseMessage, ueId string, pduSessionId string) {
	pduSessionIdInt, _ := strconv.ParseInt(pduSessionId, 10, 32)

	collName := "subscriptionData.contextData.smfRegistrations"
	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	RestfulAPIDeleteOne(collName, filter)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQuerySmfRegistration(respChan chan message.HandlerResponseMessage, ueId string, pduSessionId string) {
	pduSessionIdInt, _ := strconv.ParseInt(pduSessionId, 10, 32)

	collName := "subscriptionData.contextData.smfRegistrations"
	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	smfRegistration := RestfulAPIGetOne(collName, filter)

	if smfRegistration != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smfRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleQuerySmfRegList(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smfRegistrations"
	filter := bson.M{"ueId": ueId}

	smfRegList := RestfulAPIGetMany(collName, filter)

	if smfRegList != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smfRegList)
	} else {
		// Return empty array instead
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, []map[string]interface{}{})
	}

}

func HandleQuerySmfSelectData(respChan chan message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	smfSelectionSubscriptionData := RestfulAPIGetOne(collName, filter)

	if smfSelectionSubscriptionData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smfSelectionSubscriptionData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateSmsfContext3gpp(respChan chan message.HandlerResponseMessage, ueId string, body models.SmsfRegistration) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.contextData.smsf3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleDeleteSmsfContext3gpp(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smsf3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIDeleteOne(collName, filter)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQuerySmsfContext3gpp(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smsf3gppAccess"
	filter := bson.M{"ueId": ueId}

	smsfRegistration := RestfulAPIGetOne(collName, filter)

	if smsfRegistration != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smsfRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateSmsfContextNon3gpp(respChan chan message.HandlerResponseMessage, ueId string, body models.SmsfRegistration) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.contextData.smsfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleDeleteSmsfContextNon3gpp(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smsfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIDeleteOne(collName, filter)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQuerySmsfContextNon3gpp(respChan chan message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smsfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	smsfRegistration := RestfulAPIGetOne(collName, filter)

	if smsfRegistration != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smsfRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleQuerySmsMngData(respChan chan message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smsMngData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	smsManagementSubscriptionData := RestfulAPIGetOne(collName, filter)

	if smsManagementSubscriptionData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smsManagementSubscriptionData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleQuerySmsData(respChan chan message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smsData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	smsSubscriptionData := RestfulAPIGetOne(collName, filter)

	if smsSubscriptionData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smsSubscriptionData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePostSubscriptionDataSubscriptions(respChan chan message.HandlerResponseMessage, body models.SubscriptionDataSubscriptions) {
	udrSelf := udr_context.UDR_Self()

	newSubscriptionID := strconv.Itoa(udrSelf.SubscriptionDataSubscriptionIDGenerator)
	udrSelf.SubscriptionDataSubscriptions[newSubscriptionID] = &body
	udrSelf.SubscriptionDataSubscriptionIDGenerator++

	/* Contains the URI of the newly created resource, according
	   to the structure: {apiRoot}/subscription-data/subs-to-notify/{subsId} */
	locationHeader := fmt.Sprintf("%s/subscription-data/subs-to-notify/%s", udrSelf.GetIPv4GroupUri(udr_context.NUDR_DR), newSubscriptionID)
	headers := http.Header{
		"Location": {locationHeader},
	}

	message.SendHttpResponseMessage(respChan, headers, http.StatusCreated, body)
}

func HandleRemovesubscriptionDataSubscriptions(respChan chan message.HandlerResponseMessage, subsId string) {
	udrSelf := udr_context.UDR_Self()
	_, ok := udrSelf.SubscriptionDataSubscriptions[subsId]
	if !ok {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SUBSCRIPTION_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}
	delete(udrSelf.SubscriptionDataSubscriptions, subsId)

	message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryTraceData(respChan chan message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.traceData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	traceData := RestfulAPIGetOne(collName, filter)

	if traceData != nil {
		message.SendHttpResponseMessage(respChan, nil, http.StatusOK, traceData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}
