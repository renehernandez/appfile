package do

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/renehernandez/appfile/internal/log"
)

type DomainService struct {
	client *godo.Client
}

func NewDomainService(token string) *DomainService {
	return &DomainService{
		client: godo.NewFromToken(token),
	}
}

func (svc *DomainService) DeleteRecord(name string, domain string) error {
	ctx := context.TODO()
	record, err := svc.getCNAMERecord(name, domain)
	if err != nil {
		return err
	}
	if record.ID > 0 {
		log.Debugf("Record to delete: %++v", record)
		_, err = svc.client.Domains.DeleteRecord(ctx, domain, record.ID)
	}
	return err
}

func (svc *DomainService) getCNAMERecord(name string, domain string) (*godo.DomainRecord, error) {
	ctx := context.TODO()
	opts := &godo.ListOptions{}

	records, _, err := svc.client.Domains.RecordsByTypeAndName(ctx, domain, "CNAME", name, opts)

	if err != nil {
		return &godo.DomainRecord{}, errors.Wrapf(err, "Failed to retrieve %s record from DigitalOcean", name)
	}

	if len(records) == 0 {
		log.Warningf("%s CNAME record not found", name)
		return &godo.DomainRecord{}, nil
	} else if len(records) > 1 {
		return &godo.DomainRecord{}, fmt.Errorf("Same %s CNAME record appeared more than once", name)
	}
	return &records[0], nil
}
