package model

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"
)

type DNSRecordMeta struct {
	AutoAdd       bool   `json:"auto_added"`
	ManagedByApps bool   `json:"managed_by_apps"`
	ManagedByArgo bool   `json:"managed_by_argo_tunnel"`
	Source        string `json:"Source"`
}

func (r *DNSRecordMeta) String() string {
	buf := bytes.NewBuffer(nil)
	w := tabwriter.NewWriter(buf, 10, 20, 1, '.', tabwriter.TabIndent)

    fmt.Fprintf(w, "AutoAdd\t: %t\n", r.AutoAdd)
    fmt.Fprintf(w, "ManagedByApps\t: %t\n", r.ManagedByApps)
    fmt.Fprintf(w, "ManagedByArgo\t: %t\n", r.ManagedByArgo)
    fmt.Fprintf(w, "Source\t: %s", r.Source)

    w.Flush()
    return buf.String()
}

type DNSDeleteRequest struct {
	ID     string
	ZoneID string
}

func (r *DNSDeleteRequest) String() string {
	buf := bytes.NewBuffer(nil)
	w := tabwriter.NewWriter(buf, 10, 20, 1, '.', tabwriter.TabIndent)

    fmt.Fprintf(w, "ID\t: %s\n", r.ID)
    fmt.Fprintf(w, "ZoneID\t: %s", r.ZoneID)

    w.Flush()
    return buf.String()
}

type DNSRecordRequest struct {
	ID       string
	ZoneID   string
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	Proxied  bool   `json:"proxied"`
	Priority int    `json:"priority"`
	TTL      int    `json:"ttl"`
}

func (r *DNSRecordRequest) String() string {

	buf := bytes.NewBuffer(nil)
	w := tabwriter.NewWriter(buf, 10, 20, 1, '.', tabwriter.TabIndent)

	fmt.Fprintf(w, "ID\t: %s\n", r.ID)
	fmt.Fprintf(w, "ZoneID\t: %s\n", r.ZoneID)
	fmt.Fprintf(w, "Name\t: %s\n", r.Name)
	fmt.Fprintf(w, "Type\t: %s\n", r.Type)
	fmt.Fprintf(w, "Content\t: %s\n", r.Content)
	fmt.Fprintf(w, "Proxied\t: %v\n", r.Proxied)
	fmt.Fprintf(w, "Priority\t: %d\n", r.Priority)
	fmt.Fprintf(w, "TTL\t: %d", r.TTL)

	w.Flush()
	return buf.String()
}

func (r *DNSRecordRequest) MustSanitize() {
	if err := r.Sanitize(); err != nil {
		panic(err)
	}
}

func (r *DNSRecordRequest) Sanitize() error {
	r.Type = strings.TrimSpace(r.Type)
	if r.Type == "" {
		r.Type = "A"
	}

	r.Content = strings.TrimSpace(r.Content)
	if r.Content == "" {
		return fmt.Errorf("Content cannot be empty")
	}

	if r.ZoneID == "" {
		return fmt.Errorf("Zone ID cannot be empty")
	}


	return nil
}

type DNSRecord struct {
	ID        string         `json:"id"`
	ZoneID    string         `json:"zone_id"`
	ZoneName  string         `json:"zone_name"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Content   string         `json:"content"`
	Proxiable bool           `json:"proxiable"`
	Proxied   bool           `json:"proxied"`
	TTL       int            `json:"ttl"`
	Locked    bool           `json:"locked"`
	Created   *time.Time     `json:"Created"`
	Modified  *time.Time     `json:"Modified"`
	Meta      *DNSRecordMeta `json:"meta"`
}

func (r *DNSRecord) String() string {
    buf := bytes.NewBuffer(nil)
    w := tabwriter.NewWriter(buf, 10, 20, 1, '.', tabwriter.TabIndent)

    fmt.Fprintf(w, "ID\t: %s\n", r.ID)
    fmt.Fprintf(w, "ZoneID\t: %s\n", r.ZoneID)
    fmt.Fprintf(w, "ZoneName\t: %s\n", r.ZoneName)
    fmt.Fprintf(w, "Name\t: %s\n", r.Name)
    fmt.Fprintf(w, "Type\t: %s\n", r.Type)
    fmt.Fprintf(w, "Content\t: %s\n", r.Content)
    fmt.Fprintf(w, "Proxiable\t: %t\n", r.Proxiable)
    fmt.Fprintf(w, "Proxied\t: %t\n", r.Proxied)
    fmt.Fprintf(w, "TTL\t: %d\n", r.TTL)
    fmt.Fprintf(w, "Locked\t: %t\n", r.Locked)
    fmt.Fprintf(w, "Created\t: %s\n", r.Created)
    fmt.Fprintf(w, "Modified\t: %s\n", r.Modified)
    fmt.Fprintf(w, "Meta\n%s", r.Meta.String())

    w.Flush()
    return buf.String()
}

type Zone struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Status          string     `json:"status"`
	Paused          bool       `json:"paused"`
	Type            string     `json:"type"`
	DevelopmentMode int        `json:"development_mode"`
	NameServers     []string   `json:"name_servers"`
	OrigNameServers []string   `json:"original_name_servers"`
	OrigRegistrar   string     `json:"original_registrar"`
	Created         *time.Time `json:"created_on"`
	Modified        *time.Time `json:"modified_on"`
	Activated       *time.Time `json:"activated_on"`
}

func (z *Zone) String() string {
    buf := bytes.NewBuffer(nil)
    w := tabwriter.NewWriter(buf, 10, 20, 1, '.', tabwriter.TabIndent)

    fmt.Fprintf(w, "ID\t: %s\n", z.ID)
    fmt.Fprintf(w, "Name\t: %s\n", z.Name)
    fmt.Fprintf(w, "Status\t: %s\n", z.Status)
    fmt.Fprintf(w, "Paused\t: %s\n", z.Paused)
    fmt.Fprintf(w, "Type\t: %s\n", z.Type)
    fmt.Fprintf(w, "DevelopmentMode\t: %s\n", z.DevelopmentMode)
    fmt.Fprintf(w, "NameServers\t: %s\n", z.NameServers)
    fmt.Fprintf(w, "OrigNameServers\t: %s\n", z.OrigNameServers)
    fmt.Fprintf(w, "OrigRegistrar\t: %s\n", z.OrigRegistrar)
    fmt.Fprintf(w, "Created\t: %s\n", z.Created)
    fmt.Fprintf(w, "Modified\t: %s\n", z.Modified)
    fmt.Fprintf(w, "Activated\t: %s", z.Activated)

    w.Flush()
    return buf.String()
}
