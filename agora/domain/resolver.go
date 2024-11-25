package domain

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Resolver interface {
	Resolve(ctx context.Context, domains []string, regionPrefix string) (string, error)
}

type ResolverFunc func(ctx context.Context, domains []string, regionPrefix string) (string, error)

func (d ResolverFunc) Resolve(ctx context.Context, domains []string, regionPrefix string) (string, error) {
	return d(ctx, domains, regionPrefix)
}

type resolverImpl struct {
	logger log.Logger
	module string
}

func newResolverImpl(logger log.Logger) *resolverImpl {
	return &resolverImpl{logger: logger, module: "resolver"}
}

func (r *resolverImpl) Resolve(ctx context.Context, domains []string, regionPrefix string) (string, error) {
	var wg sync.WaitGroup

	done := make(chan struct{}, 1)
	res := make(chan string, len(domains))
	for _, d := range domains {
		wg.Add(1)

		go func(domain string, regionPrefix string) {
			defer wg.Done()
			url := regionPrefix + "." + domain
			n := time.Now()
			addrs, err := net.LookupHost(url)
			took := time.Since(n)
			r.logger.Debugf(ctx, r.module, "url:%s,IP:%s,took:%s", url, addrs, took.String())
			if err != nil {
				r.logger.Errorf(ctx, r.module, "resolve domain:%s failed,err:%s", url, err.Error())
			} else {
				res <- domain
			}
		}(d, regionPrefix)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case domain := <-res:
		return domain, nil
	case <-done:
	}
	return "", errors.New("query all dns is failed")
}
