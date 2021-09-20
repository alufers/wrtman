package main

import (
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/alufers/wrtman/ouidb"
)

// OuiHelper helps with Mac address lookups, by caching them
// TODO: add cache limit
type OuiHelper struct {
	db         *ouidb.OuiDb
	cacheMutex sync.Mutex
	cache      map[string]string
}

func NewOuiHelper() *OuiHelper {
	return &OuiHelper{
		db:    ouidb.New(viper.GetString("extra.oui_db")),
		cache: map[string]string{},
	}
}

func (oh *OuiHelper) LookupVendor(mac string) (string, error) {
	
	oh.cacheMutex.Lock()
	defer oh.cacheMutex.Unlock()
	if cached, ok := oh.cache[mac]; ok {
		return cached, nil
	}

	result, err := oh.db.VendorLookup(mac)

	if err != nil {
		return "", err
	}

	result.Organization = strings.TrimPrefix(result.Organization, "Beijing ")
	oh.cache[mac] = result.Organization
	return result.Organization, nil
}
