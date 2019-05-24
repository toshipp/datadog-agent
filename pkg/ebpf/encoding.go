package ebpf

import (
	"bytes"

	"github.com/DataDog/datadog-agent/pkg/ebpf/netlink"
	agent "github.com/DataDog/datadog-agent/pkg/process/model"
	"github.com/DataDog/datadog-agent/pkg/process/util"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

var jsonMarshaler = jsonpb.Marshaler{}

// MarshalProtobuf serializes a Connections object into a Protobuf message
func MarshalProtobuf(conns *Connections) ([]byte, error) {
	agentConns := make([]*agent.Connection, len(conns.Conns))
	for i, conn := range conns.Conns {
		agentConns[i] = FormatConnection(conn)
	}
	payload := &agent.Connections{Conns: agentConns}
	return proto.Marshal(payload)
}

// UnmarshalConnections deserializes a Protobuf message into a Connections object
func UnmarshalProtobuf(blob []byte) (*agent.Connections, error) {
	conns := new(agent.Connections)
	if err := proto.Unmarshal(blob, conns); err != nil {
		return nil, err
	}
	return conns, nil
}

// MarshalJSON serializes a Connections object into a JSON document
func MarshalJSON(conns *Connections) ([]byte, error) {
	agentConns := make([]*agent.Connection, len(conns.Conns))
	for i, conn := range conns.Conns {
		agentConns[i] = FormatConnection(conn)
	}
	payload := &agent.Connections{Conns: agentConns}
	writer := new(bytes.Buffer)
	err := jsonMarshaler.Marshal(writer, payload)
	return writer.Bytes(), err
}

// UnmarshalJSON deserializes a JSON document into a Connections object
func UnmarshalJSON(blob []byte) (*agent.Connections, error) {
	conns := new(agent.Connections)
	reader := bytes.NewReader(blob)
	if err := jsonpb.Unmarshal(reader, conns); err != nil {
		return nil, err
	}
	return conns, nil
}

func FormatConnection(conn ConnectionStats) *agent.Connection {
	return &agent.Connection{
		Pid:                int32(conn.Pid),
		Laddr:              formatAddr(conn.Source, conn.SPort),
		Raddr:              formatAddr(conn.Dest, conn.DPort),
		Family:             formatFamily(conn.Family),
		Type:               formatType(conn.Type),
		TotalBytesSent:     conn.MonotonicSentBytes,
		TotalBytesReceived: conn.MonotonicRecvBytes,
		TotalRetransmits:   conn.MonotonicRetransmits,
		LastBytesSent:      conn.LastSentBytes,
		LastBytesReceived:  conn.LastRecvBytes,
		LastRetransmits:    conn.LastRetransmits,
		Direction:          agent.ConnectionDirection(conn.Direction),
		NetNS:              conn.NetNS,
		IpTranslation:      formatIPTranslation(conn.IPTranslation),
	}
}

func formatAddr(addr util.Address, port uint16) *agent.Addr {
	if addr == nil {
		return nil
	}

	return &agent.Addr{Ip: addr.String(), Port: int32(port)}
}

func formatFamily(f ConnectionFamily) agent.ConnectionFamily {
	switch f {
	case AFINET:
		return agent.ConnectionFamily_v4
	case AFINET6:
		return agent.ConnectionFamily_v6
	default:
		return -1
	}
}

func formatType(f ConnectionType) agent.ConnectionType {
	switch f {
	case TCP:
		return agent.ConnectionType_tcp
	case UDP:
		return agent.ConnectionType_udp
	default:
		return -1
	}
}

func formatDirection(d ConnectionDirection) agent.ConnectionDirection {
	switch d {
	case INCOMING:
		return agent.ConnectionDirection_incoming
	case OUTGOING:
		return agent.ConnectionDirection_outgoing
	case LOCAL:
		return agent.ConnectionDirection_local
	default:
		return agent.ConnectionDirection_unspecified
	}
}

func formatIPTranslation(ct *netlink.IPTranslation) *agent.IPTranslation {
	if ct == nil {
		return nil
	}

	return &agent.IPTranslation{
		ReplSrcIP:   ct.ReplSrcIP,
		ReplDstIP:   ct.ReplDstIP,
		ReplSrcPort: int32(ct.ReplSrcPort),
		ReplDstPort: int32(ct.ReplDstPort),
	}
}