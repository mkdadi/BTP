package ngsi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NGSI10Client struct {
	IoTBrokerURL string
}

func CtxElement2Object(ctxElem *ContextElement) *ContextObject {
	ctxObj := ContextObject{}
	ctxObj.Entity = ctxElem.Entity

	ctxObj.Attributes = make(map[string]ValueObject)
	for _, attr := range ctxElem.Attributes {
		ctxObj.Attributes[attr.Name] = ValueObject{Type: attr.Type, Value: attr.Value}
	}

	ctxObj.Metadata = make(map[string]ValueObject)
	for _, meta := range ctxElem.Metadata {
		ctxObj.Metadata[meta.Name] = ValueObject{Type: meta.Type, Value: meta.Value}
	}

	return &ctxObj
}

func Object2CtxElement(ctxObj *ContextObject) *ContextElement {
	ctxElement := ContextElement{}

	ctxElement.Entity = ctxObj.Entity

	ctxElement.Attributes = make([]ContextAttribute, 0)
	for name, attr := range ctxObj.Attributes {
		ctxAttr := ContextAttribute{Name: name, Type: attr.Type, Value: attr.Value}
		ctxElement.Attributes = append(ctxElement.Attributes, ctxAttr)
	}

	ctxElement.Metadata = make([]ContextMetadata, 0)
	for name, meta := range ctxObj.Metadata {
		ctxMeta := ContextMetadata{Name: name, Type: meta.Type, Value: meta.Value}
		ctxElement.Metadata = append(ctxElement.Metadata, ctxMeta)
	}

	return &ctxElement
}

