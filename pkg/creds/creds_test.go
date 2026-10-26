/*
Copyright 2026 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package creds

import (
	"testing"

	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/docker/docker-credential-helpers/credentials"
)

// isolateACRHelper makes sure the ACR helper can never reach the network in
// these tests, even if the patched in-tree copy (GO-2026-6225) were lost and
// the unpatched upstream helper were linked instead. The unpatched helper
// accepts lookalike hosts and goes on to read its Azure settings; an invalid
// AZURE_ENVIRONMENT makes that fail before any token request, so the test
// fails with a non "not found" error rather than calling Azure. The patched
// helper rejects the host before reading any settings.
func isolateACRHelper(t *testing.T) {
	t.Helper()
	t.Setenv("AZURE_ENVIRONMENT", "kaniko-unit-test-no-such-cloud")
	for _, v := range []string{"AZURE_CLIENT_ID", "AZURE_CLIENT_SECRET", "AZURE_TENANT_ID", "AZURE_FEDERATED_TOKEN", "AZURE_FEDERATED_TOKEN_FILE"} {
		t.Setenv(v, "")
	}
}

// TestACRHelperIsPatched fails if the go.mod replace pointing at
// third_party/docker-credential-acr-env is dropped. It exercises the helper
// directly rather than through GetKeychain: authn.NewKeychainFromHelper turns
// any helper error into Anonymous, so a keychain-level check could not tell
// the patched and unpatched helpers apart.
func TestACRHelperIsPatched(t *testing.T) {
	isolateACRHelper(t)

	for _, h := range []string{"evil.azurecr.io.attacker.com", "foo.azurecr.io.", "foo.azurecr.io@attacker.com"} {
		_, _, err := credhelper.NewACRCredentialsHelper().Get(h)
		if !credentials.IsErrCredentialsNotFound(err) {
			t.Errorf("ACR helper Get(%q): err = %v, want credentials not found", h, err)
		}
	}
}
