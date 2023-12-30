package user

import (
	"testing"

	"github.com/Jiang-Gianni/chat/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
)

// Test both login and register
// The test assumes that the reply events
// are emitted in the same order as the request events
func TestNATSLoginRegister(t *testing.T) {
	t.Parallel()
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	db, cleanup, err := config.GetSqliteTest()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	q := *New(db)
	ns := NATSServer{
		NATS:    nc,
		Queries: q,
	}

	// Start the service
	go ns.Run()

	// A new NATS connection has to be instantiated
	// If both client and server share the same connection
	// then the same connection is publishing and subscribing
	// to the same topic -> no message is received
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		t.Fatal(err)
	}

	len := LoginEventNATS{
		Username: "my-test-username",
		Password: "my-test-password",
		ReplyTo:  nuid.Next(),
	}
	loginResp := &LoginReplyNATS{}
	ren := RegisterEventNATS{
		Username: "my-test-username",
		Password: "my-test-password",
		ReplyTo:  nuid.Next(),
	}
	registerResp := &RegisterReplyNATS{}

	loginChan := make(chan *LoginReplyNATS)
	loginSub, err := ec.Subscribe(len.ReplyTo, func(lrn *LoginReplyNATS) {
		loginChan <- lrn
	})
	if err != nil {
		t.Fatal(err)
	}
	defer loginSub.Unsubscribe()

	registerChan := make(chan *RegisterReplyNATS)
	registerSub, err := ec.Subscribe(ren.ReplyTo, func(rrn *RegisterReplyNATS) {
		registerChan <- rrn
	})
	if err != nil {
		t.Fatal(err)
	}
	defer registerSub.Unsubscribe()

	// Invalid Credentials
	err = ec.Publish(config.NATSUserLogin, len)
	if err != nil {
		t.Fatal(err)
	}
	loginResp = <-loginChan
	if loginResp.StatusCode != StatusInvalidCredentials {
		t.Fatal("should not be able to login without registering")
	}

	// Successful register
	err = ec.Publish(config.NATSUserRegister, ren)
	if err != nil {
		t.Fatal(err)
	}
	registerResp = <-registerChan
	if registerResp.StatusCode != StatusOK {
		t.Fatal("registration should have gone OK")
	}

	// Username Taken
	err = ec.Publish(config.NATSUserRegister, ren)
	if err != nil {
		t.Fatal(err)
	}
	registerResp = <-registerChan
	if registerResp.StatusCode != StatusUsernameTaken {
		t.Fatal("username should already be taken")
	}

	// Wrong Password
	wrong := len
	wrong.Password = "wrong-pw"
	err = ec.Publish(config.NATSUserLogin, wrong)
	if err != nil {
		t.Fatal(err)
	}
	loginResp = <-loginChan
	if loginResp.StatusCode != StatusInvalidCredentials {
		t.Fatal("status should be invalid credentials")
	}

	// Successful login
	err = ec.Publish(config.NATSUserLogin, len)
	if err != nil {
		t.Fatal(err)
	}
	loginResp = <-loginChan
	if loginResp.StatusCode != StatusOK {
		t.Fatal("failed to login")
	}
}
