package kv

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/docker/docker/pkg/discovery"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/go-check/check"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) ***REMOVED*** check.TestingT(t) ***REMOVED***

type DiscoverySuite struct***REMOVED******REMOVED***

var _ = check.Suite(&DiscoverySuite***REMOVED******REMOVED***)

func (ds *DiscoverySuite) TestInitialize(c *check.C) ***REMOVED***
	storeMock := &FakeStore***REMOVED***
		Endpoints: []string***REMOVED***"127.0.0.1"***REMOVED***,
	***REMOVED***
	d := &Discovery***REMOVED***backend: store.CONSUL***REMOVED***
	d.Initialize("127.0.0.1", 0, 0, nil)
	d.store = storeMock

	s := d.store.(*FakeStore)
	c.Assert(s.Endpoints, check.HasLen, 1)
	c.Assert(s.Endpoints[0], check.Equals, "127.0.0.1")
	c.Assert(d.path, check.Equals, defaultDiscoveryPath)

	storeMock = &FakeStore***REMOVED***
		Endpoints: []string***REMOVED***"127.0.0.1:1234"***REMOVED***,
	***REMOVED***
	d = &Discovery***REMOVED***backend: store.CONSUL***REMOVED***
	d.Initialize("127.0.0.1:1234/path", 0, 0, nil)
	d.store = storeMock

	s = d.store.(*FakeStore)
	c.Assert(s.Endpoints, check.HasLen, 1)
	c.Assert(s.Endpoints[0], check.Equals, "127.0.0.1:1234")
	c.Assert(d.path, check.Equals, "path/"+defaultDiscoveryPath)

	storeMock = &FakeStore***REMOVED***
		Endpoints: []string***REMOVED***"127.0.0.1:1234", "127.0.0.2:1234", "127.0.0.3:1234"***REMOVED***,
	***REMOVED***
	d = &Discovery***REMOVED***backend: store.CONSUL***REMOVED***
	d.Initialize("127.0.0.1:1234,127.0.0.2:1234,127.0.0.3:1234/path", 0, 0, nil)
	d.store = storeMock

	s = d.store.(*FakeStore)
	c.Assert(s.Endpoints, check.HasLen, 3)
	c.Assert(s.Endpoints[0], check.Equals, "127.0.0.1:1234")
	c.Assert(s.Endpoints[1], check.Equals, "127.0.0.2:1234")
	c.Assert(s.Endpoints[2], check.Equals, "127.0.0.3:1234")

	c.Assert(d.path, check.Equals, "path/"+defaultDiscoveryPath)
***REMOVED***

// Extremely limited mock store so we can test initialization
type Mock struct ***REMOVED***
	// Endpoints passed to InitializeMock
	Endpoints []string

	// Options passed to InitializeMock
	Options *store.Config
***REMOVED***

func NewMock(endpoints []string, options *store.Config) (store.Store, error) ***REMOVED***
	s := &Mock***REMOVED******REMOVED***
	s.Endpoints = endpoints
	s.Options = options
	return s, nil
***REMOVED***
func (s *Mock) Put(key string, value []byte, opts *store.WriteOptions) error ***REMOVED***
	return errors.New("Put not supported")
***REMOVED***
func (s *Mock) Get(key string) (*store.KVPair, error) ***REMOVED***
	return nil, errors.New("Get not supported")
***REMOVED***
func (s *Mock) Delete(key string) error ***REMOVED***
	return errors.New("Delete not supported")
***REMOVED***

// Exists mock
func (s *Mock) Exists(key string) (bool, error) ***REMOVED***
	return false, errors.New("Exists not supported")
***REMOVED***

// Watch mock
func (s *Mock) Watch(key string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan *store.KVPair, error) ***REMOVED***
	return nil, errors.New("Watch not supported")
***REMOVED***

// WatchTree mock
func (s *Mock) WatchTree(prefix string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan []*store.KVPair, error) ***REMOVED***
	return nil, errors.New("WatchTree not supported")
***REMOVED***

// NewLock mock
func (s *Mock) NewLock(key string, options *store.LockOptions) (store.Locker, error) ***REMOVED***
	return nil, errors.New("NewLock not supported")
***REMOVED***

// List mock
func (s *Mock) List(prefix string) ([]*store.KVPair, error) ***REMOVED***
	return nil, errors.New("List not supported")
***REMOVED***

