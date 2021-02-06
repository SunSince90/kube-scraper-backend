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

// Type is the event type that this event represents
type Type string

const (
	// EventGeneral is a general notification
	EventGeneral Type = "general"
	// EventUnexpected means that something unexpected happened to the poller
	EventUnexpected Type = "unexpected"
	// EventError means that an error happened during polling
	EventError Type = "error"
	// EventUnavailable means that the product is not available
	EventUnavailable Type = "unavailable"
	// EventSuccess means that a product is available
	EventSuccess Type = "available"
	// EventProbable means that a product is probably available
	EventProbable Type = "probable"
)
