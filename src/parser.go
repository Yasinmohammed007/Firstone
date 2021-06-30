package src

import (
	"log"

	"github.com/SophosNSG/paloma-config-service/config/dns"
	"github.com/SophosNSG/paloma-config-service/config/record"
	"github.com/openconfig/ygot/ygot"
)

func ParseAndDeleteDNSRecodConfig(rawData []byte) (string, error) {
	log.Println("Now unmarshaling the rawData into Yang struct, and will validate")

	loadedJson := &record.DelFqdn{}

	if err := record.Unmarshal([]byte(rawData), loadedJson); err != nil {
		log.Println("Cannot unmarshal JSON \n", err)
		return "", err
	} else {
		log.Println("Unmarshal of JSON Done!!!")
	}

	err := loadedJson.Validate()
	if err != nil {
		log.Println("Error occured while validating the YANG \n", err)
		return "", err
	}
	log.Println("Unmarshal of JSON Done!!!, Also the data is validated against YANG Model")

	log.Println("Printing the validated yang in json format")
	json_val, err := ygot.EmitJSON(loadedJson, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	if err != nil {
		log.Println("JSON demo error")
	}
	log.Println(json_val)

	response, err := PublishMessageToTopic(NatsTopicDelDNS, json_val)
	if err != nil {
		log.Println("Error Occured while publishing the message", err.Error())
		return "", err
	}

	log.Println("Printing response", response)
	return response, nil
}

func ParseAndWriteDNSRecodConfig(rawData []byte) (string, error) {
	log.Println("Now unmarshaling the rawData into Yang struct, and will validate")

	loadedJson := &record.Dnsrecord{}

	if err := record.Unmarshal([]byte(rawData), loadedJson); err != nil {
		log.Println("Cannot unmarshal JSON \n", err)
		return "", err
	} else {
		log.Println("Unmarshal of JSON Done!!!")
	}

	err := loadedJson.Validate()
	if err != nil {
		log.Println("Error occured while validating the YANG \n", err)
		return "", err
	} else {
		log.Println("Unmarshal of JSON Done!!!, Also the data is validated against YANG Model")
	}

	log.Println("Printing the validated yang in json format")
	json_val, err := ygot.EmitJSON(loadedJson, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	if err != nil {
		log.Println("JSON demo error")
	}
	log.Println(json_val)

	response, err := PublishMessageToTopic(NatsTopicDNS, json_val)
	if err != nil {
		log.Println("Error Occured while publishing the message", err.Error())
		return "", err
	}

	log.Println("Printing response", response)
	return response, nil
}

func ParseAndPublishDNSConfig(rawData []byte) (string, error) {
	log.Println("Now unmarshaling the rawData into Yang struct, and will validate")

	loadedJson := &dns.Dnsconfig{}
	if err := dns.Unmarshal(rawData, loadedJson); err != nil {
		log.Println("Cannot unmarshal JSON \n", err)
		return "", err
	} else {
		log.Println("Unmarshal of JSON Done!!!")
	}

	err := loadedJson.Validate()
	if err != nil {
		log.Println("Error occured while validating the YANG \n", err)
		return "", err
	} else {
		log.Println("Unmarshal of JSON Done!!!, Also the data is validated against YANG Model")
	}

	log.Println("Printing the validated yang in json format")
	json_val, err := ygot.EmitJSON(loadedJson, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	if err != nil {
		log.Println("JSON demo error")
	}
	log.Println(json_val)

	// Publish the message to the messaging queue
	response, err := PublishRecord("nameserver", json_val)
	if err != nil {
		log.Println("Error Occured while publishing the message", err.Error())
		return "", err
	}

	return response, nil

}