// DeleteTree mock
func (s *Mock) DeleteTree(prefix string) error ***REMOVED***
	return errors.New("DeleteTree not supported")
***REMOVED***

// AtomicPut mock
func (s *Mock) AtomicPut(key string, value []byte, previous *store.KVPair, opts *store.WriteOptions) (bool, *store.KVPair, error) ***REMOVED***
	return false, nil, errors.New("AtomicPut not supported")
***REMOVED***

// AtomicDelete mock
func (s *Mock) AtomicDelete(key string, previous *store.KVPair) (bool, error) ***REMOVED***
	return false, errors.New("AtomicDelete not supported")
***REMOVED***

// Close mock
func (s *Mock) Close() ***REMOVED***
***REMOVED***

func (ds *DiscoverySuite) TestInitializeWithCerts(c *check.C) ***REMOVED***
	cert := `-----BEGIN CERTIFICATE-----
MIIDCDCCAfKgAwIBAgIICifG7YeiQOEwCwYJKoZIhvcNAQELMBIxEDAOBgNVBAMT
B1Rlc3QgQ0EwHhcNMTUxMDAxMjMwMDAwWhcNMjAwOTI5MjMwMDAwWjASMRAwDgYD
VQQDEwdUZXN0IENBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1wRC
O+flnLTK5ImjTurNRHwSejuqGbc4CAvpB0hS+z0QlSs4+zE9h80aC4hz+6caRpds
+J908Q+RvAittMHbpc7VjbZP72G6fiXk7yPPl6C10HhRSoSi3nY+B7F2E8cuz14q
V2e+ejhWhSrBb/keyXpcyjoW1BOAAJ2TIclRRkICSCZrpXUyXxAvzXfpFXo1RhSb
UywN11pfiCQzDUN7sPww9UzFHuAHZHoyfTr27XnJYVUerVYrCPq8vqfn//01qz55
Xs0hvzGdlTFXhuabFtQnKFH5SNwo/fcznhB7rePOwHojxOpXTBepUCIJLbtNnWFT
V44t9gh5IqIWtoBReQIDAQABo2YwZDAOBgNVHQ8BAf8EBAMCAAYwEgYDVR0TAQH/
BAgwBgEB/wIBAjAdBgNVHQ4EFgQUZKUI8IIjIww7X/6hvwggQK4bD24wHwYDVR0j
BBgwFoAUZKUI8IIjIww7X/6hvwggQK4bD24wCwYJKoZIhvcNAQELA4IBAQDES2cz
7sCQfDCxCIWH7X8kpi/JWExzUyQEJ0rBzN1m3/x8ySRxtXyGekimBqQwQdFqlwMI
xzAQKkh3ue8tNSzRbwqMSyH14N1KrSxYS9e9szJHfUasoTpQGPmDmGIoRJuq1h6M
ej5x1SCJ7GWCR6xEXKUIE9OftXm9TdFzWa7Ja3OHz/mXteii8VXDuZ5ACq6EE5bY
8sP4gcICfJ5fTrpTlk9FIqEWWQrCGa5wk95PGEj+GJpNogjXQ97wVoo/Y3p1brEn
t5zjN9PAq4H1fuCMdNNA+p1DHNwd+ELTxcMAnb2ajwHvV6lKPXutrTFc4umJToBX
FpTxDmJHEV4bzUzh
-----END CERTIFICATE-----
`
	key := `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA1wRCO+flnLTK5ImjTurNRHwSejuqGbc4CAvpB0hS+z0QlSs4
+zE9h80aC4hz+6caRpds+J908Q+RvAittMHbpc7VjbZP72G6fiXk7yPPl6C10HhR
SoSi3nY+B7F2E8cuz14qV2e+ejhWhSrBb/keyXpcyjoW1BOAAJ2TIclRRkICSCZr
pXUyXxAvzXfpFXo1RhSbUywN11pfiCQzDUN7sPww9UzFHuAHZHoyfTr27XnJYVUe
rVYrCPq8vqfn//01qz55Xs0hvzGdlTFXhuabFtQnKFH5SNwo/fcznhB7rePOwHoj
xOpXTBepUCIJLbtNnWFTV44t9gh5IqIWtoBReQIDAQABAoIBAHSWipORGp/uKFXj
i/mut776x8ofsAxhnLBARQr93ID+i49W8H7EJGkOfaDjTICYC1dbpGrri61qk8sx
qX7p3v/5NzKwOIfEpirgwVIqSNYe/ncbxnhxkx6tXtUtFKmEx40JskvSpSYAhmmO
1XSx0E/PWaEN/nLgX/f1eWJIlxlQkk3QeqL+FGbCXI48DEtlJ9+MzMu4pAwZTpj5
5qtXo5JJ0jRGfJVPAOznRsYqv864AhMdMIWguzk6EGnbaCWwPcfcn+h9a5LMdony
MDHfBS7bb5tkF3+AfnVY3IBMVx7YlsD9eAyajlgiKu4zLbwTRHjXgShy+4Oussz0
ugNGnkECgYEA/hi+McrZC8C4gg6XqK8+9joD8tnyDZDz88BQB7CZqABUSwvjDqlP
L8hcwo/lzvjBNYGkqaFPUICGWKjeCtd8pPS2DCVXxDQX4aHF1vUur0uYNncJiV3N
XQz4Iemsa6wnKf6M67b5vMXICw7dw0HZCdIHD1hnhdtDz0uVpeevLZ8CgYEA2KCT
Y43lorjrbCgMqtlefkr3GJA9dey+hTzCiWEOOqn9RqGoEGUday0sKhiLofOgmN2B
LEukpKIey8s+Q/cb6lReajDVPDsMweX8i7hz3Wa4Ugp4Xa5BpHqu8qIAE2JUZ7bU
t88aQAYE58pUF+/Lq1QzAQdrjjzQBx6SrBxieecCgYEAvukoPZEC8mmiN1VvbTX+
QFHmlZha3QaDxChB+QUe7bMRojEUL/fVnzkTOLuVFqSfxevaI/km9n0ac5KtAchV
xjp2bTnBb5EUQFqjopYktWA+xO07JRJtMfSEmjZPbbay1kKC7rdTfBm961EIHaRj
xZUf6M+rOE8964oGrdgdLlECgYEA046GQmx6fh7/82FtdZDRQp9tj3SWQUtSiQZc
qhO59Lq8mjUXz+MgBuJXxkiwXRpzlbaFB0Bca1fUoYw8o915SrDYf/Zu2OKGQ/qa
V81sgiVmDuEgycR7YOlbX6OsVUHrUlpwhY3hgfMe6UtkMvhBvHF/WhroBEIJm1pV
PXZ/CbMCgYEApNWVktFBjOaYfY6SNn4iSts1jgsQbbpglg3kT7PLKjCAhI6lNsbk
dyT7ut01PL6RaW4SeQWtrJIVQaM6vF3pprMKqlc5XihOGAmVqH7rQx9rtQB5TicL
BFrwkQE4HQtQBV60hYQUzzlSk44VFDz+jxIEtacRHaomDRh2FtOTz+I=
-----END RSA PRIVATE KEY-----
`
	certFile, err := ioutil.TempFile("", "cert")
	c.Assert(err, check.IsNil)
	defer os.Remove(certFile.Name())
	certFile.Write([]byte(cert))
	certFile.Close()
	keyFile, err := ioutil.TempFile("", "key")
	c.Assert(err, check.IsNil)
	defer os.Remove(keyFile.Name())
	keyFile.Write([]byte(key))
	keyFile.Close()

	libkv.AddStore("mock", NewMock)
	d := &Discovery***REMOVED***backend: "mock"***REMOVED***
	err = d.Initialize("127.0.0.3:1234", 0, 0, map[string]string***REMOVED***
		"kv.cacertfile": certFile.Name(),
		"kv.certfile":   certFile.Name(),
		"kv.keyfile":    keyFile.Name(),
	***REMOVED***)
	c.Assert(err, check.IsNil)
	s := d.store.(*Mock)
	c.Assert(s.Options.TLS, check.NotNil)
	c.Assert(s.Options.TLS.RootCAs, check.NotNil)
	c.Assert(s.Options.TLS.Certificates, check.HasLen, 1)
