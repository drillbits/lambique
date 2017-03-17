//    Copyright 2017 drillbits
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package lambique

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	for i, tc := range []struct {
		path    string
		config  string
		expAddr string
		err     error
	}{
		{"generated/tmp/path", `address = "localhost:8080"`, "localhost:8080", nil},
		{"", ``, ":2697", nil},
		{"path/does/not/exist", ``, ":2697", errors.New("open path/does/not/exist: no such file or directory")},
	} {
		i, tc := i, tc
		t.Run("", func(t *testing.T) {
			t.Parallel()

			if tc.config != "" {
				tmpfile, err := ioutil.TempFile("", "testconfig")
				if err != nil {
					t.Fatal(err)
				}
				defer os.Remove(tmpfile.Name())

				_, err = tmpfile.Write([]byte(tc.config))
				if err != nil {
					t.Fatal(err)
				}

				tc.path = tmpfile.Name()
			}

			cfg, err := LoadConfig(tc.path)
			if tc.err != nil {
				if err == nil || err.Error() != tc.err.Error() {
					t.Errorf("%02d: LoadConfig(%s) causes %s, want %s", i, tc.path, err, tc.err)
				}
				if cfg != nil {
					t.Errorf("%02d: LoadConfig(%s) -> %#v, want nil", i, tc.path, cfg)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if cfg.Addr != tc.expAddr {
				t.Errorf("%02d: LoadConfig(%s).Addr -> %s, want %s", i, tc.path, cfg.Addr, tc.expAddr)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	cfg := GetConfig()
	if cfg.Addr != defaultAddr {
		t.Errorf("Config.Addr -> %s, want %s", cfg.Addr, defaultAddr)
	}

	// sigleton
	newAddr := "localhost:8080"
	cfg.Addr = newAddr
	cfg2 := GetConfig()
	if cfg2.Addr != newAddr {
		t.Errorf("Config.Addr -> %s, want %s", cfg2.Addr, newAddr)
	}
}
