package detect

//Common code to find an aurora, authenticate to it, and store credentials on the system

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/captncraig/aurora"
	"github.com/hashicorp/mdns"
)

func FindClient() (*aurora.Client, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	cacheFile := filepath.Join(u.HomeDir, ".aurora")
	cached, err := ioutil.ReadFile(cacheFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	var server, token string
	if len(cached) > 0 {
		parts := strings.Split(string(cached), "\n")
		if len(parts) == 2 {
			server, token = parts[0], parts[1]
		}
	} else {
		log.Println("No valid cached hardware info found, searching")
		server, err = discover()
		if err != nil {
			return nil, err
		}
	}
	var client *aurora.Client

	//TODO: check token if found. Set empty if invalid

	if token == "" {
		client = aurora.New(server)
		log.Println("Need to authenticate to device. Hold power button to enter pairing mode. Press enter when ready.")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		token, err = client.Authorize()
		if err != nil {
			return nil, err
		}
	} else {
		client = aurora.NewWithToken(server, token)
	}

	return client, ioutil.WriteFile(cacheFile, []byte(fmt.Sprintf("%s\n%s", server, token)), 0644)
}

func discover() (string, error) {
	//bonjour lib is chatty
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	ch := make(chan *mdns.ServiceEntry, 5)
	mdns.Lookup("_nanoleafapi._tcp", ch)

	done := time.After(2000 * time.Millisecond)
	possibilities := []string{}
L:
	for {

		select {
		case entry := <-ch:
			possibilities = append(possibilities, "http://"+entry.AddrV4.String()+":"+fmt.Sprint(entry.Port))
		case <-done:
			break L
		}
	}
	close(ch)
	if len(possibilities) > 1 {
		fmt.Println("MORE THAN ONE SERVER FOUND!!! IMPLEMENT SELECTION UI!!!!!")
		os.Exit(1)
	}
	if len(possibilities) == 0 {
		fmt.Println("NO AURORAs FOUND VIA BONJOUR")
		os.Exit(1)
	}
	return possibilities[0], nil
}
