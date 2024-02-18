package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentKey = "grpcgateway-user-agent"
	xForwardFor             = " x-forwarded-for"

	userAgentHeader = "user-agent"
	clientIpHeader  = "x-forwarded-for"
)

type MetaData struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(c context.Context) *MetaData {
	mtd := &MetaData{}
	if md, ok := metadata.FromIncomingContext(c); ok {
		if userAgents := md.Get(grpcGatewayUserAgentKey); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}

		if clientIps := md.Get(xForwardFor); len(clientIps) > 0 {
			mtd.ClientIp = clientIps[0]
		}

	}

	if per, ok := peer.FromContext(c); ok {
		mtd.ClientIp = per.Addr.String()
	}

	return mtd
}
