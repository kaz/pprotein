package event

type (
	Publisher struct {
		hub     *Hub
		message string
	}
)

func (p *Publisher) Publish() {
	p.hub.publish(p.message)
}
