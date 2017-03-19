// Copyright 2017 Florin Pățan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deep

import (
	"path/filepath"
	"strings"
)

// Package defines the package format that Deep understands
type Package struct {
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	CommitHash   string    `json:"commit_hash,omitempty"`
	License      string    `json:"license,omitempty"`
	Description  string    `json:"description,omitempty"`
	OSes         []string  `json:"oses,omitempty"`
	MinGoVer     string    `json:"min_go_ver,omitempty"`
	Dependencies []Package `json:"dependencies,omitempty"`
}

func (p Package) isStdlib() bool {
	_, ok := stdlibPackages[p.Name]
	return ok
}

// isSubPackage checks if the current package is a package contained by
// another package.
//
// We need some better detection on this one
func (p Package) isSubPackage() bool {
	if strings.HasPrefix(p.Name, "github.com") {
		// github.com/dlsniper/deep/cmd is contained by github.com/dlsniper/deep
		return strings.Count(p.Name, "/") > 2
	}

	return false
}

func (p Package) isRootPackage(currentPkg string) bool {
	if strings.HasPrefix(p.Name, currentPkg+"/vendor") {
		return false
	}

	return !p.isSubPackage()
}

func (p Package) isThirdParty(currentPkg string) bool {
	if p.Name == currentPkg {
		return false
	}

	// Standard library packages are not considered third-party packages
	if p.isStdlib() {
		return false
	}

	return true
}

func (p Package) vendoredPath(pwd string) string {
	return filepath.Clean(pwd + pathSeparatorString + "vendor" + pathSeparatorString + p.Name)
}

