package master

import (
	"fmt"
	"net"
	"net/http"
	"time"
	"net/url"
)

func MonitorSlaveHeartbeats(_ http.ResponseWriter, request *http.Request, slaveMap map[string]Slave) {
	slaveName, slaveAddress := processRequest(request)

	if _, existsInMap := slaveMap[slaveName]; existsInMap {
		slaveMap = updateSlaveHeartbeat(slaveMap, slaveAddress, slaveName)
	} else {
		fmt.Printf("Slave added with name \"%v\", IP %v", slaveName, slaveAddress)
		slaveMap[slaveName] = Slave{URL: slaveAddress, heartbeat: time.Now()}
		sendSlaveListToWebserver(webserverAddress, slaveMap)
	}
}

func processRequest(request *http.Request) (slaveName, slaveAddress string) {
	slaveName = request.PostFormValue("slaveName")
	slavePort := request.PostFormValue("slavePort")

	slaveIP,_,_ := net.SplitHostPort(request.RemoteAddr)
	slaveAddress = "http://" + slaveIP + ":" + slavePort
	return
}

func updateSlaveHeartbeat(slaveMap map[string]Slave, slaveAddress, slaveName string) map[string]Slave {
	slaveInstance := slaveMap[slaveName]
	if slaveInstance.URL != slaveAddress {
		fmt.Printf(`WARNING: Slave with name \"%v\"
			already exists with the IP address: %v. \n
			kill signal sent to slave with name \"%v\"
			with IP address: %v`,
			slaveName, slaveInstance.URL, slaveName, slaveAddress)
		sendKillSignalToSlave(slaveAddress)
	} else {
		slaveInstance.heartbeat = time.Now()
		slaveMap[slaveName] = slaveInstance
	}
	return slaveMap
}

func sendKillSignalToSlave(slaveAddress string) {
	client :=  &http.Client{}
	form := url.Values{}
	form.Set("message", "die")
	client.PostForm(slaveAddress+"/receive_killsignal",form)
}

func MonitorSlaves(timeInterval int, slaveMap map[string]Slave) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)
	for _ = range timer {
		removeDeadSlaves(timeInterval, slaveMap)
	}
}

func removeDeadSlaves(deadTime int, slaveMap map[string]Slave) {
	for slaveName, slave := range slaveMap {
		if time.Now().Sub(slave.heartbeat) > time.Duration(deadTime)*time.Second {
			fmt.Printf("\nREMOVING DEAD SLAVE: %v\n", slaveName)
			delete(slaveMap, slaveName)
			fmt.Println("Current slaves are: ")
			if len(slaveMap) == 0 {
				fmt.Println("No slaves available.")
			} else {
				for slaveName, _ := range slaveMap {
					fmt.Println(slaveName)
				}
			}
			fmt.Printf("\n\n")
			sendSlaveListToWebserver(webserverAddress, slaveMap)
		}
	}
}

func UpdateWebserverAddress(w http.ResponseWriter, r *http.Request) {
	newWebserverAddress, _ := getWebserverAddress(r)
	if webserverAddress != newWebserverAddress {
		fmt.Println("Webserver address has changed from %v to %v",webserverAddress,newWebserverAddress)
		webserverAddress = newWebserverAddress
	}
}
