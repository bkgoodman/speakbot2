package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)


var client mqtt.Client


// This is for ACL updates, which Speakbot doesn't do
func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	if (message.Topic() == "ratt/control/broadcast/acl/update") {
		fmt.Println("Got ACL Update message")
	}
}

func PingSender() {

	for {
		var topic string = fmt.Sprintf("ratt/status/node/%s/ping",cfg.ClientID)
		var message string = "{\"status\":\"ok\"}"
    t:= client.Publish(topic,0,false,message)
    //log.Printf("MQTT Publishing \"%s\" topic \"%s\"",message,topic)
    t.Wait()
    if (t.Error() != nil) {
      log.Printf("MQTT publish error to %s: %s",topic,t.Error())
    }
    //log.Printf("MQTT published %v",t)
		time.Sleep(120 * time.Second)
	}
}

func mqtt_publish(topic string, message string) {
    //log.Printf("MQTT Publishing \"%s\" topic \"%s\"",message,topic)
    t:= client.Publish(topic,0,false,message)
    t.Wait()
    if (t.Error() != nil) {
      log.Printf("MQTT publish error to %s: %s",topic,t.Error())
    }
    //log.Printf("MQTT published %v",t)
}


func mqtt_init() {

	// MQTT broker address
	broker := fmt.Sprintf("ssl://%s:%d",cfg.MqttHost,cfg.MqttPort)

	// MQTT client ID
	clientID := cfg.ClientID

	// MQTT topic to subscribe to
	//topic := "#"

  if ((cfg.ClientCert == "") || (cfg.ClientKey == "") || (cfg.CACert == "")) {
		log.Fatal("MQTT Client specified, without TLC cert, key and CA")
  }

	// Load client key pair for TLS (replace with your own paths)
	cert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
	if err != nil {
		log.Fatal("Error loading X509 Keypair: ",err)
	}

		// Load your CA certificate (replace with your own path)
	caCert, err := ioutil.ReadFile(cfg.CACert)
	if err != nil {
		log.Fatal("Error reading CA file: ",cfg.CACert,err)
	}

	// Create a certificate pool and add your CA certificate
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	// Create a TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs: caPool,
	}

	// Create an MQTT client options
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetTLSConfig(tlsConfig).
		SetDefaultPublishHandler(onMessageReceived)

	// Create an MQTT client
	client = mqtt.NewClient(opts)

	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("MQTT Connect error: ",token.Error())
	}

	// Subscribe to the topic
	//if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
	//	log.Fatal("MQTT Subscribe error: ",token.Error())
	//}

	go PingSender()
	fmt.Printf("Connected to %s\n", broker)

}

func mqtt_destroy() {
	client.Disconnect(250)
	fmt.Println("Disconnected from the MQTT broker")
}

