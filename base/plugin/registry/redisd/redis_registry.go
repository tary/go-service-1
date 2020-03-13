// Package redisd is redis registry
package redisd

import (
	"context"
	"strings"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/giant-tech/go-service/base/plugin/registry"
	"github.com/micro/mdns"
)

type redisRegistry struct {
	opts registry.Options

	sync.Mutex
	//services map[string][]*mdnsEntry
}

func newRegistry(opts ...registry.Option) registry.IRegistry {
	options := registry.Options{
		Timeout: 1 * time.Second,
	}

	registry := &redisRegistry{
		opts: options,
		//services: make(map[string][]*mdnsEntry),
	}

	for _, o := range opts {
		o(&registry.opts)
	}

	return registry
}

func (m *redisRegistry) Init(opts ...registry.Option) error {
	for _, o := range opts {
		o(&m.opts)
	}
	return nil
}

func (m *redisRegistry) Options() registry.Options {
	return m.opts
}

func (m *redisRegistry) Register(service *registry.Service, opts ...registry.RegisterOption) error {
	m.Lock()
	defer m.Unlock()
	/*var mdnsPort int
	if m.opts.Context != nil {
		if v, ok := m.opts.Context.Value(portType{}).(int); ok {
			mdnsPort = v
		}
	}
	entries, ok := m.services[service.Name]*/
	// first entry, create wildcard used for list queries
	/*if !ok {
		s, err := mdns.NewMDNSService(
			service.Name,
			"_services",
			"",
			"",
			9999,
			[]net.IP{net.ParseIP("0.0.0.0")},
			nil,
		)
		if err != nil {
			log.Error(err)
			return err
		}

		srv, err := mdns.NewServer(&mdns.Config{Zone: &mdns.DNSSDService{MDNSService: s}, Port: mdnsPort})
		if err != nil {
			log.Error(err)
			return err
		}

		// append the wildcard entry
		entries = append(entries, &mdnsEntry{id: "*", node: srv})
	}
	*/

	var gerr error

	// save
	//m.services[service.Name] = entries

	return gerr
}

func (m *redisRegistry) Deregister(service *registry.Service) error {
	m.Lock()
	defer m.Unlock()

	/*	var newEntries []*mdnsEntry

		// loop existing entries, check if any match, shutdown those that do
		for _, entry := range m.services[service.Name] {
			var remove bool

			for _, node := range service.Nodes {
				if node.ID == entry.id {
					_ = entry.node.Shutdown()
					remove = true
					break
				}
			}

			// keep it?
			if !remove {
				newEntries = append(newEntries, entry)
			}
		}

		// last entry is the wildcard for list queries. Remove it.
		if len(newEntries) == 1 && newEntries[0].id == "*" {
			_ = newEntries[0].node.Shutdown()
			delete(m.services, service.Name)
		} else {
			m.services[service.Name] = newEntries
		}
	*/
	return nil
}

func (m *redisRegistry) GetService(service string) ([]*registry.Service, error) {
	serviceMap := make(map[string]*registry.Service)
	entries := make(chan *mdns.ServiceEntry, 10)
	done := make(chan bool)

	p := mdns.DefaultParams(service)
	// set context with timeout
	var cancel context.CancelFunc
	p.Context, cancel = context.WithTimeout(context.Background(), m.opts.Timeout)
	// set entries channel
	p.Entries = entries

	go func() {
		for {
			select {
			case e := <-entries:
				// list record so skip
				if p.Service == "_services" {
					continue
				}

				if e.TTL == 0 {
					log.Errorf("node: %v, ttl is 0", e)
					continue
				}

				/*txt, err := decode(e.InfoFields)
				if err != nil {
					log.Error(err)
					continue
				}

				if txt.Service != service {
					continue
				}

				s, ok := serviceMap[txt.Version]
				if !ok {
					s = &registry.Service{
						Name:    txt.Service,
						Version: txt.Version,
					}
				}

				s.Nodes = append(s.Nodes, &registry.Node{
					ID:       strings.TrimSuffix(e.Name, "."+p.Service+"."+p.Domain+"."),
					Address:  fmt.Sprintf("%s:%d", e.AddrV4.String(), e.Port),
					Metadata: txt.Metadata,
				})

				serviceMap[txt.Version] = s*/
			case <-p.Context.Done():
				cancel()
				close(done)
				return
			}
		}
	}()

	// execute the query
	if err := mdns.Query(p); err != nil {
		log.Error(err)
		return nil, err
	}

	// wait for completion
	<-done

	// create list and return
	var services []*registry.Service

	for _, service := range serviceMap {
		services = append(services, service)
	}

	return services, nil
}

func (m *redisRegistry) ListServices() ([]*registry.Service, error) {
	serviceMap := make(map[string]bool)
	entries := make(chan *mdns.ServiceEntry, 10)
	done := make(chan bool)

	p := mdns.DefaultParams("_services")
	// set context with timeout
	var cancel context.CancelFunc
	p.Context, cancel = context.WithTimeout(context.Background(), m.opts.Timeout)
	// set entries channel
	p.Entries = entries

	var services []*registry.Service

	go func() {
		for {
			select {
			case e := <-entries:
				if e.TTL == 0 {
					continue
				}

				name := strings.TrimSuffix(e.Name, "."+p.Service+"."+p.Domain+".")
				if !serviceMap[name] {
					serviceMap[name] = true
					services = append(services, &registry.Service{Name: name})
				}
			case <-p.Context.Done():
				cancel()
				close(done)
				return
			}
		}
	}()

	// execute query
	if err := mdns.Query(p); err != nil {
		log.Error(err)
		return nil, err
	}

	// wait till done
	<-done

	return services, nil
}

func (m *redisRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}

	md := &redisWatcher{
		wo:   wo,
		ch:   make(chan *mdns.ServiceEntry, 32),
		exit: make(chan struct{}),
	}

	go func() {
		if err := mdns.Listen(md.ch, md.exit); err != nil {
			log.Error(err)
			//md.Stop()
		}
	}()

	return md, nil
}

func (m *redisRegistry) String() string {
	return "redis"
}

// NewRegistry returns a new default registry which is redis
func NewRegistry(opts ...registry.Option) registry.IRegistry {
	return newRegistry(opts...)
}
