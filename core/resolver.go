package core

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

type DomainResolver interface {
	Resolve(ctx context.Context, domains []string, regionPrefix string) (string, error)
}

type DomainResolverFunc func(ctx context.Context, domains []string, regionPrefix string) (string, error)

func (d DomainResolverFunc) Resolve(ctx context.Context, domains []string, regionPrefix string) (string, error) {
	return d(ctx, domains, regionPrefix)
}

type resolverImpl struct {
	logger Logger
	module string
}

func newResolverImpl(logger Logger) *resolverImpl {
	return &resolverImpl{logger: logger, module: "resolver"}
}

func (r *resolverImpl) Resolve(ctx context.Context, domains []string, regionPrefix string) (string, error) {
	var wg sync.WaitGroup

	done := make(chan struct{}, 1)
	res := make(chan string, 1)
	errCh := make(chan error, len(domains))
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
				errCh <- err
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
