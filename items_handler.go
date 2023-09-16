package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/a-h/templ"
	"github.com/danielmmetz/templ/templates"
)

type ItemsHandler struct {
	mu    sync.Mutex
	items []string
}

func (h *ItemsHandler) List(_ *http.Request) (templ.Component, error) {
	return templates.List(h.Items()), nil
}

func (h *ItemsHandler) AddItem(r *http.Request) (templ.Component, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("parsing form: %w", err)
	}
	item := r.PostFormValue("item")
	if item == "" {
		fmt.Println("empty item")
		return templ.NopComponent, nil
	}
	h.addItem(item)
	return templates.AddItem(item), nil
}

func (h *ItemsHandler) DeleteItem(r *http.Request) (templ.Component, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("parsing form: %w", err)
	}
	item := r.FormValue("item")
	h.deleteItem(item)
	return templ.NopComponent, nil
}

func (h *ItemsHandler) Items() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	clone := make([]string, len(h.items))
	copy(clone, h.items)
	return clone
}

func (h *ItemsHandler) addItem(item string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.items = append(h.items, item)
}

func (h *ItemsHandler) deleteItem(item string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for i := range h.items {
		if item == h.items[i] {
			h.items = append(h.items[:i], h.items[i+1:]...)
			return
		}
	}
}
