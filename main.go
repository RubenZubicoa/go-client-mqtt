package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
    "time"
    "os"
)

func writeFile(text string){
    f, err := os.OpenFile("file.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    if _, err = f.WriteString(text); err != nil {
        panic(err)
    }
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    message := "Received message: "+ string(msg.Payload()) +" from topic: "+ string(msg.Topic())+"\n"
    fmt.Printf(message)
    writeFile(message)
}
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client){
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error){
	fmt.Printf("Connect lost: %v", err)
}

func sub(client mqtt.Client, topic string) {
    token := client.Subscribe(topic, 1, nil)
    token.Wait()
    fmt.Printf("Subscribed to topic %s\n", topic)
}

func wait(){
    time.Sleep(20 * time.Second)
}

func main(){
	var broker = "127.0.0.1"
	var port = 1883
	opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetClientID("go_mqtt_client")
    opts.SetUsername("emqx")
    opts.SetPassword("public")
    opts.SetDefaultPublishHandler(messagePubHandler)
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
  }
  topic := "topic/test"
  sub(client, topic)
  topic2 := "topic/prueba"
  sub(client, topic2)
  wait()
  client.Disconnect(250)
}