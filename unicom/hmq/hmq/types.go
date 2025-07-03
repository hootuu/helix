package hmq

type Topic string

type Channel string

type Payload []byte

type Message struct {
	ID        string  `json:"id" bson:"id"`
	Timestamp int64   `json:"timestamp" bson:"timestamp"` //Nano
	Payload   Payload `json:"payload" bson:"payload"`
}

func NewMessage(id string, timestamp int64, payload Payload) *Message {
	return &Message{
		ID:        id,
		Timestamp: timestamp,
		Payload:   payload,
	}
}