func (nc *NGSI10Client) UpdateContext(ctxObj *ContextObject) error {
	elem := Object2CtxElement(ctxObj)

	updateCtxReq := &UpdateContextRequest{
		ContextElements: []ContextElement{*elem},
		UpdateAction:    "UPDATE",
	}

	body, err := json.Marshal(updateCtxReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/updateContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	updateCtxResp := UpdateContextResponse{}
	err = json.Unmarshal(text, &updateCtxResp)
	if err != nil {
		return err
	}

	if updateCtxResp.ErrorCode.Code == 200 {
		return nil
	} else {
		err = errors.New(updateCtxResp.ErrorCode.ReasonPhrase)
		return err
	}
}

func (nc *NGSI10Client) DeleteContext(eid *EntityId) error {
	element := ContextElement{}

	entity := EntityId{}
	entity.ID = eid.ID
	entity.Type = eid.Type
	entity.IsPattern = eid.IsPattern

	element.Entity = entity

	updateCtxReq := &UpdateContextRequest{
		ContextElements: []ContextElement{element},
		UpdateAction:    "DELETE",
	}

	body, err := json.Marshal(updateCtxReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/updateContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)

	updateCtxResp := UpdateContextResponse{}
	err = json.Unmarshal(text, &updateCtxResp)
	if err != nil {
		return err
	}

	if updateCtxResp.ErrorCode.Code == 200 {
		return nil
	} else {
		err = errors.New(updateCtxResp.ErrorCode.ReasonPhrase)
		return err
	}
}

func (nc *NGSI10Client) NotifyContext(elem *ContextElement) error {
	elementResponse := ContextElementResponse{}
	elementResponse.ContextElement = *elem
	elementResponse.StatusCode.Code = 200
	elementResponse.StatusCode.ReasonPhrase = "OK"

	notifyCtxReq := &NotifyContextRequest{
		ContextResponses: []ContextElementResponse{elementResponse},
	}

	body, err := json.Marshal(notifyCtxReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/notifyContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	notifyCtxResp := NotifyContextResponse{}
	err = json.Unmarshal(text, &notifyCtxResp)
	if err != nil {
		return err
	}

	if notifyCtxResp.ResponseCode.Code == 200 {
		return nil
	} else {
		err = errors.New(notifyCtxResp.ResponseCode.ReasonPhrase)
		return err
	}
}

func (nc *NGSI10Client) QueryContext(query *QueryContextRequest, headers *map[string]string) ([]*ContextObject, error) {
	body, err := json.Marshal(*query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/queryContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// if there is any additional header to add
	if headers != nil {
		for key, val := range *headers {
			req.Header.Add(key, val)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	queryCtxResp := QueryContextResponse{}
	err = json.Unmarshal(text, &queryCtxResp)
	if err != nil {
		return nil, err
	}

	ctxObjectList := make([]*ContextObject, 0)
	for _, contextElementResponse := range queryCtxResp.ContextResponses {
		ctxObj := CtxElement2Object(&contextElementResponse.ContextElement)
		ctxObjectList = append(ctxObjectList, ctxObj)
	}

	return ctxObjectList, nil
}

func (nc *NGSI10Client) InternalQueryContext(query *QueryContextRequest, headers *map[string]string) ([]ContextElement, error) {
	body, err := json.Marshal(*query)
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(body))

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/queryContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	for key, val := range *headers {
		req.Header.Add(key, val)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	queryCtxResp := QueryContextResponse{}
	err = json.Unmarshal(text, &queryCtxResp)
	if err != nil {
		return nil, err
	}

	ctxElements := make([]ContextElement, 0)
	for _, contextElementResponse := range queryCtxResp.ContextResponses {
		ctxElements = append(ctxElements, contextElementResponse.ContextElement)
	}

	return ctxElements, nil
}

func (nc *NGSI10Client) SubscribeContext(sub *SubscribeContextRequest, requireReliability bool) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}

	//fmt.Println(string(body))
	//fmt.Println(nc.IoTBrokerURL + "/subscribeContext")

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/subscribeContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if requireReliability == true {
		req.Header.Add("Require-Reliability", "true")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	subscribeCtxResp := SubscribeContextResponse{}
	err = json.Unmarshal(text, &subscribeCtxResp)
	if err != nil {
		return "", err
	}

	if subscribeCtxResp.SubscribeResponse.SubscriptionId != "" {
		return subscribeCtxResp.SubscribeResponse.SubscriptionId, nil
	} else {
		err = errors.New(subscribeCtxResp.SubscribeError.ErrorCode.ReasonPhrase)
		return "", err
	}
}

func (nc *NGSI10Client) UnsubscribeContext(sid string) error {
	unsubscription := &UnsubscribeContextRequest{
		SubscriptionId: sid,
	}

	body, err := json.Marshal(unsubscription)
	if err != nil {
		return err
	}

	//fmt.Println(string(body))

	req, err := http.NewRequest("POST", nc.IoTBrokerURL+"/unsubscribeContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	unsubscribeCtxResp := UnsubscribeContextResponse{}
	err = json.Unmarshal(text, &unsubscribeCtxResp)
	if err != nil {
		return err
	}

	if unsubscribeCtxResp.StatusCode.Code == 200 {
		return nil
	} else {
		err = errors.New(unsubscribeCtxResp.StatusCode.ReasonPhrase)
		return err
	}
}

type NGSI9Client struct {
	IoTDiscoveryURL string
}

func (nc *NGSI9Client) RegisterContext(registerCtxReq *RegisterContextRequest) (string, error) {
	body, err := json.Marshal(registerCtxReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/registerContext", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	registerCtxResp := RegisterContextResponse{}
	err = json.Unmarshal(text, &registerCtxResp)
	if err != nil {
		return "", err
	}

	if registerCtxResp.ErrorCode.Code == 200 {
		return registerCtxResp.RegistrationId, nil
	} else {
		err = errors.New(registerCtxResp.ErrorCode.ReasonPhrase)
		return "", err
	}
}

func (nc *NGSI9Client) UnregisterEntity(eid string) error {
	req, err := http.NewRequest("DELETE", nc.IoTDiscoveryURL+"/registration/"+eid, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (nc *NGSI9Client) DiscoverContextAvailability(discoverCtxAvailabilityReq *DiscoverContextAvailabilityRequest) ([]ContextRegistration, error) {
	body, err := json.Marshal(discoverCtxAvailabilityReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/discoverContextAvailability", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)

	discoverCtxAvailResp := DiscoverContextAvailabilityResponse{}
	err = json.Unmarshal(text, &discoverCtxAvailResp)
	if err != nil {
		return nil, err
	}

	registrationList := make([]ContextRegistration, 0)
	for _, registration := range discoverCtxAvailResp.ContextRegistrationResponses {
		registrationList = append(registrationList, registration.ContextRegistration)
	}

	return registrationList, nil
}

func (nc *NGSI9Client) SubscribeContextAvailability(sub *SubscribeContextAvailabilityRequest) (string, error) {
	body, err := json.Marshal(*sub)
	if err != nil {
		return "", err
	}

	fmt.Println(string(body))

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/subscribeContextAvailability", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(text))

	subscribeCtxAvailResp := SubscribeContextAvailabilityResponse{}
	err = json.Unmarshal(text, &subscribeCtxAvailResp)
	if err != nil {
		return "", err
	}

	if subscribeCtxAvailResp.SubscriptionId != "" {
		return subscribeCtxAvailResp.SubscriptionId, nil
	} else {
		err = errors.New(subscribeCtxAvailResp.ErrorCode.ReasonPhrase)
		return "", err
	}
}

func (nc *NGSI9Client) UnsubscribeContextAvailability(sid string) error {
	unsubscription := &UnsubscribeContextAvailabilityRequest{
		SubscriptionId: sid,
	}

	body, err := json.Marshal(unsubscription)
	if err != nil {
		return err
	}

	fmt.Println("unsubscribe the context availability from IoT Discovery")
	fmt.Println(string(body))

	req, err := http.NewRequest("POST", nc.IoTDiscoveryURL+"/unsubscribeContextAvailability", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(text))

	unsubscribeCtxAvailResp := UnsubscribeContextAvailabilityResponse{}
	err = json.Unmarshal(text, &unsubscribeCtxAvailResp)
	if err != nil {
		return err
	}

	if unsubscribeCtxAvailResp.StatusCode.Code == 200 {
		return nil
	} else {
		err = errors.New(unsubscribeCtxAvailResp.StatusCode.ReasonPhrase)
		return err
	}
}

func (nc *NGSI9Client) DiscoveryNearbyIoTBroker(nearby NearBy) (string, error) {
	discoverReq := DiscoverContextAvailabilityRequest{}

	entity := EntityId{}
	entity.Type = "IoTBroker"
	entity.IsPattern = true
	discoverReq.Entities = make([]EntityId, 0)

	discoverReq.Entities = append(discoverReq.Entities, entity)

	scope := OperationScope{}
	scope.Type = "nearby"
	scope.Value = nearby

	discoverReq.Restriction.Scopes = make([]OperationScope, 0)
	discoverReq.Restriction.Scopes = append(discoverReq.Restriction.Scopes, scope)

	registerationList, err := nc.DiscoverContextAvailability(&discoverReq)

	if err != nil {
		return "", err
	}

	if registerationList == nil {
		return "", nil
	} else {
		for _, reg := range registerationList {
			return reg.ProvidingApplication, nil
		}
	}
	return "", nil
}