***REMOVED***

func (ds *DiscoverySuite) TestWatch(c *check.C) ***REMOVED***
	mockCh := make(chan []*store.KVPair)

	storeMock := &FakeStore***REMOVED***
		Endpoints:  []string***REMOVED***"127.0.0.1:1234"***REMOVED***,
		mockKVChan: mockCh,
	***REMOVED***

	d := &Discovery***REMOVED***backend: store.CONSUL***REMOVED***
	d.Initialize("127.0.0.1:1234/path", 0, 0, nil)
	d.store = storeMock

	expected := discovery.Entries***REMOVED***
		&discovery.Entry***REMOVED***Host: "1.1.1.1", Port: "1111"***REMOVED***,
		&discovery.Entry***REMOVED***Host: "2.2.2.2", Port: "2222"***REMOVED***,
	***REMOVED***
	kvs := []*store.KVPair***REMOVED***
		***REMOVED***Key: path.Join("path", defaultDiscoveryPath, "1.1.1.1"), Value: []byte("1.1.1.1:1111")***REMOVED***,
		***REMOVED***Key: path.Join("path", defaultDiscoveryPath, "2.2.2.2"), Value: []byte("2.2.2.2:2222")***REMOVED***,
	***REMOVED***

	stopCh := make(chan struct***REMOVED******REMOVED***)
	ch, errCh := d.Watch(stopCh)

	// It should fire an error since the first WatchTree call failed.
	c.Assert(<-errCh, check.ErrorMatches, "test error")
	// We have to drain the error channel otherwise Watch will get stuck.
	go func() ***REMOVED***
		for range errCh ***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Push the entries into the store channel and make sure discovery emits.
	mockCh <- kvs
	c.Assert(<-ch, check.DeepEquals, expected)

	// Add a new entry.
	expected = append(expected, &discovery.Entry***REMOVED***Host: "3.3.3.3", Port: "3333"***REMOVED***)
	kvs = append(kvs, &store.KVPair***REMOVED***Key: path.Join("path", defaultDiscoveryPath, "3.3.3.3"), Value: []byte("3.3.3.3:3333")***REMOVED***)
	mockCh <- kvs
	c.Assert(<-ch, check.DeepEquals, expected)

	close(mockCh)
	// Give it enough time to call WatchTree.
	time.Sleep(3 * time.Second)

	// Stop and make sure it closes all channels.
	close(stopCh)
	c.Assert(<-ch, check.IsNil)
	c.Assert(<-errCh, check.IsNil)
