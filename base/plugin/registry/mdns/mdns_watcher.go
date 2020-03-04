package mdns

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/registry"
	"github.com/micro/mdns"
)

type mdnsWatcher struct {
	wo   registry.WatchOptions
	ch   chan *mdns.ServiceEntry
	exit chan struct{}
}

func (m *mdnsWatcher) Next() (*registry.Result, error) {
	for {
		select {
		case e := <-m.ch:
			txt, err := decode(e.InfoFields)
			if err != nil {
				log.Error(err)
				continue
			}

			if len(txt.Service) == 0 || len(txt.Version) == 0 {
				continue
			}

			// Filter watch options
			// wo.Service: Only keep services we care about
			if len(m.wo.Service) > 0 && txt.Service != m.wo.Service {
				continue
			}

			var action string

			if e.TTL == 0 {
				action = "delete"
			} else {
				action = "create"
			}

			service := &registry.Service{
				Name:    txt.Service,
				Version: txt.Version,
			}

			if !strings.HasSuffix(e.Name, "."+service.Name+".local.") {
				continue
			}

			service.Nodes = append(service.Nodes, &registry.Node{
				ID:       strings.TrimSuffix(e.Name, "."+service.Name+".local."),
				Address:  fmt.Sprintf("%s:%d", e.AddrV4.String(), e.Port),
				Metadata: txt.Metadata,
			})

			return &registry.Result{
				Action:  action,
				Service: service,
			}, nil
		case <-m.exit:
			return nil, errors.New("watcher stopped")
		}
	}
}

func (m *mdnsWatcher) Stop() {
	select {
	case <-m.exit:
		return
	default:
		close(m.exit)
	}
}
