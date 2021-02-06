// Copyright © 2021 Elis Lulja
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event

import "time"

// Options represents an option to add to the event
type Options func(*Event)

// WithCustomMessage adds a message to the event
func WithCustomMessage(msg string) Options {
	return func(e *Event) {
		e.Message = msg
	}
}

// WithProductID adds the product ID to the event
func WithProductID(id string) Options {
	return func(e *Event) {
		e.ProductID = id
	}
}

// WithCustomDate adds a date to the event
func WithCustomDate(date time.Time) Options {
	return func(e *Event) {
		e.Date = date
	}
}

// WithURL adds a URL to the event
func WithURL(url string) Options {
	return func(e *Event) {
		e.URL = url
	}
}

// WithTelegramNotification adds a chat that must be notified with telegram
func WithTelegramNotification(IDs ...string) Options {
	return func(e *Event) {
		e.ToTelegramChat = IDs
	}
}

// WithAttributes adds map of attributes to the event
func WithAttributes(attrs map[string]string) Options {
	return func(e *Event) {
		e.Attributes = attrs
	}
}

// New generates a new event with the given options
func New(eventType Type, websiteName string, opts ...Options) *Event {
	event := &Event{
		Type:           eventType,
		WebsiteName:    websiteName,
		ToTelegramChat: []string{},
	}

	for _, opt := range opts {
		opt(event)
	}

	return event
}
