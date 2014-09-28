package dht

import (
	"fmt"
	log "github.com/cihub/seelog"
	"strings"
	"testing"
	"time"
)

func validateTestResults(test *testing.T, validResults []string) {
	waitForValidTest := make(chan bool)
	logReceiver := &CustomReceiver{channel: waitForValidTest}
	log.RegisterReceiver("tester", logReceiver)
	select {
	case <-waitForValidTest:
	case <-time.After(timeoutSeconds * 5):
		test.Error("Test timed out")
	}
}

// Receiver that receives log messages and puts them in an array
// This can be used to test expected behavior of network requests
type CustomReceiver struct { // implements seelog.CustomReceiver
	validLogs []string
	channel   chan bool
	prefix    string
	test      *testing.T
}

func (ar *CustomReceiver) ReceiveMessage(message string, level log.LogLevel, context log.LogContextInterface) error {
	message = message[:len(message)-1]
	if message == ar.validLogs[0] {
		ar.validLogs = append(ar.validLogs[:0], ar.validLogs[0+1:]...)
	} else {
		fmt.Printf("Error, expecting '%s', got '%s'\n", ar.validLogs[0], message)
	}
	return nil
}

func (ar *CustomReceiver) AfterParse(initArgs log.CustomReceiverInitArgs) error {
	var ok bool
	var validString string
	validString, ok = initArgs.XmlCustomAttrs["prefix"]
	ar.validLogs = strings.Split(validString, ",")
	if !ok {
		ar.prefix = "No prefix"
	}
	return nil
}

func (ar *CustomReceiver) Flush() {

}

func (ar *CustomReceiver) Close() error {
	return nil
}

func configTestValidator(valid []string) {
	testValidator := &CustomReceiver{}
	log.RegisterReceiver("tester", testValidator)
	testConfig := `
		<seelog type="sync">
			<outputs>
				<file formatid="onlytime" path="logfile.log"/>
				<custom name="tester" formatid="onlymsg" data-prefix="` + strings.Join(valid, ",") + `"/>
			</outputs>
			<formats>
				<format id="default" format="%Date %Time [%LEVEL] %Msg%n"/>
				<format id="onlytime" format="%Time [%LEVEL] %Msg%n"/>
				<format id="onlymsg" format="%Msg%n"/>
			</formats>
		</seelog>
	`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)
}
