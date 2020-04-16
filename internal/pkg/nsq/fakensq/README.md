# Fake NSQ

Fake NSQ Consumer and Producer for `backend/internal/pkg/nsq`

Built to test the correctness of the nsqio wrapper

Limitations:

- Consumer must registered first before publishing message.
- Do not expecting message to be stored, all message directly consumed.
- Message published before any active consumer will be lost.
- Message requeue not working
- Message is always finished