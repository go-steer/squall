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

package version

import "testing"

func TestFormat(t *testing.T) {
	origVersion, origCommit := Version, Commit
	t.Cleanup(func() {
		Version = origVersion
		Commit = origCommit
	})

	tests := []struct {
		name    string
		binary  string
		version string
		commit  string
		want    string
	}{
		{
			name:    "default dev build",
			binary:  "squall",
			version: "dev",
			commit:  "unknown",
			want:    "squall dev",
		},
		{
			name:    "empty commit",
			binary:  "squalld",
			version: "v0.1.0",
			commit:  "",
			want:    "squalld v0.1.0",
		},
		{
			name:    "stamped release build",
			binary:  "squall",
			version: "v0.1.0",
			commit:  "a8719f9",
			want:    "squall v0.1.0 (a8719f9)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			Version = tc.version
			Commit = tc.commit
			if got := Format(tc.binary); got != tc.want {
				t.Fatalf("Format(%q) = %q, want %q", tc.binary, got, tc.want)
			}
		})
	}
}
