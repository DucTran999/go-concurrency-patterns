package privatechat

import (
	"log"

	"github.com/DucTran999/shared-pkg/scrypto/caesar"
)

type person struct {
	name string // Name of the person

	sendChan     chan string // Channel used by the person to send messages
	receivedChan chan string // Channel used by the person to receive messages
	SecretChan   chan int    // Channel through which the person receives the nonce (secret value) from a generator

	turn     int      // Indicates the person's turn in a communication or protocol sequence
	messages []string // Stores messages sent by the person
}

func NewPerson(
	name string, messages []string, sendChan, receivedChan chan string, turn int,
) *person {
	return &person{
		name:         name,
		messages:     messages,
		sendChan:     sendChan,
		receivedChan: receivedChan,
		SecretChan:   make(chan int),
		turn:         turn,
	}
}

func (p *person) Chat() {
	sent := p.turn
	msgIndex := 0

	for nonce := range p.SecretChan {
		if sent%2 == 0 { // my turn to SEND
			if msgIndex >= len(p.messages) {
				log.Printf("%s: all messages sent - closing chat", p.name)
				close(p.sendChan)
				return
			}
			p.sendChan <- caesar.CaesarEncrypt(p.messages[msgIndex], nonce)
			msgIndex++
		} else { // my turn to RECEIVE
			cipher := <-p.receivedChan
			log.Printf("%s got cipher: %q", p.name, cipher)
			log.Printf("%s received message: %q", p.name, caesar.CaesarDecrypt(cipher, nonce))
			log.Println("---------------------------------------------------------")
		}
		sent++
	}
}
