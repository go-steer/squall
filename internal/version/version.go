// Copyright 2026 Google LLC
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

// Package version holds build identity stamped via -ldflags at release time.
package version

import "fmt"

// Build variables overridden via -ldflags -X at build time.
var (
	Version = "dev"
	Commit  = "unknown"
)

// Format returns a human-readable build string for CLI and daemon banners.
func Format(binary string) string {
	if Commit == "" || Commit == "unknown" {
		return fmt.Sprintf("%s %s", binary, Version)
	}
	return fmt.Sprintf("%s %s (%s)", binary, Version, Commit)
}
