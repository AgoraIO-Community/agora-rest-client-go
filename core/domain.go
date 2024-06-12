package core

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	ChinaMainlandMajorDomain = "sd-rtn.com"
	OverseaMajorDomain       = "agora.io"
)

const GlobalDomainPrefix = "api"

type RegionArea int

const (
	USRegionArea = iota
	EURegionArea
	APRegionArea
	CNRegionArea
)

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

var RegionDomain = map[RegionArea]Domain{
	USRegionArea: {
		RegionDomainPrefixes: []string{
			USWestRegionDomainPrefix,
			USEastRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChinaMainlandMajorDomain,
		},
	},
	EURegionArea: {
		RegionDomainPrefixes: []string{
			EUWestRegionDomainPrefix,
			EUCentralRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChinaMainlandMajorDomain,
		},
	},
	APRegionArea: {
		RegionDomainPrefixes: []string{
			APSoutheastRegionDomainPrefix,
			APNortheastRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChinaMainlandMajorDomain,
		},
	},
	CNRegionArea: {
		RegionDomainPrefixes: []string{
			CNEastRegionDomainPrefix,
			CNNorthRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			ChinaMainlandMajorDomain,
			OverseaMajorDomain,
		},
	},
}

type DomainPool struct {
	domainArea            RegionArea
	domainSuffixes        []string
	currentDomain         string
	regionPrefixes        []string
	currentRegionPrefixes []string
	locker                *sync.Mutex

	resolver   DomainResolver
	lastUpdate time.Time
	logger     Logger
	module     string
}

func NewDomainPool(domainArea RegionArea, logger Logger) *DomainPool {
	d := &DomainPool{
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

	return d
}

const updateDuration = 30 * time.Second

func (d *DomainPool) domainNeedUpdate() bool {
	return time.Since(d.lastUpdate) > updateDuration
}

func (d *DomainPool) SelectBestDomain(ctx context.Context) error {
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

func (d *DomainPool) NextRegion() {
	d.locker.Lock()
	defer d.locker.Unlock()

	d.currentRegionPrefixes = d.currentRegionPrefixes[1:]
	if len(d.currentRegionPrefixes) == 0 {
		d.currentRegionPrefixes = d.regionPrefixes
	}
}

func (d *DomainPool) selectDomain(domain string) {
	if Contains(d.domainSuffixes, domain) {
		d.currentDomain = domain
		d.lastUpdate = time.Now()
	}
}

func (d *DomainPool) GetCurrentUrl() string {
	d.locker.Lock()
	defer d.locker.Unlock()

	currentRegion := d.currentRegionPrefixes[0]
	currentDomain := d.currentDomain

	return fmt.Sprintf("https://%s.%s", currentRegion, currentDomain)
}