var stdlibPackages = map[string]struct{}{
	"archive/tar":                                 {},
	"archive/zip":                                 {},
	"bufio":                                       {},
	"builtin":                                     {},
	"bytes":                                       {},
	"cmd/addr2line":                               {},
	"cmd/asm":                                     {},
	"cmd/asm/internal/arch":                       {},
	"cmd/asm/internal/asm":                        {},
	"cmd/asm/internal/flags":                      {},
	"cmd/asm/internal/lex":                        {},
	"cmd/cgo":                                     {},
	"cmd/compile":                                 {},
	"cmd/compile/internal/amd64":                  {},
	"cmd/compile/internal/arm":                    {},
	"cmd/compile/internal/arm64":                  {},
	"cmd/compile/internal/gc":                     {},
	"cmd/compile/internal/mips":                   {},
	"cmd/compile/internal/mips64":                 {},
	"cmd/compile/internal/ppc64":                  {},
	"cmd/compile/internal/s390x":                  {},
	"cmd/compile/internal/ssa":                    {},
	"cmd/compile/internal/syntax":                 {},
	"cmd/compile/internal/test":                   {},
	"cmd/compile/internal/x86":                    {},
	"cmd/cover":                                   {},
	"cmd/dist":                                    {},
	"cmd/doc":                                     {},
	"cmd/fix":                                     {},
	"cmd/go":                                      {},
	"cmd/gofmt":                                   {},
	"cmd/internal/bio":                            {},
	"cmd/internal/browser":                        {},
	"cmd/internal/dwarf":                          {},
	"cmd/internal/gcprog":                         {},
	"cmd/internal/goobj":                          {},
	"cmd/internal/obj":                            {},
	"cmd/internal/obj/arm":                        {},
	"cmd/internal/obj/arm64":                      {},
	"cmd/internal/obj/mips":                       {},
	"cmd/internal/obj/ppc64":                      {},
	"cmd/internal/obj/s390x":                      {},
	"cmd/internal/obj/x86":                        {},
	"cmd/internal/objfile":                        {},
	"cmd/internal/sys":                            {},
	"cmd/link":                                    {},
	"cmd/link/internal/amd64":                     {},
	"cmd/link/internal/arm":                       {},
	"cmd/link/internal/arm64":                     {},
	"cmd/link/internal/ld":                        {},
	"cmd/link/internal/mips":                      {},
	"cmd/link/internal/mips64":                    {},
	"cmd/link/internal/ppc64":                     {},
	"cmd/link/internal/s390x":                     {},
	"cmd/link/internal/x86":                       {},
	"cmd/nm":                                      {},
	"cmd/objdump":                                 {},
	"cmd/pack":                                    {},
	"cmd/pprof":                                   {},
	"cmd/pprof/internal/commands":                 {},
	"cmd/pprof/internal/driver":                   {},
	"cmd/pprof/internal/fetch":                    {},
	"cmd/pprof/internal/plugin":                   {},
	"cmd/pprof/internal/report":                   {},
	"cmd/pprof/internal/svg":                      {},
	"cmd/pprof/internal/symbolizer":               {},
	"cmd/pprof/internal/symbolz":                  {},
	"cmd/pprof/internal/tempfile":                 {},
	"cmd/trace":                                   {},
	"cmd/vendor/golang.org/x/arch/arm/armasm":     {},
	"cmd/vendor/golang.org/x/arch/ppc64/ppc64asm": {},
	"cmd/vendor/golang.org/x/arch/x86/x86asm":     {},
	"cmd/vet":                           {},
	"cmd/vet/internal/cfg":              {},
	"cmd/vet/internal/whitelist":        {},
	"compress/bzip2":                    {},
	"compress/flate":                    {},
	"compress/gzip":                     {},
	"compress/lzw":                      {},
	"compress/zlib":                     {},
	"container/heap":                    {},
	"container/list":                    {},
	"container/ring":                    {},
	"context":                           {},
	"crypto":                            {},
	"crypto/aes":                        {},
	"crypto/cipher":                     {},
	"crypto/des":                        {},
	"crypto/dsa":                        {},
	"crypto/ecdsa":                      {},
	"crypto/elliptic":                   {},
	"crypto/hmac":                       {},
	"crypto/internal/cipherhw":          {},
	"crypto/md5":                        {},
	"crypto/rand":                       {},
	"crypto/rc4":                        {},
	"crypto/rsa":                        {},
	"crypto/sha1":                       {},
	"crypto/sha256":                     {},
	"crypto/sha512":                     {},
	"crypto/subtle":                     {},
	"crypto/tls":                        {},
	"crypto/x509":                       {},
	"crypto/x509/pkix":                  {},
	"database/sql":                      {},
	"database/sql/driver":               {},
	"debug/dwarf":                       {},
	"debug/elf":                         {},
	"debug/gosym":                       {},
	"debug/macho":                       {},
	"debug/pe":                          {},
	"debug/plan9obj":                    {},
	"encoding":                          {},
	"encoding/ascii85":                  {},
	"encoding/asn1":                     {},
	"encoding/base32":                   {},
	"encoding/base64":                   {},
	"encoding/binary":                   {},
	"encoding/csv":                      {},
	"encoding/gob":                      {},
	"encoding/hex":                      {},
	"encoding/json":                     {},
	"encoding/pem":                      {},
	"encoding/xml":                      {},
	"errors":                            {},
	"expvar":                            {},
	"flag":                              {},
	"fmt":                               {},
	"go/ast":                            {},
	"go/build":                          {},
	"go/constant":                       {},
	"go/doc":                            {},
	"go/format":                         {},
	"go/importer":                       {},
	"go/internal/gccgoimporter":         {},
	"go/internal/gcimporter":            {},
	"go/parser":                         {},
	"go/printer":                        {},
	"go/scanner":                        {},
	"go/token":                          {},
	"go/types":                          {},
	"hash":                              {},
	"hash/adler32":                      {},
	"hash/crc32":                        {},
	"hash/crc64":                        {},
	"hash/fnv":                          {},
	"html":                              {},
	"html/template":                     {},
	"image":                             {},
	"image/color":                       {},
	"image/color/palette":               {},
	"image/draw":                        {},
	"image/gif":                         {},
	"image/internal/imageutil":          {},
	"image/jpeg":                        {},
	"image/png":                         {},
	"index/suffixarray":                 {},
	"internal/nettrace":                 {},
	"internal/pprof/profile":            {},
	"internal/race":                     {},
	"internal/singleflight":             {},
	"internal/syscall/unix":             {},
	"internal/syscall/windows":          {},
	"internal/syscall/windows/registry": {},
	"internal/syscall/windows/sysdll":   {},
	"internal/testenv":                  {},
	"internal/trace":                    {},
	"io":                                {},
	"io/ioutil":                         {},
	"log":                               {},
	"log/syslog":                        {},
	"math":                              {},
	"math/big":                          {},
	"math/cmplx":                        {},
	"math/rand":                         {},
	"mime":                              {},
	"mime/multipart":                    {},
	"mime/quotedprintable":              {},
	"net":                               {},
	"net/http":                          {},
	"net/http/cgi":                      {},
	"net/http/cookiejar":                {},
	"net/http/fcgi":                     {},
	"net/http/httptest":                 {},
	"net/http/httptrace":                {},
	"net/http/httputil":                 {},
	"net/http/internal":                 {},
	"net/http/pprof":                    {},
	"net/internal/socktest":             {},
	"net/mail":                          {},
	"net/rpc":                           {},
	"net/rpc/jsonrpc":                   {},
	"net/smtp":                          {},
	"net/textproto":                     {},
	"net/url":                           {},
	"os":                                {},
	"os/exec":                           {},
	"os/signal":                         {},
	"os/user":                           {},
	"path":                              {},
	"path/filepath":                     {},
	"plugin":                            {},
	"reflect":                           {},
	"regexp":                            {},
	"regexp/syntax":                     {},
	"runtime":                           {},
	"runtime/cgo":                       {},
	"runtime/debug":                     {},
	"runtime/internal/atomic":           {},
	"runtime/internal/sys":              {},
	"runtime/pprof":                     {},
	"runtime/pprof/internal/protopprof": {},
	"runtime/race":                      {},
	"runtime/trace":                     {},
	"sort":                              {},
	"strconv":                           {},
	"strings":                           {},
	"sync":                              {},
	"sync/atomic":                       {},
	"syscall":                           {},
	"testing":                           {},
	"testing/internal/testdeps":         {},
	"testing/iotest":                    {},
	"testing/quick":                     {},
	"text/scanner":                      {},
	"text/tabwriter":                    {},
	"text/template":                     {},
	"text/template/parse":               {},
	"time":                              {},
	"unicode":                           {},
	"unicode/utf16":                     {},
	"unicode/utf8":                      {},
	"unsafe":                            {},
	"vendor/golang_org/x/crypto/chacha20poly1305":                   {},
	"vendor/golang_org/x/crypto/chacha20poly1305/internal/chacha20": {},
	"vendor/golang_org/x/crypto/curve25519":                         {},
	"vendor/golang_org/x/crypto/poly1305":                           {},
	"vendor/golang_org/x/net/http2/hpack":                           {},
	"vendor/golang_org/x/net/idna":                                  {},
	"vendor/golang_org/x/net/lex/httplex":                           {},
	"vendor/golang_org/x/text/transform":                            {},
	"vendor/golang_org/x/text/unicode/norm":                         {},
	"vendor/golang_org/x/text/width":                                {},
}