***REMOVED***

// FakeStore implements store.Store methods. It mocks all store
// function in a simple, naive way.
type FakeStore struct ***REMOVED***
	Endpoints  []string
	Options    *store.Config
	mockKVChan <-chan []*store.KVPair

	watchTreeCallCount int
***REMOVED***

func (s *FakeStore) Put(key string, value []byte, options *store.WriteOptions) error ***REMOVED***
	return nil
***REMOVED***

func (s *FakeStore) Get(key string) (*store.KVPair, error) ***REMOVED***
	return nil, nil
***REMOVED***

func (s *FakeStore) Delete(key string) error ***REMOVED***
	return nil
***REMOVED***

func (s *FakeStore) Exists(key string) (bool, error) ***REMOVED***
	return true, nil
***REMOVED***

func (s *FakeStore) Watch(key string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan *store.KVPair, error) ***REMOVED***
	return nil, nil
***REMOVED***

// WatchTree will fail the first time, and return the mockKVchan afterwards.
// This is the behavior we need for testing.. If we need 'moar', should update this.
func (s *FakeStore) WatchTree(directory string, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan []*store.KVPair, error) ***REMOVED***
	if s.watchTreeCallCount == 0 ***REMOVED***
		s.watchTreeCallCount = 1
		return nil, errors.New("test error")
	***REMOVED***
	// First calls error
	return s.mockKVChan, nil
***REMOVED***

func (s *FakeStore) NewLock(key string, options *store.LockOptions) (store.Locker, error) ***REMOVED***
	return nil, nil
***REMOVED***

func (s *FakeStore) List(directory string) ([]*store.KVPair, error) ***REMOVED***
	return []*store.KVPair***REMOVED******REMOVED***, nil
***REMOVED***

func (s *FakeStore) DeleteTree(directory string) error ***REMOVED***
	return nil
***REMOVED***

func (s *FakeStore) AtomicPut(key string, value []byte, previous *store.KVPair, options *store.WriteOptions) (bool, *store.KVPair, error) ***REMOVED***
	return true, nil, nil
***REMOVED***

func (s *FakeStore) AtomicDelete(key string, previous *store.KVPair) (bool, error) ***REMOVED***
	return true, nil
***REMOVED***

func (s *FakeStore) Close() ***REMOVED***
***REMOVED***
