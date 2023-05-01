package main

import (
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func getValuesFromApi() {
	req, err := http.Get(CbrApiUrl)
	if err != nil {
		log.Fatal("Request error: " + err.Error())
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	var result CBRApiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal("Unmarshall error: " + err.Error())
	}

	for _, curr := range defaultCurrenciesArray {
		currency := result.Valute[curr]
		msg := fmt.Sprintf("%s = %f RUB", curr, currency.Value)
		showNotification(msg)
	}
}

func showNotification(msg string) {
	cmd := exec.Command("notify-send", msg)
	_, err := cmd.Output()
	if err != nil {
		log.Fatal("Notification error: " + err.Error())
	}
}

func start(showError bool) {
	getValuesFromApi()
}

func main() {
	runtime.GOMAXPROCS(1)

	parser := argparse.NewParser("currencyChangeNotifier", "Program notifies you about changes in the exchange rate")
	repeatApiRequestArg := parser.Int("r", "repeat", &argparse.Options{Required: false, Help: "Check currency data every some minutes. Default is 5 minutes", Default: 5})
	notifyAboutErrorsArg := parser.Flag("s", "show_error", &argparse.Options{Required: false, Help: "Notify about program errors", Default: false})
	//currenciesListArg := parser.List("a", "additional_currencies", &argparse.Options{Required: false, Help: fmt.Sprintf("Additional list of currencies. Default are %s", strings.Join(defaultCurrenciesArray[:], ", ")), Default: defaultCurrenciesArray})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	start(*notifyAboutErrorsArg)

	ticker := time.NewTicker(time.Duration(*repeatApiRequestArg) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			start(*notifyAboutErrorsArg)
		}
	}

}
