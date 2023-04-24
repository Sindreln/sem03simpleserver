package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "strconv"
    "strings"
    "sync"

    "github.com/uia-worker/is105sem03/mycrypt"
	
)

func celsiusToFahrenheit(celsius float64) float64 {
    return (celsius * 9 / 5) + 32
}

func main() {
    var wg sync.WaitGroup
    server, err := net.Listen("tcp", "172.17.0.3:8888")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("bundet til %s", server.Addr().String())
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            log.Println("før server.Accept() kallet")
            conn, err := server.Accept()
            if err != nil {
                return
            }
            go func(c net.Conn) {
                defer c.Close()
                for {
                    buf := make([]byte, 1024)
                    n, err := c.Read(buf)
                    if err != nil {
                        if err != io.EOF {
                            log.Println(err)
                        }
                        return // fra for løkke
                    }
                    
                  dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
                  log.Println("Dekryptert melding: ", string(dekryptertMelding))
                      
                    
                    switch msg := strings.TrimSpace(string(dekryptertMelding)); msg {
                    case "ping":
                    	_, err = c.Write([]byte("pong\x00"))
                    default:
                        if strings.HasPrefix(msg, "temperature:") {
                            parts := strings.Split(msg, ":")
                            if len(parts) != 2 {
                                log.Println("Ugyldig temperaturmelding:", msg)
                                continue
                            }
                            celsius, err := strconv.ParseFloat(parts[1], 64)
                            if err != nil {
                                log.Println("Kunne ikke parse Celsius-temperatur:", err)
                                continue
                            }
                            fahrenheit := celsiusToFahrenheit(celsius)
                            response := fmt.Sprintf("temperature:%.1f", fahrenheit)
                            _, err = c.Write([]byte(response))
                        } else {
				log.Println("Ukjent input", )
				continue
                        }
                    }
                    if err != nil {
                        if err != io.EOF {
                            log.Println(err)
                        }
                        return // fra for løkke
                    }
                }
            }(conn)
        }
    }()
    wg.Wait()

}