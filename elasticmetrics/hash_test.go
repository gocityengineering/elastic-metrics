package elasticmetrics

import (
	"testing"
)

const (
	stringA = "2021-11-02T11:42:54.894Z eventType=Warning involvedObject=Pod name=filebeat-filebeat-tn588 namespace=elastic-system reason=FailedScheduling message=0/26 nodes are available: 1 Insufficient cpu, 22 node(s) didn't match Pod's node affinity, 3 node(s) had taint {node-role.kubernetes.io/master: }, that the pod didn't tolerate."
	stringB = "2020-11-02T11:42:54.894Z eventType=Warning involvedObject=Pod name=filebeat-filebeat-tn588 namespace=elastic-system reason=FailedScheduling message=0/26 nodes are available: 1 Insufficient cpu, 22 node(s) didn't match Pod's node affinity, 3 node(s) had taint {node-role.kubernetes.io/master: }, that the pod didn't tolerate."
)

func TestHash(t *testing.T) {
	var tests = []struct {
		description string
		a           string
		b           string
		expected    bool
	}{
		{"identityAA", stringA, stringA, true},
		{"contrastAB", stringA, stringB, false},
		{"identityBB", stringB, stringB, true},
		{"contrastBA", stringB, stringA, false},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			hashA := hash(test.a)
			hashB := hash(test.b)
			b := hashA == hashB
			if b != test.expected {
				t.Errorf("Unexpected result matches=%t for hash values %d, %d", b, hashA, hashB)
			}
		})
	}
}
