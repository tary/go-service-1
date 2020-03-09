// Package mdns is a multicast dns registry
package mdns

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/giant-tech/go-service/base/plugin/registry"
	"github.com/micro/mdns"
	hash "github.com/mitchellh/hashstructure"
)

type mdnsTxt struct {
	Service  string
	Version  string
	Metadata map[string]string
}

type mdnsEntry struct {
	hash uint64
	id   string
	node *mdns.Server
}

type mdnsRegistry struct {
	opts registry.Options

	sync.Mutex
	services map[string][]*mdnsEntry
}

func newRegistry(opts ...registry.Option) registry.IRegistry {
	options := registry.Options{
		Timeout: 1 * time.Second,
	}

	registry := &mdnsRegistry{
		opts:     options,
		services: make(map[string][]*mdnsEntry),
	}

	for _, o := range opts {
		o(&registry.opts)
	}

	return registry
}

func (m *mdnsRegistry) Init(opts ...registry.Option) error {
	for _, o := range opts {
		o(&m.opts)
	}
	return nil
}

func (m *mdnsRegistry) Options() registry.Options {
	return m.opts
}

func (m *mdnsRegistry) Register(service *registry.Service, opts ...registry.RegisterOption) error {
	m.Lock()
	defer m.Unlock()
	var mdnsPort int
	if m.opts.Context != nil {
		if v, ok := m.opts.Context.Value(portType{}).(int); ok {
			mdnsPort = v
		}
	}
	entries, ok := m.services[service.Name]
	// first entry, create wildcard used for list queries
	if !ok {
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

	var gerr error

	for _, node := range service.Nodes {
		// create hash of service; uint64
		h, err := hash.Hash(node, nil)
		if err != nil {
			log.Error(err)
			gerr = err
			continue
		}

		var seen bool
		var e *mdnsEntry
		for _, entry := range entries {
			if node.ID == entry.id {
				seen = true
				e = entry
				break
			}
		}

		// already registered, continue
		if seen && e.hash == h {
			continue
			// hash doesn't match, shutdown
		} else if seen {
			log.Infof("id:%s, node hash:%s, old hash:%s. node will restart ...", node.ID, h, e.hash)
			_ = e.node.Shutdown()
			// doesn't exist
		} else {
			e = &mdnsEntry{hash: h}
		}

		txt, err := encode(&mdnsTxt{
			Service:  service.Name,
			Version:  service.Version,
			Metadata: node.Metadata,
		})

		if err != nil {
			log.Error(err)
			gerr = err
			continue
		}

		//
		host, pt, err := net.SplitHostPort(node.Address)
		if err != nil {
			log.Error(err)
			gerr = err
			continue
		}
		port, _ := strconv.Atoi(pt)

		// we got here, new node
		s, err := mdns.NewMDNSService(
			node.ID,
			service.Name,
			"",
			"",
			port,
			[]net.IP{net.ParseIP(host)},
			txt,
		)
		if err != nil {
			log.Error(err)
			gerr = err
			continue
		}

		srv, err := mdns.NewServer(&mdns.Config{Zone: s, Port: mdnsPort})
		if err != nil {
			log.Error(err)
			gerr = err
			continue
		}

		e.id = node.ID
		e.node = srv
		entries = append(entries, e)
	}

	// save
	m.services[service.Name] = entries

	return gerr
}

func (m *mdnsRegistry) Deregister(service *registry.Service) error {
	m.Lock()
	defer m.Unlock()

	var newEntries []*mdnsEntry

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

	return nil
}

func (m *mdnsRegistry) GetService(service string) ([]*registry.Service, error) {
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

				txt, err := decode(e.InfoFields)
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

				serviceMap[txt.Version] = s
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

func (m *mdnsRegistry) ListServices() ([]*registry.Service, error) {
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

/*func (m *mdnsRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}

	md := &mdnsWatcher{
		wo:   wo,
		ch:   make(chan *mdns.ServiceEntry, 32),
		exit: make(chan struct{}),
	}

	go func() {
		if err := mdns.Listen(md.ch, md.exit); err != nil {
			log.Error(err)
			md.Stop()
		}
	}()

	return md, nil
}
*/

func (m *mdnsRegistry) String() string {
	return "mdns"
}

// NewRegistry returns a new default registry which is mdns
func NewRegistry(opts ...registry.Option) registry.IRegistry {
	return newRegistry(opts...)
}
