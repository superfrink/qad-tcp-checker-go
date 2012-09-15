package main

import "bufio"
import "fmt"
import "log"
import "math/rand"
import "net"
import "os"
import "strconv"
import "time"

func format_statistics(stats map[string]map[string]int) string {

	str := ""

	for host, host_map := range stats {

		str += fmt.Sprintf("%15s", host);
		for status, count := range host_map {
			str += fmt.Sprintf(" %s %d\t", status, count);
		}
		str += fmt.Sprintln();
	}

	return str
}

func aggregator(input chan []string, dump_stats chan bool) {
	statistics := make(map[string]map[string]int)

	go func() {
		var message []string
		for {
			select {
			case message = <- input :
				host   := message[0]
				status := message[1]

				if("success" != status) {
					fmt.Println(message[0] + " " + message[1])
				}

				if _, ok := statistics[host]; !ok {
					statistics[host] = make(map[string]int)
				}

				if _, ok := statistics[host][status]; !ok {
					statistics[host][status] = 0
				}
				statistics[host][status] ++;

			case <- dump_stats:
				fmt.Printf("%v\n", statistics)

				fmt.Print(format_statistics(statistics))
			}
		}
	}()
}

func check_host(output chan []string, host string) {
	port := 80  // FIXME: make a parameter
	go func() {
		for {
			// check the network status for the host
			conn, err := net.Dial("tcp", host + ":" + strconv.Itoa(port))
			if err != nil {
				// CLAIM: error
				output <- []string{host, "failure"}
			} else {
				// CLAIM: success
				output <- []string{host, "success"}
				conn.Close()
			}

			// wait for a bit before checking again
			time.Sleep(time.Second * time.Duration(rand.Intn(1) + 1))
		}
	}()
}

func check_host_list(output chan []string, host_list []string) {

	for _, host := range host_list {
		check_host(output, host)
	}
}

func main() {

	// Create the stats aggregator
	status := make(chan []string)
	dump_stats := make(chan bool)
	aggregator(status, dump_stats)
	status <- []string{"test", "success"}
	status <- []string{"test", "failure"}
	status <- []string{"test", "success"}

	// Start collecting stats
	host_list := []string{"a", "superfrink.net", "127.0.0.1"}
	check_host_list(status, host_list)

	reader := bufio.NewReader(os.Stdin)
	for {
		_, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		dump_stats <- true
	}
	// sit_forever := make(chan bool)
	// <- sit_forever

	fmt.Println("Good-bye.")
}