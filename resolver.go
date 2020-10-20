package mcping

import (
	"context"
	"net"
	"strconv"
)

// ResolveMinecraftHostPort looks up srv record
func ResolveMinecraftHostPort(ctx context.Context, resolver *net.Resolver, server string) (servers []string, err error) {
	_, addrs, err := resolver.LookupSRV(ctx, "minecraft", "tcp", server)
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		return []string{server}, nil
	}

	servers = make([]string, len(addrs))
	for i := range servers {
		port := strconv.FormatUint(uint64(addrs[i].Port), 10)
		servers[i] = net.JoinHostPort(addrs[i].Target, port)
	}
	return servers, nil
}
