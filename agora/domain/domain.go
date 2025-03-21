package domain

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/utils"
)

const (
	ChineseMainlandMajorDomain = "sd-rtn.com"
	OverseaMajorDomain         = "agora.io"
)

const GlobalDomainPrefix = "api"

const (
	USWestRegionDomainPrefix = "api-us-west-1"
	USEastRegionDomainPrefix = "api-us-east-1"
)

const (
	APSoutheastRegionDomainPrefix = "api-ap-southeast-1"
	APNortheastRegionDomainPrefix = "api-ap-northeast-1"
)

const (
	EUWestRegionDomainPrefix    = "api-eu-west-1"
	EUCentralRegionDomainPrefix = "api-eu-central-1"
)

const (
	CNEastRegionDomainPrefix  = "api-cn-east-1"
	CNNorthRegionDomainPrefix = "api-cn-north-1"
)

type Domain struct {
	RegionDomainPrefixes []string
	MajorDomainSuffixes  []string
}

var RegionDomain = map[Area]Domain{
	US: {
		RegionDomainPrefixes: []string{
			USWestRegionDomainPrefix,
			USEastRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChineseMainlandMajorDomain,
		},
	},
	EU: {
		RegionDomainPrefixes: []string{
			EUWestRegionDomainPrefix,
			EUCentralRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChineseMainlandMajorDomain,
		},
	},
	AP: {
		RegionDomainPrefixes: []string{
			APSoutheastRegionDomainPrefix,
			APNortheastRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChineseMainlandMajorDomain,
		},
	},
	CN: {
		RegionDomainPrefixes: []string{
			CNEastRegionDomainPrefix,
			CNNorthRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			ChineseMainlandMajorDomain,
			OverseaMajorDomain,
		},
	},
}

type Pool struct {
	domainArea            Area
	domainSuffixes        []string
	currentDomain         string
	regionPrefixes        []string
	currentRegionPrefixes []string
	locker                *sync.Mutex

	resolver   Resolver
	lastUpdate time.Time
	logger     log.Logger
	module     string
}

func NewPool(domainArea Area, logger log.Logger) (*Pool, error) {
	if _, ok := RegionDomain[domainArea]; !ok {
		return nil, errors.New("invalid domain area")
	}
	d := &Pool{
		domainArea:     domainArea,
		domainSuffixes: RegionDomain[domainArea].MajorDomainSuffixes,
		resolver:       newResolverImpl(logger),
		logger:         logger,
		locker:         &sync.Mutex{},
		module:         "domain pool",
	}

	d.regionPrefixes = append(d.regionPrefixes, RegionDomain[domainArea].RegionDomainPrefixes...)

	d.currentRegionPrefixes = d.regionPrefixes
	d.currentDomain = d.domainSuffixes[0]

	return d, nil
}

const updateDuration = 30 * time.Second

func (d *Pool) domainNeedUpdate() bool {
	return time.Since(d.lastUpdate) > updateDuration
}

func (d *Pool) SelectBestDomain(ctx context.Context) error {
	if !d.domainNeedUpdate() {
		return nil
	}

	d.locker.Lock()
	defer d.locker.Unlock()

	if d.domainNeedUpdate() {
		d.logger.Debug(ctx, d.module, "need update domainPool")
		domain, err := d.resolver.Resolve(ctx, d.domainSuffixes, d.currentRegionPrefixes[0])
		if err != nil {
			return err
		}
		d.logger.Debugf(ctx, d.module, "select best domain:%s", domain)
		d.selectDomain(domain)
	}
	return nil
}

func (d *Pool) NextRegion() {
	d.locker.Lock()
	defer d.locker.Unlock()

	d.currentRegionPrefixes = d.currentRegionPrefixes[1:]
	if len(d.currentRegionPrefixes) == 0 {
		d.currentRegionPrefixes = d.regionPrefixes
	}
}

func (d *Pool) selectDomain(domain string) {
	if utils.Contains(d.domainSuffixes, domain) {
		d.currentDomain = domain
		d.lastUpdate = time.Now()
	}
}

func (d *Pool) GetCurrentUrl() string {
	d.locker.Lock()
	defer d.locker.Unlock()

	currentRegion := d.currentRegionPrefixes[0]
	currentDomain := d.currentDomain

	return fmt.Sprintf("https://%s.%s", currentRegion, currentDomain)
}
