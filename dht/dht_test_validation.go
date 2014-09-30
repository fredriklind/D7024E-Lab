package dht

import (
	"testing"

	log "github.com/cihub/seelog"
)

// Sets up a custom log receiver that compares log messages
// to the supplied valid logs and fails the test if it finds a
// mismatch.

type TestContext struct {
	validLogs []string
	test      *testing.T
}

type CustomReceiver struct {
	context TestContext
}

func (cr *CustomReceiver) ReceiveMessage(message string, level log.LogLevel, context log.LogContextInterface) error {
	message = message[:len(message)-1]
	if message == cr.context.validLogs[0] {
		cr.context.validLogs = append(cr.context.validLogs[:0], cr.context.validLogs[0+1:]...)
	} else {
		cr.context.test.Fatalf("Expecting '%s', got '%s'\n", cr.context.validLogs[0], message)
	}
	return nil
}
func (cr *CustomReceiver) AfterParse(initArgs log.CustomReceiverInitArgs) error {
	return nil
}
func (cr *CustomReceiver) Flush() {

}
func (cr *CustomReceiver) Close() error {
	return nil
}

func setupTest(t *testing.T, valids []string) {
	c := TestContext{test: t, validLogs: valids}
	testConfig := `
<seelog>
    <outputs>
        <custom name="myreceiver" formatid="test"/>
    </outputs>
    <formats>
        <format id="test" format="%Msg%n"/>
    </formats>
</seelog>
`
	parserParams := &log.CfgParseParams{
		CustomReceiverProducers: map[string]log.CustomReceiverProducer{
			"myreceiver": func(log.CustomReceiverInitArgs) (log.CustomReceiver, error) {
				return &CustomReceiver{c}, nil
			},
		},
	}
	logger, err := log.LoggerFromParamConfigAsString(testConfig, parserParams)
	if err != nil {
		panic(err)
	}
	defer logger.Flush()
	err = log.ReplaceLogger(logger)
	if err != nil {
		panic(err)
	}
}
